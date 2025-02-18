package queue

import (
	"github.com/RyoLena/Gadget"
	"github.com/RyoLena/Gadget/internal/queue"
	"github.com/stretchr/testify/assert"
	"testing"
)

func compare() Gadget.Comparator[int] {
	return Gadget.ComparatorRealNumber[int]
}

func TestNewPriorityQueue(t *testing.T) {
	testCases := []struct {
		name     string
		initSize int
		compare  Gadget.Comparator[int]
		wantErr  error
	}{
		{
			name:     "compare is nil",
			initSize: 8,
			compare:  nil,
		},
		{
			name:     "compare is ok",
			initSize: 8,
			compare:  compare(),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_ = NewPriorityQueue[int](tc.initSize, tc.compare)
		})
	}
}

func TestPriorityQueue_Len(t *testing.T) {
	testCases := []struct {
		name     string
		initSize int
		compare  Gadget.Comparator[int]
		wantLen  int
	}{
		{
			name:     "no err is ok",
			initSize: 8,
			compare:  compare(),
			wantLen:  0,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pq := NewPriorityQueue[int](tc.initSize, tc.compare)
			assert.Equal(t, tc.wantLen, pq.Len())
		})
	}
}

func TestPriorityQueue_Peek(t *testing.T) {
	testCases := []struct {
		name       string
		initSize   int
		compare    Gadget.Comparator[int]
		wantResult int
		wantErr    error
	}{
		{
			name:     "no err is ok",
			initSize: 8,
			compare:  compare(),
			wantErr:  queue.ErrEmptyQueue,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pq := NewPriorityQueue[int](tc.initSize, tc.compare)
			result, err := pq.Peek()
			assert.Equal(t, tc.wantResult, result)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestPriorityQueue_Enqueue(t *testing.T) {
	testCases := []struct {
		name        string
		initSize    int
		compare     Gadget.Comparator[int]
		enqueueData int
		wantErr     error
	}{
		{
			name:     "no err is ok",
			initSize: 8,
			compare:  compare(),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pq := NewPriorityQueue[int](tc.initSize, tc.compare)
			err := pq.Enqueue(tc.enqueueData)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestPriorityQueue_Dequeue(t *testing.T) {
	testCases := []struct {
		name       string
		initSize   int
		compare    Gadget.Comparator[int]
		wantResult int
		wantErr    error
	}{
		{
			name:     "no err is ok",
			initSize: 8,
			compare:  compare(),
			wantErr:  queue.ErrEmptyQueue,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pq := NewPriorityQueue[int](tc.initSize, tc.compare)
			result, err := pq.Dequeue()
			assert.Equal(t, tc.wantResult, result)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
