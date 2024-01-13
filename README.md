# fibonacci-heap

## What is this?

Package `fheap` is an implementation of a Fibonacci heap. It mostly follows the pseudocode from [Introduction to Algorithms (3rd edition, 2009) by CLRS](https://dahlan.unimal.ac.id/files/ebooks/2009%20Introduction%20to%20Algorithms%20Third%20Ed.pdf).

Supported operations:

| Function                       | Effect                                       |
| :----------------------------- | :------------------------------------------- |
| `New[V, P](...) *fheap[V, P]`  | Creates an empty Fibonacci heap              |
| `Size() (int, error)`          | Return how many values are in the heap       |
| `Push(v, p) error`             | Add value `v` with priority `p` to heap      |
| `Pop() (V, error)`             | Pop the highest-priority value from the heap |
| `IncreasePriority(v, p) error` | Increase the priority of value `v` to `p`    |
| `Delete(v) error`              | Delete value `v` from the heap               |

Exported errors:

| Error                 | When                                                   |
| :-------------------- | :----------------------------------------------------- |
| `ErrNilHeap`          | The heap pointer is `nil`                              |
| `ErrEmptyHeap`        | The heap is empty                                      |
| `ErrReservedPriority` | The supplied priority is the sentinel highest-priority |

## Installation

`go get github.com/iyassou/fibonacci-heap`

## Example usage

Min-heap with `string` values and `int` priorities:

```go
package main

import (
	"log"
	"math"

	fheap "github.com/iyassou/fibonacci-heap"
)

func main() {
	higherThan := func(x, y int) bool { return x < y }
	sentinel := math.MinInt
	h := fheap.New[string, int](higherThan, sentinel)

	// Pushing to the heap.
	h.Push("low priority", 100)
	h.Push("meh priority", 50)
	h.Push("medium priority", 30)
	h.Push("high priority", 1)

	s, _ := h.Size()
	log.Print("Values in the heap: ", s) // 4

	// Popping and deleting values.
	v, _ := h.Pop()
	log.Print("Highest-priority value was: ", v) // "high priority"
	h.Delete("low priority")

	// Increasing priority.
	h.IncreasePriority("meh priority", -100)
	v, _ = h.Pop()
	log.Print("Highest-priority value was: ", v) // "meh priority"
}

```

## Notable implementation details

A Fibonacci heap consists of heap-ordered trees. A node in these trees has a priority and a value. The priority is what determines the node's ordering in the tree, and its value what we wish to order. Both the priority and value of a node are generic, with the following type constraints.

In order to support deletion and increasing the priority of a value in the heap within the theoretical amortised time complexities of `O(log n)` and `O(1)` respectively, this implementation makes use of a `map` associating each value to its node's pointer. For this reason, values in the heap must be `comparable` and unique.

Priorities on the other hand are of type `any`. `New` requires a user-defined priority comparison function. This function, `R`, must be a [connected binary relation](https://en.wikipedia.org/wiki/Connected_relation) on the priority's type `P`, i.e.

$$
\forall x,y \in P: x \neq y \implies xRy \text{ or } yRx
$$

This fact is used to check for priority equality, namely to keep the sentinel highest-priority value for internal use.

The sentinel highest-priority is reserved due to the implementation of `Delete`, which consists of first increasing the priority of the value to delete to this highest-priority, and then popping and discarding the value from the heap.
