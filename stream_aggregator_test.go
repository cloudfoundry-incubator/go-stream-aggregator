package streamaggregator_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/apoydence/onpar"
	. "github.com/apoydence/onpar/expect"
	. "github.com/apoydence/onpar/matchers"
	streamaggregator "github.com/apoydence/stream-aggregator"
)

type TSA struct {
	*testing.T
	a *streamaggregator.StreamAggregator
}

func TestStreamAggregator(t *testing.T) {
	t.Parallel()
	o := onpar.New()
	defer o.Run(t)

	o.BeforeEach(func(t *testing.T) TSA {
		return TSA{
			T: t,
			a: streamaggregator.New(),
		}
	})

	o.Group("adding consumers", func() {
		o.Group("producers are added", func() {
			o.Spec("it does not activate a producer without any consumers", func(t TSA) {
				ctxs := make(chan context.Context, 100)
				t.a.AddProducer("producer-1", streamaggregator.ProducerFunc(func(ctx context.Context, r interface{}, c chan<- interface{}) {
					ctxs <- ctx
				}))
				Expect(t, ctxs).To(HaveLen(0))
			})

			o.Spec("it activates a producer when added", func(t TSA) {
				clientCtx := context.Background()

				requests := make(chan interface{}, 100)
				ctxs := make(chan context.Context, 100)
				t.a.AddProducer("producer-1", streamaggregator.ProducerFunc(func(ctx context.Context, r interface{}, c chan<- interface{}) {
					ctxs <- ctx
					requests <- r

					// block to force each Produce to be on a different go-routine
					var wg sync.WaitGroup
					wg.Add(1)
					wg.Wait()
				}))

				Expect(t, ctxs).To(Always(HaveLen(0)))
				t.a.Consume(clientCtx, "some-request-stuff")

				Expect(t, ctxs).To(ViaPolling(HaveLen(1)))
				Expect(t, requests).To(Chain(Receive(), Equal("some-request-stuff")))
			})

			o.Spec("each producer writes to the returned channel", func(t TSA) {
				clientCtx := context.Background()

				t.a.AddProducer("producer-1", streamaggregator.ProducerFunc(func(ctx context.Context, r interface{}, c chan<- interface{}) {
					c <- 1
				}))

				t.a.AddProducer("producer-2", streamaggregator.ProducerFunc(func(ctx context.Context, r interface{}, c chan<- interface{}) {
					c <- 2
				}))

				data := t.a.Consume(clientCtx, "some-request-stuff", streamaggregator.WithConsumeChannelLength(5))

				Expect(t, data).To(ViaPolling(HaveLen(2)))
				Expect(t, data).To(Chain(Receive(), Or(Equal(1), Equal(2))))
				Expect(t, data).To(Chain(Receive(), Or(Equal(1), Equal(2))))
			})
		})

		o.Group("no producers yet", func() {
			o.Spec("it activates a producer when added", func(t TSA) {
				clientCtx := context.Background()
				t.a.Consume(clientCtx, "some-request-stuff")

				ctxs := make(chan context.Context, 100)
				requests := make(chan interface{}, 100)
				t.a.AddProducer("producer-1", streamaggregator.ProducerFunc(func(ctx context.Context, r interface{}, c chan<- interface{}) {
					ctxs <- ctx
					requests <- r

					// block to force each Produce to be on a different go-routine
					var wg sync.WaitGroup
					wg.Add(1)
					wg.Wait()

				}))

				Expect(t, ctxs).To(ViaPolling(HaveLen(1)))
				Expect(t, requests).To(Chain(Receive(), Equal("some-request-stuff")))
			})

			o.Spec("it does not activate a producer when added after context is cancelled", func(t TSA) {
				clientCtx, cancel := context.WithCancel(context.Background())
				t.a.Consume(clientCtx, "some-request-stuff")
				cancel()

				// Give time to allow the cancellation to be noticed
				time.Sleep(250 * time.Millisecond)

				ctxs := make(chan context.Context, 100)
				t.a.AddProducer("producer-1", streamaggregator.ProducerFunc(func(ctx context.Context, r interface{}, c chan<- interface{}) {
					ctxs <- ctx
				}))

				Expect(t, ctxs).To(Always(HaveLen(0)))
			})

			o.Spec("it closes the channel the context is cancelled and the producers exit", func(t TSA) {
				clientCtx, cancel := context.WithCancel(context.Background())
				data := t.a.Consume(clientCtx, "some-request-stuff")

				// Give time to allow the cancellation to be noticed
				time.Sleep(250 * time.Millisecond)

				t.a.AddProducer("producer-1", streamaggregator.ProducerFunc(func(ctx context.Context, r interface{}, c chan<- interface{}) {
					<-ctx.Done()
				}))

				t.a.AddProducer("producer-2", streamaggregator.ProducerFunc(func(ctx context.Context, r interface{}, c chan<- interface{}) {
					<-ctx.Done()
				}))

				cancel()
				Expect(t, data).To(ViaPolling(BeClosed()))
			})
		})
	})

	o.Group("removing producers", func() {
		o.Spec("it closes contexts to the corresponding producer", func(t TSA) {
			clientCtx := context.Background()
			t.a.Consume(clientCtx, "some-request-stuff")

			ctxs := make(chan context.Context, 100)
			t.a.AddProducer("producer-1", streamaggregator.ProducerFunc(func(ctx context.Context, r interface{}, c chan<- interface{}) {
				ctxs <- ctx
			}))

			var ctx context.Context
			Expect(t, ctxs).To(ViaPolling(HaveLen(1)))
			Expect(t, ctxs).To(Chain(Receive(), Fetch(&ctx)))

			t.a.RemoveProducer("producer-1")
			Expect(t, ctx.Done()).To(BeClosed())
		})

		o.Spec("it does not use the producer for new consumers", func(t TSA) {
			ctxs := make(chan context.Context, 100)
			t.a.AddProducer("producer-1", streamaggregator.ProducerFunc(func(ctx context.Context, r interface{}, c chan<- interface{}) {
				ctxs <- ctx
			}))

			t.a.RemoveProducer("producer-1")
			// Wait for the producer to be removed
			time.Sleep(250 * time.Millisecond)

			t.a.Consume(context.Background(), "some-request-stuff")
			Expect(t, ctxs).To(ViaPolling(Always(HaveLen(0))))
		})
	})
}
