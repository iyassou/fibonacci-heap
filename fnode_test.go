package fheap

import (
	"flag"
	"testing"
)

var ListSize = flag.Int("listsize", 100, "size of arbitrary list when testing")

func TestFNodeInsertLeft(t *testing.T) {
	f := newFnode(0, 0)
	N := *ListSize
	for i := 1; i <= N; i++ {
		if err := f.insertLeft(newFnode(i, i)); err != nil {
			t.Fatal(err)
		}
	}
	// walk left to right
	counter := 0
	for iter := f; ; iter = iter.right {
		val := iter.Value
		if val != counter {
			t.Fatalf("expected %d, got %d", counter, val)
		}
		counter++
		if iter.right == f {
			break
		}
	}
	// walk right to left
	counter = N
	for iter := f.left; ; iter = iter.left {
		val := iter.Value
		if val != counter {
			t.Fatalf("expected %d, got %d", counter, val)
		}
		counter--
		if iter.left == f.left {
			break
		}
	}
}

func TestFNodeInsertChild(t *testing.T) {
	f := newFnode(0, 0)
	N := *ListSize
	for i := 1; i <= N; i++ {
		if err := f.insertChild(newFnode(i, i)); err != nil {
			t.Fatal(err)
		}
		deg := f.degree
		if deg != i {
			t.Fatalf("expected degree=%d, got %d", i, deg)
		}
	}
	// walk children left to right
	counter := 1
	for iter := f.children; ; iter = iter.right {
		val := iter.Value
		if val != counter {
			t.Fatalf("expected %d, got %d", counter, val)
		}
		counter++
		if iter.right == f.children {
			break
		}
	}
	// walk children right to left
	counter = N
	for iter := f.children.left; ; iter = iter.left {
		val := iter.Value
		if val != counter {
			t.Fatalf("expected %d, got %d", counter, val)
		}
		counter--
		if iter.left == f.children.left {
			break
		}
	}
}

func TestFNodePopChild(t *testing.T) {
	f := newFnode(0, 0)
	N := *ListSize
	for i := 1; i <= N; i++ {
		if err := f.insertChild(newFnode(i, i)); err != nil {
			t.Fatal(err)
		}
	}
	for counter := N; ; counter-- {
		child, err := f.popChild()
		if err != nil {
			if err == errBarrenFnode {
				break
			}
			t.Fatal(err)
		}
		val := child.Value
		if val != counter {
			t.Fatalf("expected %d, got %d", counter, val)
		}
	}
}
