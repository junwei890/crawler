package data_structures

import (
	"errors"
	"reflect"
	"testing"
)

func TestQueue(t *testing.T) {
	queue := &Queue{"a", "b", "c", "d", "e"}

	size := queue.Size()
	if size != 5 {
		t.Errorf("Queue: test case 1 failed, %d != %d", size, 5)
	}

	queue.Enqueue("f")
	if comp := reflect.DeepEqual(*queue, Queue{"a", "b", "c", "d", "e", "f"}); !comp {
		t.Errorf("Queue: test case 2 failed: %v != %v", *queue, Queue{"a", "b", "c", "d", "e", "f"})
	}

	size = queue.Size()
	if size != 6 {
		t.Errorf("Queue: test case 3 failed, %d != %d", size, 6)
	}

	popped, err := queue.Dequeue()
	if err != nil {
		t.Errorf("Queue: test case 4 failed, unexpected error: %v", err)
	}
	if popped != "a" {
		t.Errorf("Queue: test case 5 failed, %s != %s", popped, "a")
	}
	if comp := reflect.DeepEqual(*queue, Queue{"b", "c", "d", "e", "f"}); !comp {
		t.Errorf("Queue: test case 6 failed, %v != %v", *queue, Queue{"b", "c", "d", "e", "f"})
	}

	first, err := queue.Peek()
	if err != nil {
		t.Errorf("Queue: test case 7 failed, unexpected error: %v", err)
	}
	if first != "b" {
		t.Errorf("Queue: test case 8 failed, %s != %s", first, "b")
	}

	for size := queue.Size(); size > 0; size-- {
		queue.Dequeue()
	}

	empty := queue.CheckEmpty()
	if !empty {
		t.Errorf("Queue: test case 9 failed, %v != %v", empty, true)
	}

	_, err = queue.Dequeue()
	if (err != nil) != true {
		t.Errorf("Queue: test case 10 failed, expected error: %s", errors.New("queue empty"))
	}

	_, err = queue.Peek()
	if (err != nil) != true {
		t.Errorf("Queue: test case 11 failed, expected error: %s", errors.New("queue empty"))
	}
}
