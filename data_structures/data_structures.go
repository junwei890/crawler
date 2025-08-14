package data_structures

import (
	"errors"
	"slices"
)

type Queue []string

type QueueOps interface {
	Enqueue(string)
	Dequeue() (string, error)
	Peek() (string, error)
	Empty() bool
	Size() int
}

func (q *Queue) Enqueue(url string) {
	*q = append(*q, url)
}

func (q *Queue) Dequeue() (string, error) {
	if len(*q) == 0 {
		return "", errors.New("can't pop from an empty queue")
	}

	popped := (*q)[0]
	*q = slices.Delete(*q, 0, 1)

	return popped, nil
}

func (q *Queue) Peek() (string, error) {
	if len(*q) == 0 {
		return "", errors.New("can't peek into an empty queue")
	}

	return (*q)[0], nil
}

func (q *Queue) CheckEmpty() bool {
	return len(*q) == 0
}

func (q *Queue) Size() int {
	return len(*q)
}
