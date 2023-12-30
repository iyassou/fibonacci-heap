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
type fnode[P, V any] struct {
	priority                      P
	Value                         V
	bereaved                      bool
	parent, children, left, right *fnode[P, V]
	degree                        int
}

var errNilFnode = errors.New("nil node")
var errBarrenFnode = errors.New("barren node")

// newFnode creates a new fnode given a priority and a value.
// The node's parent and children pointers are nil, and its
// left and right pointers are set to itself.
func newFnode[P, V any](priority P, value V) *fnode[P, V] {
	f := &fnode[P, V]{priority: priority, Value: value}
	f.left = f
	f.right = f
	return f
}

// insertLeft inserts an fnode to the left of the current fnode.
func (fn *fnode[P, V]) insertLeft(other *fnode[P, V]) error {
	if fn == nil || other == nil {
		return errNilFnode
	}
	fn.left = other
	other.right = fn
	return nil
}

// insertRight inserts an fnode to the right of the current fnode.
func (fn *fnode[P, V]) insertRight(other *fnode[P, V]) error {
	if fn == nil || other == nil {
		return errNilFnode
	}
	fn.right = other
	other.left = fn
	return nil
}

// insertChildLeft inserts an fnode to the left of the current fnode's
// children pointer, preserving the pre-existing sibling chain, if present.
// The child fnode's parent pointer is updated, and the parent fnode's
// degree incremented.
func (fn *fnode[P, V]) insertChildLeft(other *fnode[P, V]) (err error) {
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

// popLeftChild removes an fnode's left child from the parent fnode's
// children and its siblings, preserving the pre-existing sibling chain,
// if present. The parent's degree is updated, the child's parent pointer
// set to nil, and the child's left and right pointers to itself.
func (fn *fnode[P, V]) popLeftChild() (child *fnode[P, V], err error) {
	if fn == nil {
		err = errNilFnode
		return
	}
	if fn.children == nil {
		err = errBarrenFnode
		return
	}
	defer func() {
		if err == nil {
			child.parent = nil
			child.left = child
			child.right = child
			fn.degree--
		}
	}()
	child = fn.children.left
	if child == fn.children {
		fn.children = nil
		return
	}
	err = fn.children.insertLeft(child.left)
	return
}

// link "combines two item-disjoint [heap-ordered] trees into one.
// Given two trees with roots x and y, we link them by comparing
// the keys of items in x and y.
// If the item in x has the smaller key, we make y a child of x;
// otherwise, we make x a child of y."
// When calling this function, it is assumed that the function argument
// is to be made a child of the current fnode.
func (fn *fnode[P, V]) link(other *fnode[P, V]) error {
	if fn == nil || other == nil {
		return errNilFnode
	}
	if other.parent != nil {
		return fmt.Errorf("%[1]v@%[1]p", other)
		// if err := other.parent.removeChild(other); err != nil {
		// 	return err
		// }
	}
	return fn.insertChildLeft(other) // wlog
}
