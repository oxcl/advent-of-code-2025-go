package main

import (
	"errors"
)

// Queue represents a generic FIFO data structure
type Queue[T any] struct {
	items []T
}

// NewQueue creates and returns a new empty queue
func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{
		items: []T{},
	}
}

// Enqueue adds an item to the back of the queue
func (q *Queue[T]) Enqueue(item T) {
	q.items = append(q.items, item)
}

// Dequeue removes and returns the front item from the queue
func (q *Queue[T]) Dequeue() (T, error) {
	var zero T
	if q.IsEmpty() {
		return zero, errors.New("queue is empty")
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item, nil
}

// Peek returns the front item without removing it
func (q *Queue[T]) Peek() (T, error) {
	var zero T
	if q.IsEmpty() {
		return zero, errors.New("queue is empty")
	}
	return q.items[0], nil
}

// IsEmpty returns true if the queue is empty
func (q *Queue[T]) IsEmpty() bool {
	return len(q.items) == 0
}

// Size returns the number of items in the queue
func (q *Queue[T]) Size() int {
	return len(q.items)
}

// Clear removes all items from the queue
func (q *Queue[T]) Clear() {
	q.items = []T{}
}
