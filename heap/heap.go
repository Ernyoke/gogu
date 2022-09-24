// Package heap provides an implementation of the binary heap data structure.
// A common implementation of the heap is the binary tree, where each node
// of the subtree satisfies the heap property: each node of the subtree
// is greather or equal then the parent node in case of min heap,
// and less or equal than the parent node in case of max heap.
// The heap package implements each type, where
// the conditional function of the constructor defines the heap type.
package heap

import (
	"fmt"
	"sync"

	"github.com/esimov/gogu"
)

type Heap[T comparable] struct {
	mu   *sync.RWMutex
	data []T
	comp gogu.CompFn[T]
}

// NewHeap creates a new heap data structure having two components:
// a data slice holding the concrete values and a comparator function.
// The comparison sign decides if the heap is a max heap or min heap.
func NewHeap[T comparable](comp gogu.CompFn[T]) *Heap[T] {
	return &Heap[T]{
		mu:   new(sync.RWMutex),
		data: make([]T, 0),
		comp: comp,
	}
}

// Size returns the heap size.
func (h *Heap[T]) Size() int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return len(h.data)
}

// IsEmpty returns true if the heap is empty, otherwise false.
func (h *Heap[T]) IsEmpty() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return h.Size() == 0
}

// Clear removes all the elements from the heap.
func (h *Heap[T]) Clear() {
	if h.Size() == 0 {
		return
	}

	h.mu.Lock()
	h.data = h.data[:0]
	h.mu.Unlock()
}

// Peek returns the first element of the heap.
// This can be the minimum or maximum value depending on the heap type.
func (h *Heap[T]) Peek() T {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if h.Size() == 0 {
		defer h.mu.RLock()
		var t T
		return t
	}

	return h.data[0]
}

// GetValues returns the heap values.
func (h *Heap[T]) GetValues() []T {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return h.data
}

// Push inserts new elements at the end of the heap and calls the heapify algorithm to reorder
// the existing elements in ascending or descending order, depending on the heap type.
func (h *Heap[T]) Push(val ...T) {
	for _, v := range val {
		h.mu.Lock()
		h.data = append(h.data, v)
		h.mu.Unlock()

		h.moveUp(h.Size() - 1)
	}
}

// Pop removes the first element from the heap and reorder the existing elements.
// The removed element can be the minimum or maximum depending on the heap type.
func (h *Heap[T]) Pop() T {
	var val T
	len := h.Size()
	if h.Size() == 0 {
		return val
	}
	val = h.Peek()

	h.mu.Lock()
	h.data[0] = h.data[len-1]
	h.data = h.data[:len-1]
	h.mu.Unlock()

	h.moveDown(h.Size(), 0)

	return val
}

// Delete removes an element from the heap. In case the element does not exists it returns false and an error.
// After removal it reorders the heap following the heap specific rules.
func (h *Heap[T]) Delete(val T) (bool, error) {
	len := h.Size()
	if len == 0 {
		return false, fmt.Errorf("heap empty")
	}

	idx, ok := h.getIndex(h.data, val)
	if !ok {
		return false, fmt.Errorf("value not found in the heap: %v", val)
	}

	swap(h.mu, h.data, idx, len-1)
	h.data = h.data[:len-1]
	h.moveDown(len, 0)

	return true, nil
}

// Convert a min heap to max heap and vice versa.
func (h *Heap[T]) Convert(comp gogu.CompFn[T]) {
	h.mu.Lock()
	h.comp = comp
	h.mu.Unlock()

	// Start from bottom-rightmost internal mode and reorder all internal nodes.
	for i := (h.Size() - 2) / 2; i >= 0; i-- {
		h.moveDown(h.Size(), i)
	}
}

// FromSlice imports the slice elements into a new heap using the comparator function.
func FromSlice[T comparable](mu *sync.RWMutex, data []T, comp gogu.CompFn[T]) *Heap[T] {
	for i := len(data)/2 - 1; i >= 0; i-- {
		for {
			l, r := 2*i+1, 2*i+2
			if l >= len(data) || l < 0 {
				break
			}

			current := l
			if r < len(data) && comp(data[r], data[l]) {
				current = r
			}

			if !comp(data[current], data[i]) {
				break
			}

			swap(mu, data, i, current)
			i = current
		}
	}

	return &Heap[T]{
		mu:   mu,
		data: data,
		comp: comp,
	}
}

// Merge joins two heaps into a new one preserving the original heaps.
func (h *Heap[T]) Merge(h2 *Heap[T]) *Heap[T] {
	newHeap := NewHeap(h.comp)

	for i := 0; i < h.Size(); i++ {
		newHeap.Push(h.data[i])
	}

	for i := 0; i < h2.Size(); i++ {
		newHeap.Push(h2.data[i])
	}

	return newHeap
}

// Meld merge two heaps into a new one containing all the
// elements of both and destroying the original heaps.
func (h *Heap[T]) Meld(h2 *Heap[T]) *Heap[T] {
	newHeap := NewHeap(h.comp)

	for i := 0; i < h.Size(); i++ {
		newHeap.Push(h.data[i])
	}

	for i := 0; i < h2.Size(); i++ {
		newHeap.Push(h2.data[i])
	}
	h.data = nil
	h2.data = nil

	return newHeap
}

// moveDown moves the element at the position i down to its
// correct position in the heap following the heap rules.
func (h *Heap[T]) moveDown(n, i int) {
	h.mu.RLock()

	left := h.leftChild(i)
	right := h.rightChild(i)
	current := i

	if left < n && h.comp(h.data[left], h.data[current]) {
		current = left
	}

	if right < n && h.comp(h.data[right], h.data[current]) {
		current = right
	}

	if current != i {
		h.mu.RUnlock()
		swap(h.mu, h.data, i, current)
		h.moveDown(n, current)
		return
	}
	h.mu.RUnlock()
}

// moveUp moves the element from index i up to its
// correct position in the heap following the heap rules.
func (h *Heap[T]) moveUp(i int) {
	h.mu.RLock()
	if h.comp(h.data[i], h.data[h.parent(i)]) {
		h.mu.RUnlock()
		swap(h.mu, h.data, i, h.parent(i))

		h.mu.RLock()
		i = h.parent(i)
		h.mu.RUnlock()

		h.moveUp(i)
		return
	}
	h.mu.RUnlock()
}

// leftChild returns the index of the left child of node at index i.
func (h *Heap[T]) leftChild(i int) int {
	return 2*i + 1
}

// rightChild returns the index of the right child of node at index i.
func (h *Heap[T]) rightChild(i int) int {
	return 2*i + 2
}

// parent returns the index of the child node parent at index i.
func (h *Heap[T]) parent(i int) int {
	return (i - 1) / 2
}

// swap swaps the position of elements at index i and j.
func swap[T any](mu *sync.RWMutex, data []T, i, j int) {
	mu.Lock()
	data[i], data[j] = data[j], data[i]
	mu.Unlock()
}

func (h *Heap[T]) getIndex(slice []T, val T) (int, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for i := 0; i < len(slice); i++ {
		if slice[i] == val {
			return i, true
		}
	}

	return -1, false
}
