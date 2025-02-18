package queue

import (
	"fmt"
	"github.com/RyoLena/Gadget"
)

func ExampleNewPriorityQueue() {
	// 容量，并且队列里面放的是 int
	pq := NewPriorityQueue(10, Gadget.ComparatorRealNumber[int])
	_ = pq.Enqueue(10)
	_ = pq.Enqueue(9)
	val, _ := pq.Dequeue()
	fmt.Println(val)
	// Output:
	// 9
}
