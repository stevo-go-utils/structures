package structures

import "fmt"

type Node[T any] struct {
	Data T
	Next *Node[T]
}

// String returns the string representation of the node's data.
// time-complexity: O(1)
func (n *Node[T]) String() string {
	return fmt.Sprint(n.Data)
}
