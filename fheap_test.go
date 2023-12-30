package fheap

import (
	"flag"
	"math/rand"
	"testing"
)

var HeapSize = flag.Int("heapsize", 20, "size of arbitrary heap when testing")

func intMinHeap[V comparable]() *fheap[int, V] {
	return NewFHeap[int, V](func(x, y int) bool { return x < y })
}

func TestFHeap_pathological(t *testing.T) {
	var h *fheap[int, int]
	e := ErrNilHeap
	if _, err := h.Size(); err != e {
		t.Fatalf("[Size] expected %v, got %v", e, err)
	}
	if _, err := h.Contains(0); err != e {
		t.Fatalf("[Contains] expected %v, got %v", e, err)
	}
	if err := h.Push(1, 1); err != e {
		t.Fatalf("[Push] expected %v, got %v", e, err)
	}
	if _, err := h.Pop(); err != e {
		t.Fatalf("[Pop] expected %v, got %v", e, err)
	}
	if err := h.IncreasePriority(2, 7); err != e {
		t.Fatalf("[IncreasePriority] expected %v, got %v", e, err)
	}
}

func TestFHeapPush(t *testing.T) {
	h := intMinHeap[int]()
	N := *HeapSize
	for i := N; i > 0; i-- {
		p, v := i*2, i
		if err := h.Push(p, v); err != nil {
			t.Fatalf("Push(%d, %d) failed with %v", p, v, err)
		}
		if contains, err := h.Contains(N); err != nil {
			t.Fatal(err)
		} else if !contains {
			t.Fatalf("value=%v missing from heap", N)
		}
		if peek, err := h.Peek(); err != nil {
			t.Fatalf("[p=%d, v=%d] Peek failed with %v", p, v, err)
		} else if peek.Value != v {
			t.Fatalf("expected highest-priority element's value to be %d, got %d", v, peek.Value)
		}
	}
	if size, err := h.Size(); err != nil {
		t.Fatal(err)
	} else if size != N {
		t.Fatalf("expected size=1, got %d", size)
	}
}

func TestFHeapPop(t *testing.T) {
	h := intMinHeap[int]()
	N := *HeapSize
	// push one, pop one
	v := 34
	if err := h.Push(12, v); err != nil {
		t.Fatal(err)
	}
	if actual, err := h.Pop(); err != nil {
		t.Fatal(err)
	} else if actual != v {
		t.Fatalf("expected highest-priority value to be %d, got %d", v, actual)
	}
	if size, err := h.Size(); err != nil {
		t.Fatal(err)
	} else if size != 0 {
		t.Fatalf("expected empty heap, got size=%d", size)
	}
	// push all, then pop
	perm := rand.Perm(N)
	t.Logf("perm: %v", perm)
	for _, p := range perm {
		if err := h.Push(p, p); err != nil {
			t.Fatalf("Push(%[1]d, %[1]d) failed with %v", p, err)
		}
	}
	for expected := 0; expected < N; expected++ {
		if actual, err := h.Pop(); err != nil {
			t.Fatal(err)
		} else if actual != expected {
			t.Fatalf("[i=%[1]d] expected value=%[1]d, got %[2]d", expected, actual)
		}
	}
}
