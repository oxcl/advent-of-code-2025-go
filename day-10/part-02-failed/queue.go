package main

import (
	"errors"
)

type QueueNode[T any] struct {
	value *T
	cost  float64
}

type CostQueue[T any] struct {
	items []QueueNode[T]
}

func NewQueue[T any]() *CostQueue[T] {
	return &CostQueue[T]{
		items: []QueueNode[T]{},
	}
}

func (q *CostQueue[T]) Add(item *T, cost float64) {
	q.items = append(q.items, QueueNode[T]{
		cost:  cost,
		value: item,
	})
	hippifyUp(q, len(q.items)-1)
}

func hippifyUp[T any](q *CostQueue[T], index int) {
	if index == 0 || q.items[index].cost >= q.items[(index-1)/2].cost {
		return
	}
	q.items[index], q.items[(index-1)/2] = q.items[(index-1)/2], q.items[index]
	hippifyUp(q, (index-1)/2)
}

func (q *CostQueue[T]) Pop() (*T, error) {
	if len(q.items) == 0 {
		return nil, errors.New("queue is empty")
	}
	item := q.items[0]
	q.items[0] = q.items[len(q.items)-1]
	q.items = q.items[0 : len(q.items)-1]
	hippifyDown(q, 0)
	return item.value, nil
}

func hippifyDown[T any](q *CostQueue[T], index int) {
	if index*2+1 >= len(q.items) {
		return
	}
	var smallerChildIdx int
	if index*2+2 < len(q.items) && q.items[index*2+2].cost < q.items[index*2+1].cost {
		smallerChildIdx = index*2 + 2
	} else {
		smallerChildIdx = index*2 + 1
	}
	if q.items[index].cost > q.items[smallerChildIdx].cost {
		q.items[index], q.items[smallerChildIdx] = q.items[smallerChildIdx], q.items[index]
		hippifyDown(q, smallerChildIdx)
	}
}

// Size returns the number of items in the queue
func (q *CostQueue[T]) Size() int {
	return len(q.items)
}
