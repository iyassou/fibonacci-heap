package fheap

import (
	"errors"
	"flag"
	"fmt"
	"math"
	"math/rand"
	"testing"
)

var HeapSize = flag.Int("heapsize", 100, "size of arbitrary heap when testing")

func isOrderedHeap[V comparable, P any](h *fheap[V, P], n *fnode[V, P]) error {
	if n == nil {
		return errNilFnode
	}
	prefix := fmt.Sprintf("%[1]v @ %[1]p", n)
	numChildren := 0
	if n.children != nil {
		for child := n.children; ; child = child.right {
			if h.higherThan(child.priority, n.priority) {
				return fmt.Errorf("%s: parent (v=%v) priority %v lower than child's (v=%v) %v", prefix, n.Value, n.priority, child.Value, child.priority)
			}
			if child.parent != n {
				return fmt.Errorf("%s: child's parent = %v", prefix, child.parent)
			}
			if err := isOrderedHeap(h, child); err != nil {
				return err
			}
			numChildren++
			if child.right == n.children {
				break
			}
		}
	}
	if numChildren == n.degree {
		return nil
	}
	return fmt.Errorf("%s: counted %d children, but degree=%d", prefix, numChildren, n.degree)
}

func isFibonacciHeap[V comparable, P any](h *fheap[V, P]) error {
	if h == nil {
		return nil
	}
	prefix := fmt.Sprintf("%[1]v @ %[1]p", h)
	if h.prioritaire == nil {
		if len(h.values) == 0 {
			return nil
		}
		return fmt.Errorf("%s: prioritaire=nil but %d nodes", prefix, len(h.values))
	}
	for root := h.prioritaire; ; root = root.right {
		if root.bereaved {
			return fmt.Errorf("%s: bereaved", prefix)
		}
		if err := isOrderedHeap(h, root); err != nil {
			return err
		}
		if root.right == h.prioritaire {
			break
		}
	}
	return nil
}

func intMinHeap[V comparable]() *fheap[V, int] {
	return New[V, int](func(x, y int) bool { return x < y }, math.MinInt)
}

func Push[V comparable, P any](h *fheap[V, P], v V, p P, name string) error {
	if err := h.Push(v, p); err != nil {
		return fmt.Errorf("[%s] Push(p=%v, v=%v) failed with %w", name, p, v, err)
	}
	if err := isFibonacciHeap(h); err != nil {
		return fmt.Errorf("[%s] Push(p=%v, v=%v), err=%v", name, p, v, err)
	}
	return nil
}

func Pop[V comparable, P any](h *fheap[V, P], name string) (V, error) {
	v, err := h.Pop()
	if err != nil {
		return v, fmt.Errorf("[%s] Pop() failed with %w", name, err)
	}
	if err := isFibonacciHeap(h); err != nil {
		return v, fmt.Errorf("[%s] Pop(), err=%v", name, err)
	}
	return v, nil
}

func IncreasePriority[V comparable, P any](h *fheap[V, P], v V, p P, name string) error {
	if err := h.IncreasePriority(v, p); err != nil {
		return fmt.Errorf("[%s] IncreasePriority(v=%v, p=%v) failed with %w", name, v, p, err)
	}
	if err := isFibonacciHeap(h); err != nil {
		return fmt.Errorf("[%s] IncreasePriority(v=%v, p=%v), err=%v", name, v, p, err)
	}
	return nil
}

func Delete[V comparable, P any](h *fheap[V, P], v V, name string) error {
	if err := h.Delete(v); err != nil {
		return fmt.Errorf("[%s] Delete(v=%v) failed with %w", name, v, err)
	}
	if err := isFibonacciHeap(h); err != nil {
		return fmt.Errorf("[%s] Delete(v=%v), err=%v", name, v, err)
	}
	return nil
}

func TestFHeap_NilHeap(t *testing.T) {
	var h *fheap[int, int]
	e := ErrNilHeap
	msg := fmt.Sprintf("[%s] expected %v, got %v", "%s", e, "%v")
	if _, err := h.Size(); err != e {
		t.Fatalf(msg, "Size", err)
	}
	if err := h.Push(1, 1); err != e {
		t.Fatalf(msg, "Push", err)
	}
	if _, err := h.Pop(); err != e {
		t.Fatalf(msg, "Pop", err)
	}
	if err := h.IncreasePriority(2, 7); err != e {
		t.Fatalf(msg, "IncreasePriority", err)
	}
	if err := h.Delete(12); err != e {
		t.Fatalf(msg, "Delete", err)
	}
}

func TestFHeap_EmptyHeap(t *testing.T) {
	h := intMinHeap[string]()
	if s, err := h.Size(); err != nil {
		t.Fatalf("[Size] failed with %v", err)
	} else if s != 0 {
		t.Fatalf("[Size] expected size=0, got %d", s)
	}
	if _, err := h.Pop(); err != ErrEmptyHeap {
		t.Fatalf("[Pop] expected ErrEmptyHeap, got err=%v", err)
	}
	value := "3-2-1 girls wanna have fun, if the man don't dance he's done"
	if err := h.IncreasePriority(value, 12); err != ErrEmptyHeap {
		t.Fatalf("[IncreasePriority] expected ErrEmptyHeap, got err=%v", err)
	}
	if err := h.Delete("wesh gros"); err != ErrEmptyHeap {
		t.Fatalf("[Delete] expected ErrEmptyHeap, got err=%v", err)
	}
}

func TestFHeapPush(t *testing.T) {
	h := intMinHeap[int]()
	N := *HeapSize
	for i := N; i > 0; i-- {
		if err := Push(h, i, i*2, t.Name()); err != nil {
			t.Fatal(err)
		}
	}
	if size, err := h.Size(); err != nil {
		t.Fatal(err)
	} else if size != N {
		t.Fatalf("expected size=1, got %d", size)
	}
	expected := fmt.Sprintf("duplicate value=%v", N)
	err := Push(h, N, 123123123, t.Name())
	if err == nil {
		t.Fatal("expected duplicate value error")
	}
	if actual := errors.Unwrap(err).Error(); actual != expected {
		t.Fatalf("expected %q, got %q", expected, actual)
	}
}

func TestFHeapPop_OneInOneOut(t *testing.T) {
	h := intMinHeap[int]()
	v := 34
	if err := Push(h, v, 12, t.Name()); err != nil {
		t.Fatal(err)
	}
	if actual, err := Pop(h, t.Name()); err != nil {
		t.Fatal(err)
	} else if actual != v {
		t.Fatalf("expected highest-priority value to be %d, got %d", v, actual)
	}
	if size, err := h.Size(); err != nil {
		t.Fatal(err)
	} else if size != 0 {
		t.Fatalf("expected empty heap, got size=%d", size)
	}
}

func TestFHeapPop_RandomPermutation(t *testing.T) {
	h := intMinHeap[int]()
	N := *HeapSize
	perm := rand.Perm(N)
	for _, p := range perm {
		if err := Push(h, p, p, t.Name()); err != nil {
			t.Fatal(err)
		}
	}
	for expected := 0; expected < N; expected++ {
		if actual, err := Pop(h, t.Name()); err != nil {
			t.Fatal(err)
		} else if actual != expected {
			t.Fatalf("[i=%[1]d] expected value=%[1]d, got %[2]d", expected, actual)
		}
	}
}

func TestFHeapIncreasePriority(t *testing.T) {
	type testcase struct {
		name           string
		values         []int
		priorities     []int
		pops           int // to modify heap structure
		incValues      []int
		incPriorities  []int
		expectedErrors []error
	}
	testcases := []testcase{
		{
			"1 node, 0 pops, increase priority",
			[]int{13},
			[]int{1_000_000},
			0,
			[]int{13},
			[]int{999_999},
			[]error{nil},
		},
		{
			"1 node, 0 pops, decrease priority",
			[]int{9},
			[]int{1},
			0,
			[]int{9},
			[]int{2},
			[]error{errors.New("old priority 1 is higher than new 2")},
		},
		{
			"7 nodes, 1 pop, increase lowest priority to highest",
			[]int{1, 2, 3, 4, 5, 6, 7},
			[]int{1, 2, 3, 4, 5, 6, 7},
			1,
			[]int{7},
			[]int{0},
			[]error{nil},
		},
		{
			"10 nodes, 2 pops, increase lowest priority to mid",
			[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			2,
			[]int{10},
			[]int{5},
			[]error{nil},
		},
	}
	for _, tc := range testcases {
		h := intMinHeap[int]()
		prefix := fmt.Sprintf("[%s | %s]", t.Name(), tc.name)
		vs := tc.values
		ps := tc.priorities
		if len(vs) != len(ps) {
			t.Fatalf("%s bruh", prefix)
		}
		for i, v := range vs {
			if err := Push(h, v, ps[i], prefix); err != nil {
				t.Fatal(err)
			}
		}
		for ; tc.pops > 0; tc.pops-- {
			if _, err := Pop(h, prefix); err != nil {
				t.Fatal(err)
			}
		}
		vs = tc.incValues
		newPs := tc.incPriorities
		errs := tc.expectedErrors
		if len(vs) != len(newPs) || len(vs) != len(errs) {
			t.Fatalf("%s reuf", prefix)
		}
		for i, v := range vs {
			np := newPs[i]
			expected := errs[i]
			err := IncreasePriority(h, v, np, prefix)
			bad := expected == nil && err != nil ||
				expected != nil && err == nil ||
				expected != nil && err != nil &&
					expected.Error() != errors.Unwrap(err).Error()
			if bad {
				t.Fatalf("%s expected err=%q for v=%v, newP=%v, got %q",
					prefix, expected, v, np, err)
			}
		}
	}
}

func TestFHeapDelete(t *testing.T) {
	type testcase struct {
		name           string
		values         []int
		priorities     []int
		pops           int // to modify heap structure
		deleteValues   []int
		expectedErrors []error
	}
	testcases := []testcase{
		{
			"1 node, 0 pops, delete highest-priority",
			[]int{99},
			[]int{8},
			0,
			[]int{99},
			[]error{nil},
		},
		{
			"3 nodes, 1 pop, delete highest-priority",
			[]int{1, 2, 3},
			[]int{1, 2, 3},
			1,
			[]int{2},
			[]error{nil},
		},
		{
			"3 nodes, 1 pop, delete lower priority",
			[]int{1, 2, 3},
			[]int{1, 2, 3},
			1,
			[]int{3},
			[]error{nil},
		},
		{
			"3 nodes, 1 pop, delete high-low-fromEmpty",
			[]int{1, 2, 3},
			[]int{1, 2, 3},
			1,
			[]int{2, 3, 4},
			[]error{nil, nil, ErrEmptyHeap},
		},
		{
			"10 nodes, 3 pops, delete high-mid-low",
			[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			3,
			[]int{4, 7, 10},
			[]error{nil, nil, nil},
		},
	}
	for _, tc := range testcases {
		h := intMinHeap[int]()
		prefix := fmt.Sprintf("[%s | %s]", t.Name(), tc.name)
		vs := tc.values
		ps := tc.priorities
		if len(vs) != len(ps) {
			t.Fatalf("%s bruh", prefix)
		}
		for i, v := range vs {
			if err := Push(h, v, ps[i], prefix); err != nil {
				t.Fatal(err)
			}
		}
		for ; tc.pops > 0; tc.pops-- {
			if _, err := Pop(h, prefix); err != nil {
				t.Fatal(err)
			}
		}
		vs = tc.deleteValues
		errs := tc.expectedErrors
		if len(vs) != len(errs) {
			t.Fatalf("%s reuf", prefix)
		}
		for i, v := range vs {
			expected := errs[i]
			err := Delete(h, v, prefix)
			bad := expected == nil && err != nil ||
				expected != nil && err == nil ||
				expected != nil && err != nil &&
					expected.Error() != errors.Unwrap(err).Error()
			if bad {
				t.Fatalf("%s expected err=%q for v=%v, got %q",
					prefix, expected, v, err)
			}
		}
	}
}
