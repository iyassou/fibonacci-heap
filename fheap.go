package fheap

import (
	"errors"
	"fmt"
	"math"
)

// fheap is a Fibonacci heap, consisting of a:
//   - pointer to the highest-priority element
//   - map of values to fnodes
//   - priority comparison function
//   - the highest priority an element can have
//
// The map is used by `IncreasePriority` to find a value's corresponding node.
// To this end, this implementation requires heap values to be comparable
// and doesn't allow duplicate values.
// `higherThan` determines if the first priority is higher than the second.
// `Delete` requires `higherThan` to be a connected relation on the priority
// set, i.e. for priorities x, y, if x != y then either x is higher than y
// or y is higher than x (https://en.wikipedia.org/wiki/Connected_relation).
// `highestPriority` is the highest possible priority a value can have. It will
// be reserved for internal use by `Delete`.
type fheap[V comparable, P any] struct {
	prioritaire     *fnode[V, P]
	values          map[V]*fnode[V, P]
	higherThan      func(x, y P) bool
	highestPriority P
}

var ErrNilHeap = errors.New("nil heap")
var ErrEmptyHeap = errors.New("empty heap")
var ErrReservedPriority = errors.New("highest priority is reserved for internal use")

// New creates an empty Fibonacci heap.
func New[V comparable, P any](higherThan func(x, y P) bool, highestPriority P) *fheap[V, P] {
	return &fheap[V, P]{
		values:          map[V]*fnode[V, P]{},
		higherThan:      higherThan,
		highestPriority: highestPriority}
}

// Size returns the number of elements in the heap.
func (fh *fheap[V, P]) Size() (int, error) {
	if fh == nil {
		return 0, ErrNilHeap
	}
	return len(fh.values), nil
}

// Push inserts a given value with the supplied priority into the heap.
func (fh *fheap[V, P]) Push(value V, priority P) error {
	if fh == nil {
		return ErrNilHeap
	}
	if fh.prioritiesEqual(priority, fh.highestPriority) {
		return ErrReservedPriority
	}
	if _, ok := fh.values[value]; ok {
		return fmt.Errorf("duplicate value=%v", value)
	}
	node := newFnode(value, priority)
	fh.values[value] = node
	if fh.prioritaire == nil {
		fh.prioritaire = node
		return nil
	}
	if err := fh.prioritaire.insertLeft(node); err != nil {
		return err
	}
	if fh.higherThan(priority, fh.prioritaire.priority) {
		fh.prioritaire = node
	}
	return nil
}

// Pop removes and returns the highest-priority element from the heap
// after consolidating the heap.
func (fh *fheap[V, P]) Pop() (value V, err error) {
	defer func() {
		if err == nil {
			delete(fh.values, value)
		}
	}()
	if fh == nil {
		return value, ErrNilHeap
	}
	if fh.prioritaire == nil {
		return value, ErrEmptyHeap
	}
	value = fh.prioritaire.Value
	// foster out prioritaire's children
	var child *fnode[V, P]
	for {
		child, err = fh.prioritaire.popChild()
		if err != nil {
			if err == errBarrenFnode {
				err = nil
				break
			}
			return
		}
		child.parent = nil
		child.left = child
		child.right = child
		child.bereaved = false
		if err = fh.prioritaire.insertLeft(child); err != nil {
			return
		}
	}
	// remove prioritaire from the heap's root list
	if fh.prioritaire.left == fh.prioritaire.right && fh.prioritaire.left == fh.prioritaire {
		fh.prioritaire = nil
	} else {
		fh.prioritaire.left.right = fh.prioritaire.right
		fh.prioritaire.right.left = fh.prioritaire.left
		fh.prioritaire = fh.prioritaire.right
		err = fh.consolidate()
	}
	return
}

// IncreasePriority increases a value's priority in the heap, if present.
// An error is returned if the priority is the heap's `highestPriority`.
func (fh *fheap[V, P]) IncreasePriority(value V, priority P) error {
	if fh == nil {
		return ErrNilHeap
	}
	if fh.prioritaire == nil {
		return ErrEmptyHeap
	}
	if fh.prioritiesEqual(priority, fh.highestPriority) {
		return ErrReservedPriority
	}
	return fh.increasePriority(value, priority)
}

// Delete deletes a value from the heap, if present. Operation consists
// of increasing its priority to the highest priority before popping the
// highest-priority element (itself).
func (fh *fheap[V, P]) Delete(value V) error {
	if fh == nil {
		return ErrNilHeap
	}
	if fh.prioritaire == nil {
		return ErrEmptyHeap
	}
	if err := fh.increasePriority(value, fh.highestPriority); err != nil {
		return err
	}
	_, err := fh.Pop()
	return err
}

// consolidate reduces the number of trees in the heap.
func (fh *fheap[V, P]) consolidate() error {
	D := int(math.Ceil(math.Log2(float64(len(fh.values)))))
	A := make([]*fnode[V, P], D+1)
	end := fh.prioritaire.left
	for w := fh.prioritaire; ; {
		next := w.right
		x := w
		d := x.degree
		for A[d] != nil {
			y := A[d]
			if fh.higherThan(y.priority, x.priority) {
				x, y = y, x
			}
			if err := fh.link(y, x); err != nil {
				return err
			}
			A[d] = nil
			d++
		}
		A[d] = x
		if w == end {
			break
		}
		w = next
	}
	// find the new minimum
	fh.prioritaire = nil
	for _, root := range A {
		if root == nil {
			continue
		}
		root.left.right = root.right
		root.right.left = root.left
		root.left = root
		root.right = root
		if fh.prioritaire == nil {
			fh.prioritaire = root
			continue
		}
		if err := fh.prioritaire.insertLeft(root); err != nil {
			return err
		}
		if fh.higherThan(root.priority, fh.prioritaire.priority) {
			fh.prioritaire = root
		}
	}
	return nil
}

// link removes y from the root list, and makes y a child of x.
func (fh *fheap[V, P]) link(y, x *fnode[V, P]) error {
	// remove y from the root list of H
	y.left.right = y.right
	y.right.left = y.left
	y.left = y
	y.right = y
	// make y a child of x
	if err := x.insertChild(y); err != nil {
		return err
	}
	// unmark y
	y.bereaved = false
	return nil
}

// prioritiesEqual determines if two priorities are equal.
func (fh *fheap[V, P]) prioritiesEqual(a, b P) bool {
	// R := `higherThan` is a connected binary relation, so
	//					x != y 	=>	xRy || yRx
	// hence
	//			!(xRy || yRx)	=>	!(x != y)
	//		<=>	!xRy && !yRx	=>	x = y
	return !fh.higherThan(a, b) && !fh.higherThan(b, a)
}

// increasePriority increases the priority of a value in the heap, if present.
// For internal use, as it allows setting the priority to `highestPriority`.
func (fh *fheap[V, P]) increasePriority(value V, priority P) error {
	x, ok := fh.values[value]
	if !ok {
		return fmt.Errorf("value %v missing from heap", value)
	}
	if fh.higherThan(x.priority, priority) {
		return fmt.Errorf("old priority %v is higher than new %v", x.priority, priority)
	}
	x.priority = priority
	y := x.parent
	if y != nil && fh.higherThan(x.priority, y.priority) {
		if err := fh.cut(x, y); err != nil {
			return err
		}
		if err := fh.cascadingCut(y); err != nil {
			return err
		}
	}
	if fh.higherThan(x.priority, fh.prioritaire.priority) {
		fh.prioritaire = x
	}
	return nil
}

// cut severs the link between x and its parent y, and turns x into a root.
func (fh *fheap[V, P]) cut(x, y *fnode[V, P]) error {
	if err := y.removeChild(x); err != nil {
		return err
	}
	x.left = x
	x.right = x
	x.parent = nil
	x.bereaved = false
	return fh.prioritaire.insertLeft(x)
}

// cascadingCut handles the ancestral consequences of cutting a node.
func (fh *fheap[V, P]) cascadingCut(y *fnode[V, P]) error {
	z := y.parent
	if z != nil {
		if !y.bereaved {
			y.bereaved = true
		} else {
			if err := fh.cut(y, z); err != nil {
				return err
			}
			if err := fh.cascadingCut(z); err != nil {
				return err
			}
		}
	}
	return nil
}
