package streamaggregator

import (
  "github.com/apoydence/pubsub"
  "fmt"
)
type notificationTrav struct{}
 func NewNotificationTrav()notificationTrav{ return notificationTrav{} }

func (s notificationTrav) Traverse(data interface{}, currentPath []string) pubsub.Paths {
	return s._added(data, currentPath)
}

	func (s notificationTrav) done(data interface{}, currentPath []string) pubsub.Paths {
	return pubsub.FlatPaths(nil)
}

func (s notificationTrav) _added(data interface{}, currentPath []string) pubsub.Paths {
	
  return pubsub.NewPathsWithTraverser([]string{"", fmt.Sprintf("%v", data.(producerNotification).added)}, pubsub.TreeTraverserFunc(s._key))
}

func (s notificationTrav) _key(data interface{}, currentPath []string) pubsub.Paths {
	
  return pubsub.NewPathsWithTraverser([]string{"", fmt.Sprintf("%v", data.(producerNotification).key)}, pubsub.TreeTraverserFunc(s.done))
}

type producerNotificationFilter struct{
added *bool
key *string

}

func (g notificationTrav) CreatePath(f *producerNotificationFilter) []string {
if f == nil {
	return nil
}
var path []string




var count int
if count > 1 {
	panic("Only one field can be set")
}


if f.added != nil {
	path = append(path, fmt.Sprintf("%v", *f.added))
}else{
	path = append(path, "")
}

if f.key != nil {
	path = append(path, fmt.Sprintf("%v", *f.key))
}else{
	path = append(path, "")
}





return path
}
