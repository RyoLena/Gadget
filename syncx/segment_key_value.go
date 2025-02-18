package syncx

import (
	"hash/fnv"
	"sync"
)

// SegmentKeysLock 部分key lock结构定义
type SegmentKeysLock struct {
	locks []*sync.RWMutex
	size  uint32
}

// NewSegmentKeysLock 创建 SegmentKeysLock 示例
func NewSegmentKeysLock(size uint32) *SegmentKeysLock {
	locks := make([]*sync.RWMutex, size)
	for i := range locks {
		locks[i] = &sync.RWMutex{}
	}
	return &SegmentKeysLock{
		locks: locks,
		size:  size,
	}
}

// hash 索引锁的hash函数
func (s *SegmentKeysLock) hash(key string) uint32 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(key))
	return h.Sum32()
}

// RLock 读锁加锁
func (s *SegmentKeysLock) RLock(key string) {
	s.getLock(key).RLock()
}

// TryRLock 试着加读锁，加锁成功会返回
func (s *SegmentKeysLock) TryRLock(key string) bool {
	return s.getLock(key).TryRLock()
}

// RUnlock 读锁解锁
func (s *SegmentKeysLock) RUnlock(key string) {
	s.getLock(key).RUnlock()
}

// Lock 写锁加锁
func (s *SegmentKeysLock) Lock(key string) {
	s.getLock(key).Lock()
}

// TryLock 试着加锁，加锁成功会返回 true
func (s *SegmentKeysLock) TryLock(key string) bool {
	return s.getLock(key).TryLock()
}

// Unlock 写锁解锁
func (s *SegmentKeysLock) Unlock(key string) {
	s.getLock(key).Unlock()
}

func (s *SegmentKeysLock) getLock(key string) *sync.RWMutex {
	hash := s.hash(key)
	return s.locks[hash%s.size]
}
