package queue

import (
	"Gadget"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewPriorityQueue(t *testing.T) {
	type args[T any] struct {
		capacity int
		compare  Gadget.Comparator
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want *PriorityQueue[T]
	}
	tests := []testCase[ /* TODO: Insert concrete types here */ ]{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewPriorityQueue(tt.args.capacity, tt.args.compare), "NewPriorityQueue(%v, %v)", tt.args.capacity, tt.args.compare)
		})
	}
}

func TestPriorityQueue_Dequeue(t *testing.T) {
	type testCase[T any] struct {
		name    string
		pq      PriorityQueue[T]
		want    T
		wantErr assert.ErrorAssertionFunc
	}
	tests := []testCase[ /* TODO: Insert concrete types here */ ]{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.pq.Dequeue()
			if !tt.wantErr(t, err, fmt.Sprintf("Dequeue()")) {
				return
			}
			assert.Equalf(t, tt.want, got, "Dequeue()")
		})
	}
}

func TestPriorityQueue_Enqueue(t *testing.T) {
	type args[T any] struct {
		t T
	}
	type testCase[T any] struct {
		name    string
		pq      PriorityQueue[T]
		args    args[T]
		wantErr assert.ErrorAssertionFunc
	}
	tests := []testCase[ /* TODO: Insert concrete types here */ ]{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, tt.pq.Enqueue(tt.args.t), fmt.Sprintf("Enqueue(%v)", tt.args.t))
		})
	}
}

func TestPriorityQueue_Len(t *testing.T) {
	type testCase[T any] struct {
		name string
		pq   PriorityQueue[T]
		want int
	}
	tests := []testCase[ /* TODO: Insert concrete types here */ ]{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.pq.Len(), "Len()")
		})
	}
}

func TestPriorityQueue_Peek(t *testing.T) {
	type testCase[T any] struct {
		name    string
		pq      PriorityQueue[T]
		want    T
		wantErr assert.ErrorAssertionFunc
	}
	tests := []testCase[ /* TODO: Insert concrete types here */ ]{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.pq.Peek()
			if !tt.wantErr(t, err, fmt.Sprintf("Peek()")) {
				return
			}
			assert.Equalf(t, tt.want, got, "Peek()")
		})
	}
}
