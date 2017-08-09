package traverser

import (
  "github.com/apoydence/pubsub"
  "fmt"
)
type Traverser struct{}
 func NewTraverser()Traverser{ return Traverser{} }

func (s Traverser) Traverse(data interface{}, currentPath []string) pubsub.Paths {
	return s._Added(data, currentPath)
}

	func (s Traverser) done(data interface{}, currentPath []string) pubsub.Paths {
	return pubsub.FlatPaths(nil)
}

func (s Traverser) _Added(data interface{}, currentPath []string) pubsub.Paths {
	
  return pubsub.NewPathsWithTraverser([]string{"", fmt.Sprintf("%v", data.(Notification).Added)}, pubsub.TreeTraverserFunc(s._Key))
}

func (s Traverser) _Key(data interface{}, currentPath []string) pubsub.Paths {
	
  return pubsub.NewPathsWithTraverser([]string{"", fmt.Sprintf("%v", data.(Notification).Key)}, pubsub.TreeTraverserFunc(s.done))
}

type NotificationFilter struct{
Added *bool
Key *string

}

func (g Traverser) CreatePath(f *NotificationFilter) []string {
if f == nil {
	return nil
}
var path []string




var count int
if count > 1 {
	panic("Only one field can be set")
}


if f.Added != nil {
	path = append(path, fmt.Sprintf("%v", *f.Added))
}else{
	path = append(path, "")
}

if f.Key != nil {
	path = append(path, fmt.Sprintf("%v", *f.Key))
}else{
	path = append(path, "")
}





return path
}
