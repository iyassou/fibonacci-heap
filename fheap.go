package fheap

import (
	"errors"
	"fmt"
)

// fheap is a Fibonacci heap, consisting of a:
//   - pointer to the highest-priority element
//   - priority comparison function
//   - map of values to fnodes
//
// The map is used to find the corresponding node when calling
// IncreasePriority or Remove. As such, this implementation
// requires heap values to be comparable and doesn't allow
// duplicate values.
type fheap[P any, V comparable] struct {
	prioritaire *fnode[P, V]
	higher      func(x, y P) bool
	values      map[V]*fnode[P, V]
}

var ErrNilHeap = errors.New("nil heap")
var ErrEmptyHeap = errors.New("empty heap")

// NewFHeap creates a new, empty Fibonacci heap.
func NewFHeap[P any, V comparable](higher func(x, y P) bool) *fheap[P, V] {
	return &fheap[P, V]{higher: higher, values: map[V]*fnode[P, V]{}}
}

// Size returns the number of elements in the heap.
func (fh *fheap[P, V]) Size() (int, error) {
	if fh == nil {
		return 0, ErrNilHeap
	}
	return len(fh.values), nil
}

// Contains determines if the supplied value is present in the heap.
func (fh *fheap[P, V]) Contains(value V) (bool, error) {
	if fh == nil {
		return false, ErrNilHeap
	}
	_, contains := fh.values[value]
	return contains, nil
}

// Push creates a new node with the supplied priority and value,
// and inserts it into the heap.
func (fh *fheap[P, V]) Push(priority P, value V) error {
	if fh == nil {
		return ErrNilHeap
	}
	if _, ok := fh.values[value]; ok {
		return fmt.Errorf("implementation does not allow duplicate values (value=%v)", value)
	}
	node := newFnode(priority, value)
	fh.values[value] = node
	if fh.prioritaire == nil {
		fh.prioritaire = node
		return nil
	}
	if err := fh.prioritaire.insertLeft(node); err != nil { // wlog
		return err
	}
	if fh.higher(priority, fh.prioritaire.priority) {
		fh.prioritaire = node
	}
	return nil
}

// Peek returns a pointer to the highest-priority element.
func (fh *fheap[P, V]) Peek() (*fnode[P, V], error) {
	if fh == nil {
		return nil, ErrNilHeap
	}
	if fh.prioritaire == nil {
		return nil, ErrEmptyHeap
	}
	return fh.prioritaire, nil
}

// Pop retrieves and removes the highest-priority element from
// the heap, consolidates the heap, and returns that element.
func (fh *fheap[P, V]) Pop() (value V, _ error) {
	if fh == nil {
		return value, ErrNilHeap
	}
	if fh.prioritaire == nil {
		return value, ErrEmptyHeap
	}
	// TODO
	return value, nil
}

// IncreasePriority increases the priority of a given value's node,
// if the value is present in the heap.
func (fh *fheap[P, V]) IncreasePriority(value V, priority P) error {
	if fh == nil {
		return ErrNilHeap
	}
	if fh.prioritaire == nil {
		return ErrEmptyHeap
	}
	// TODO
	return nil
}

// meld combines the root lists of two heaps into a single list, and sets
// the highest-priority node of the new heap to be the more prioritised
// highest-priority node of the two original heaps.
// The calling heap is overwritten with the melded heap.
func (fh *fheap[P, V]) meld(other *fheap[P, V]) error {
	if fh == nil || other == nil {
		return ErrNilHeap
	}
	if other.prioritaire == nil {
		return nil
	}
	if fh.prioritaire == nil {
		*fh = *other
		return nil
	}
	// check for duplicate values
	a := fh.values
	b := other.values
	if len(fh.values) > len(other.values) {
		a, b = b, a
	}
	for v := range a {
		if _, ok := b[v]; ok {
			return fmt.Errorf("value=%v present in both heaps", v)
		}
	}
	// merge value-node maps
	for v, node := range a {
		b[v] = node
	}
	fh.values = b
	// combine root lists
	if err := fh.prioritaire.insertLeft(other.prioritaire); err != nil {
		return err
	}
	if fh.higher(other.prioritaire.priority, fh.prioritaire.priority) {
		fh.prioritaire = other.prioritaire
	}
	return nil
}
