# heap

[![GoDoc](https://pkg.go.dev/badge/github.com/gammazero/heap)](https://pkg.go.dev/github.com/gammazero/heap)
[![Build Status](https://github.com/gammazero/heap/actions/workflows/go.yml/badge.svg)](https://github.com/gammazero/heap/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/gammazero/heap)](https://goreportcard.com/report/github.com/gammazero/heap)
[![codecov](https://codecov.io/gh/gammazero/heap/branch/master/graph/badge.svg)](https://codecov.io/gh/gammazero/heap)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

Generic implementation of a binary heap.

A binary heap is a tree with the property that each node is the minimum-valued node in its subtree. This implementation allows the caller to provide a `less` function that determined how the heap is ordered.

The minimum element in the tree is the root, at index 0.

A heap is a common way to implement a priority queue. To build a priority queue, create a Heap of the type of elements it will hold and specify a "less" function that orders the elements by priority. Use `Push` to add items and `Pop` to remove the item with the greatest precedence.

## Installation

```
$ go get github.com/gammazero/heap
```

## Reading Empty Heap

Since it is OK for the heap to contain an element's zero-value, it is necessary to either panic or return a second boolean value to indicate the heap is empty, when reading or removing an element. This heap panics when reading from an empty heap. This is a run-time check to help catch programming errors, which may be missed if a second return value is ignored. Simply check `Heap.Len()` before reading from the heap.

## Generics

Heap uses generics to create a Heap that contains items of the type specified. To create a Heap that holds a specific type, provide a type argument with the `Heap` variable declaration. For example:
```go
    intHeap := heap.New(func(a, b int) bool {
        return a < b})
```

## Example

```go
package main

import (
    "fmt"
    "strings"
    
    "github.com/gammazero/heap"
)

func main() {
	h := heap.New(func(a, b string) bool {
		return strings.Compare(a, b) < 0
	})
	h.Push("foo")
	h.Push("bar")
	h.Push("baz")

	fmt.Println(h.Len())
	fmt.Println(h.Peek())
	fmt.Println(h.Pop())
	fmt.Println(h.Pop())
	fmt.Println(h.Pop())

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
```
