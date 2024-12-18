// Package heap provides a generic implementation of a heap. A heap is a tree
// with the property that each node is the minimum-valued node in its subtree.
//
// The minimum element in the tree is the root, at index 0.
//
// A heap is a common way to implement a priority queue. To build a priority
// queue, create a Heap of the type of elements it will hold and specify a
// "less" function that orders the elements by priority. Use `Push` to add
// items and `Pop` to remove the item with the greatest precedence.
package heap

// Heap implements a binary heap.
type Heap[T any] struct {
	data []T
	less func(a, b T) bool
}

// New returns a new heap with the given less function. The less function
// returns whether 'a' is less than 'b'.
func New[T any](less func(a, b T) bool) *Heap[T] {
	return &Heap[T]{
		less: less,
	}
}

// NewFrom returns a new heap with the given less function and initial data.
func NewFrom[T any](less func(a, b T) bool, data ...T) *Heap[T] {
	n := len(data)
	h := &Heap[T]{
		less: less,
		data: data,
	}
	for i := n/2 - 1; i >= 0; i-- {
		h.down(i)
	}
	return h
}

// Len returns the number of elements in the heap.
func (h *Heap[T]) Len() int {
	return len(h.data)
}

// Push pushes the given element onto the heap.
func (h *Heap[T]) Push(x T) {
	h.data = append(h.data, x)
	h.up(len(h.data) - 1)
}

// Pop removes and returns the minimum element from the heap. If the heap is
// empty, it returns zero value and false.
func (h *Heap[T]) Pop() T {
	if len(h.data) == 0 {
		panic("heap: Pop called on empty heap")
	}

	var zero T
	x := h.data[0]
	n := len(h.data) - 1
	h.data[0] = h.data[n]
	h.data[n] = zero
	h.data = h.data[:n]
	h.down(0)

	return x
}

// Peek returns the minimum element from the heap without removing it. if the
// heap is empty, it returns zero value and false.
func (h *Heap[T]) Peek() T {
	if len(h.data) == 0 {
		panic("heap: Peek called on empty heap")
	}

	return h.data[0]
}

// Remove removes and returns the element at index i from the heap.
// The complexity is O(log n) where n = h.Len().
func (h *Heap[T]) Remove(i int) T {
	n := len(h.data) - 1
	if i < 0 || i > n {
		panic("heap: Remove index out of range")
	}
	if i == 0 {
		return h.Pop()
	}

	var zero T
	x := h.data[i]
	if n != i {
		h.data[i] = h.data[n]
		h.data[n] = zero
		h.data = h.data[:n]
		if !h.down(i) {
			h.up(i)
		}
	} else {
		h.data[n] = zero
		h.data = h.data[:n]
	}
	return x
}

// At returns the element at index i from the heap.
func (h *Heap[T]) At(i int) T {
	if i < 0 || i >= len(h.data) {
		panic("heap: At index out of range")
	}
	return h.data[i]
}

// Set replaces the element at index i in the heap and then calls [Fix] to
// restore the heap condition.
func (h *Heap[T]) Set(i int, x T) {
	if i < 0 || i >= len(h.data) {
		panic("heap: Set index out of range")
	}
	h.data[i] = x
	h.Fix(i)
}

// Fix re-establishes the heap ordering after the element at index i has changed its value.
// Changing the value of the element at index i and then calling Fix is equivalent to,
// but less expensive than, calling [Remove](i) followed by a Push of the new value.
// The complexity is O(log n) where n = h.Len().
func (h *Heap[T]) Fix(i int) {
	if i < 0 || i >= len(h.data) {
		panic("heap: Fix index out of range")
	}
	if !h.down(i) {
		h.up(i)
	}
}

func (h *Heap[T]) down(i int) bool {
	data := h.data
	n := len(data)
	less := h.less
	i0 := i
	for {
		left := 2*i + 1
		if left >= n || left < 0 { // left < 0 after int overflow
			break
		}
		j := left
		// find the smallest child
		if right := left + 1; right < n && less(data[right], data[left]) {
			j = right
		}
		if !less(data[j], data[i]) {
			break
		}
		data[i], data[j] = data[j], data[i]
		i = j
	}
	return i > i0
}

func (h *Heap[T]) up(i int) {
	data := h.data
	less := h.less
	for {
		parent := (i - 1) / 2
		if i == 0 || !less(data[i], data[parent]) {
			break
		}

		data[i], data[parent] = data[parent], data[i]
		i = parent
	}
}
