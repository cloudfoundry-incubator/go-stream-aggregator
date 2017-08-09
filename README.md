# stream-aggregator
[![GoDoc][go-doc-badge]][go-doc][![travis][travis-badge]][travis]

A stream aggregator is used to aggregate data in the [fan-in](https://blog.golang.org/pipelines) scheme. 

### Dynamic Producers and Consumers
It manages producers and consumers to be added and removed freely.

When a new producer is added, the existing consumers will start to receive data from the producer. When a new consumer is added, all previoulsy added producers will start writing to the new consumer.

### Lifecycle Management
[Contexts](https://golang.org/pkg/context/) are used to manage the lifecycle of the consumer. Cancelling the consumer's context will indicated to the producers to drain (and stop any reconnect logic). Once each producer exits, the given channel is closed. Contexts were chosen due to their tight integration with [gRPC](https://github.com/grpc/grpc-go).

[go-doc-badge]:             https://godoc.org/github.com/apoydence/stream-aggregator?status.svg
[go-doc]:                   https://godoc.org/github.com/apoydence/stream-aggregator
[travis-badge]:             https://travis-ci.org/apoydence/stream-aggregator.svg?branch=master
[travis]:                   https://travis-ci.org/apoydence/stream-aggregator?branch=master
