package fheap

import (
	"errors"
	"fmt"
)

// fnode is a Fibonacci heap node, consisting of a:
//   - priority
//   - value
//   - bereavement flag
//   - parent, children, left, right node pointers
//   - degree
//
// fnode siblings are doubly-linked.
// Since fnodes are only used by fheaps, the implemented
// methods are, wlog, left-centric.
type fnode[V, P any] struct {
	Value                         V
	priority                      P
	bereaved                      bool
	parent, children, left, right *fnode[V, P]
	degree                        int
}

var errNilFnode = errors.New("nil node")
var errBarrenFnode = errors.New("barren node")

// newFnode creates a new fnode given a priority and a value.
// The node's parent and children pointers are nil, and its
// left and right pointers are set to itself.
func newFnode[V, P any](value V, priority P) *fnode[V, P] {
	f := &fnode[V, P]{priority: priority, Value: value}
	f.left = f
	f.right = f
	return f
}

// insertLeft inserts an fnode to the left of the current fnode.
func (fn *fnode[V, P]) insertLeft(other *fnode[V, P]) error {
	if fn == nil || other == nil {
		return errNilFnode
	}
	fn.left.right = other
	other.left = fn.left
	other.right = fn
	fn.left = other
	return nil
}

// insertChild inserts an fnode to the left of the current fnode's
// children pointer.
// The child fnode's parent pointer is updated, and the parent fnode's
// degree incremented.
func (fn *fnode[V, P]) insertChild(other *fnode[V, P]) (err error) {
	if fn == nil || other == nil {
		return errNilFnode
	}
	defer func() {
		if err == nil {
			other.parent = fn
			fn.degree++
		}
	}()
	if fn.children == nil {
		fn.children = other
	} else {
		err = fn.children.insertLeft(other)
	}
	return
}

// removeChild removes a child node from its parent.
func (fn *fnode[V, P]) removeChild(child *fnode[V, P]) error {
	if fn == nil {
		return errNilFnode
	}
	if fn.children == nil {
		return errBarrenFnode
	}
	if child.parent != fn {
		return fmt.Errorf("child %v is unrelated to node %v", child, fn)
	}
	fn.degree--
	if fn.children == child {
		fn.children = fn.children.right // wlog
	}
	if child.left == child.right && child == child.left {
		fn.children = nil
	} else {
		child.left.right = child.right
		child.right.left = child.left
	}
	return nil
}

// popChild pops this node's child's left sibling.
func (fn *fnode[V, P]) popChild() (*fnode[V, P], error) {
	if fn == nil {
		return nil, errNilFnode
	}
	if fn.children == nil {
		return nil, errBarrenFnode
	}
	child := fn.children.left
	err := fn.removeChild(child)
	return child, err
}
