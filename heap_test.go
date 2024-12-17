package heap_test

import (
	"cmp"
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/gammazero/heap"
)

type testElem struct {
	key      string
	priority int
	index    int
}

var prioCmp = func(i, j *testElem) bool {
	return i.priority < j.priority
}

func TestPop(t *testing.T) {
	q := heap.New[*testElem](prioCmp)
	tasks := []testElem{
		{key: "a", priority: 9},
		{key: "b", priority: 4},
		{key: "c", priority: 3},
		{key: "d", priority: 0},
		{key: "e", priority: 6},
	}
	for _, e := range tasks {
		q.Push(&e)
	}
	var priorities []int
	var peekPriorities []int
	for q.Len() > 0 {
		i := q.Peek().priority
		t.Logf("peeked %v", i)
		peekPriorities = append(peekPriorities, i)
		j := q.Pop().priority
		t.Logf("popped %v", j)
		priorities = append(priorities, j)
	}
	if !sort.IntsAreSorted(peekPriorities) {
		t.Fatal("the values were not returned in sorted order")
	}
	if !sort.IntsAreSorted(priorities) {
		t.Fatal("the popped values were not returned in sorted order")
	}
}

func TestHeap(t *testing.T) {
	less := cmp.Less[int]
	h := heap.New(less)

	for i := 20; i > 10; i-- {
		h.Push(i)
	}
	verifyIntHeap(t, h, 0, less)

	for i := 10; i > 0; i-- {
		h.Push(i)
		verifyIntHeap(t, h, 0, less)
	}

	for i := 1; h.Len() > 0; i++ {
		x := h.Pop()
		if i < 20 {
			h.Push(20 + i)
		}
		if x != i {
			t.Errorf("%d.th pop got %d; want %d", i, x, i)
		}
	}
}

func TestRemove0(t *testing.T) {
	less := cmp.Less[int]
	h := heap.New(less)
	for i := 0; i < 10; i++ {
		h.Push(i)
	}
	verifyIntHeap(t, h, 0, less)

	for h.Len() > 0 {
		i := h.Len() - 1
		x := h.Remove(i)
		if x != i {
			t.Errorf("Remove(%d) got %d; want %d", i, x, i)
		}
		verifyIntHeap(t, h, 0, less)
	}
}

func TestRemove1(t *testing.T) {
	less := cmp.Less[int]
	h := heap.New(less)
	for i := 0; i < 10; i++ {
		h.Push(i)
	}

	for i := 0; h.Len() > 0; i++ {
		x := h.Remove(0)
		if x != i {
			t.Errorf("Remove(0) got %d; want %d", x, i)
		}
		verifyIntHeap(t, h, 0, less)
	}
}

func TestRemove2(t *testing.T) {
	N := 10

	less := cmp.Less[int]
	h := heap.New(less)
	for i := 0; i < N; i++ {
		h.Push(i)
	}

	m := make(map[int]bool)
	for h.Len() > 0 {
		m[h.Remove((h.Len()-1)/2)] = true
		verifyIntHeap(t, h, 0, less)
	}

	if len(m) != N {
		t.Errorf("len(m) = %d; want %d", len(m), N)
	}
	for i := 0; i < len(m); i++ {
		if !m[i] {
			t.Errorf("m[%d] doesn't exist", i)
		}
	}
}

func TestCreateHeapFromSlice(t *testing.T) {
	cases := []struct {
		name   string
		slice  []int
		sorted []int
		less   func(int, int) bool
	}{
		{
			name:  "empty",
			slice: []int{},
			less:  func(a, b int) bool { return a < b },
		},
		{
			name:   "non-empty (minheap)",
			slice:  []int{6, 3, 7, 5, 2, 4, 1},
			sorted: []int{1, 2, 3, 4, 5, 6, 7},
			less:   func(a, b int) bool { return a < b },
		},
		{
			name:   "non-empty (maxheap)",
			slice:  []int{-3, 5, 7, 9, 2, -1},
			sorted: []int{9, 7, 5, 2, -1, -3},
			less:   func(a, b int) bool { return a > b },
		},
	}

	for _, c := range cases {
		t.Run(c.name+" From", func(t *testing.T) {
			heap := heap.NewFrom(c.less, c.slice...)

			remaining := len(c.sorted)
			for i, v := range c.sorted {
				if heap.Len() != remaining {
					t.Errorf("heap does not have expectd number of elements, got %d want %d", heap.Len(), remaining)
				}
				val := heap.Pop()
				if val != v {
					t.Errorf("peek not equal, idx: %v", i)
				}
				remaining--
			}
		})
	}
}

func TestAtOutOfRangePanics(t *testing.T) {
	h := heap.New(cmp.Less[int])

	h.Push(1)
	h.Push(2)
	h.Push(3)

	assertPanics(t, "should panic when negative index", func() {
		h.At(-4)
	})

	assertPanics(t, "should panic when index greater than length", func() {
		h.At(3)
	})
}

func TestRemoveOutOfRangePanics(t *testing.T) {
	h := heap.New(cmp.Less[int])

	h.Push(1)
	h.Push(2)
	h.Push(3)

	assertPanics(t, "should panic when negative index", func() {
		h.Remove(-4)
	})

	assertPanics(t, "should panic when index greater than length", func() {
		h.Remove(3)
	})
}

func TestPopEmptyPanics(t *testing.T) {
	h := heap.New(cmp.Less[int])

	assertPanics(t, "should panic when popping empty heap", func() {
		h.Pop()
	})
	assertPanics(t, "should panic when peeking empty heap", func() {
		h.Peek()
	})

	h.Push(1)
	h.Pop()

	assertPanics(t, "should panic when popping emptied heap", func() {
		h.Pop()
	})
	assertPanics(t, "should panic when peeking emptied heap", func() {
		h.Peek()
	})
}

func assertPanics(t *testing.T, name string, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("%s: didn't panic as expected", name)
		}
	}()

	f()
}

func verifyIntHeap(t *testing.T, h *heap.Heap[int], i int, less func(a, b int) bool) {
	t.Helper()
	n := h.Len()
	j1 := 2*i + 1
	j2 := 2*i + 2
	if j1 < n {
		if less(h.At(j1), h.At(i)) {
			t.Errorf("heap invariant invalidated [%d] = %d > [%d] = %d", i, h.At(i), j1, h.At(j1))
			return
		}
		verifyIntHeap(t, h, j1, less)
	}
	if j2 < n {
		if less(h.At(j2), h.At(i)) {
			t.Errorf("heap invariant invalidated [%d] = %d > [%d] = %d", i, h.At(i), j1, h.At(j2))
			return
		}
		verifyIntHeap(t, h, j2, less)
	}
}

func Example() {
	h := heap.New(func(a, b int) bool { return a < b })

	h.Push(103)
	h.Push(101)
	h.Push(102)

	fmt.Println(h.Len())

	v := h.Pop()
	fmt.Println(v)

	v = h.Peek()
	fmt.Println(v)
	// Output:
	// 3
	// 101
	// 102

}

func ExampleNewFrom() {
	h := heap.NewFrom(func(a, b int) bool { return a < b }, 5, 2, 3)

	v := h.Pop()
	fmt.Println(v)

	v = h.Peek()
	fmt.Println(v)
	// Output:
	// 2
	// 3
}

func Example_strings() {
	h := heap.New(func(a, b string) bool {
		return strings.Compare(a, b) < 0
	})
	h.Push("foo")
	h.Push("bar")
	h.Push("baz")

	fmt.Println(h.Len())  // Prints: 3
	fmt.Println(h.Peek()) // Prints: bar
	fmt.Println(h.Pop())  // Prints: bar
	fmt.Println(h.Pop())  // Prints: baz
	fmt.Println(h.Pop())  // Prints: foo

	h.Push("hello")
	h.Push("world")

	// Consume heap and print elements.
	for h.Len() != 0 {
		fmt.Println(h.Pop())
	}

	// Output:
	// 3
	// bar
	// bar
	// baz
	// foo
	// hello
	// world
}