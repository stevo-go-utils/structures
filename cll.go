package structures

import (
	"fmt"
	"strings"
)

type CircularLinkedList[T comparable] struct {
	tail *Node[T]
	Size int
}

// New constructs and returns an empty circularly linked-list.
// time-complexity: O(1)
func NewCircularLinkedList[T comparable]() CircularLinkedList[T] {
	return CircularLinkedList[T]{}
}

// IsEmpty returns true if the linked-list doesn't contain any nodes.
// time-complexity: O(1)
func (c *CircularLinkedList[T]) IsEmpty() bool {
	return c.Size == 0
}

// First returns the first element of the list. It returns false if the list is empty.
// time-complexity: O(1)
func (c *CircularLinkedList[T]) First() (data T, ok bool) {
	if c.IsEmpty() {
		return
	}
	return c.tail.Next.Data, true
}

// Last returns the last element of the list. It returns false if the list is empty.
// time-complexity: O(1)
func (c *CircularLinkedList[T]) Last() (data T, ok bool) {
	if c.IsEmpty() {
		return
	}
	return c.tail.Data, true
}

// Rotate rotates the list. It moves the first element to the end.
// time-complexity: O(1)
func (c *CircularLinkedList[T]) Rotate() {
	if c.tail != nil {
		c.tail = c.tail.Next
	}
}

// AddFirst adds a new node to the beginning of the list.
// time-complexity: O(1)
func (c *CircularLinkedList[T]) AddFirst(data T) {
	if c.IsEmpty() {
		c.tail = &Node[T]{Data: data}
		c.tail.Next = c.tail
	} else {
		n := Node[T]{Data: data, Next: c.tail.Next}
		c.tail.Next = &n
	}
	c.Size++
}

// AddLast adds a new node to the end of the list.
// time-complexity: O(1)
func (c *CircularLinkedList[T]) AddLast(data T) {
	c.AddFirst(data)
	c.tail = c.tail.Next
}

// RemoveFirst removes and returns the first element of the list. It returns false if the list is empty.
// time-complexity: O(1)
func (c *CircularLinkedList[T]) RemoveFirst() (val T, ok bool) {
	if c.IsEmpty() {
		return
	}

	head := c.tail.Next

	val = head.Data

	if head == c.tail {
		c.tail = nil
	} else {
		c.tail.Next = head.Next
	}

	c.Size--

	return val, true
}

// RemoveLast removes and returns the last element of the list. It returns false if the list empty.
// time-complexity: O(n)
func (c *CircularLinkedList[T]) RemoveLast() (val T, ok bool) {
	if c.IsEmpty() {
		return
	}

	val = c.tail.Data

	current := c.tail.Next
	for ; current.Next != c.tail; current = current.Next {
	}

	current.Next = c.tail.Next
	c.tail = current

	c.Size--

	return val, true
}

func (c *CircularLinkedList[T]) Remove(val T) {
	if c.IsEmpty() {
		return
	}

	var prev *Node[T]
	current := c.tail.Next
	for current != c.tail {
		if current.Data == val {
			if prev == nil {
				c.tail.Next = current.Next
			} else {
				prev.Next = current.Next
			}
			c.Size--
			return
		}
		prev = current
		current = current.Next
	}

	if current.Data == val {
		if prev == nil {
			c.tail = nil
		} else {
			prev.Next = current.Next
			c.tail = prev
		}
		c.Size--
	}
}

// String returns the string representation of the list.
// time-complexity: O(n)
func (c *CircularLinkedList[T]) Vals() (vals []T) {
	if c.IsEmpty() {
		return
	}

	for current := c.tail.Next; current != c.tail; current = current.Next {
		vals = append(vals, current.Data)
	}
	vals = append(vals, c.tail.Data)
	return
}

func (c *CircularLinkedList[T]) Add(data T) {
	c.AddFirst(data)
}

// String returns the string representation of the list.
// time-complexity: O(n)
func (c *CircularLinkedList[T]) String() string {
	if c.IsEmpty() {
		return "[ ]"
	}

	var b strings.Builder
	b.WriteString("[ ")

	for current := c.tail.Next; current != c.tail; current = current.Next {
		b.WriteString(fmt.Sprint(current.Data))
		b.WriteString(" ")
	}

	b.WriteString(fmt.Sprint(c.tail.Data))
	b.WriteString(" ]")

	return b.String()
}
