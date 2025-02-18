package syncx

import "fmt"

func ExampleNewSegmentKeysLock() {
	// 参数就是分多少段，你也可以理解为总共有多少锁
	// 锁越多，并发竞争越低，但是消耗内存；
	// 锁越少，并发竞争越高，但是内存消耗少；
	lock := NewSegmentKeysLock(100)
	// 对应的还有 TryLock
	// RLock 和 RUnlock
	lock.Lock("key1")
	defer lock.Unlock("key1")
	fmt.Println("OK")
	// Output:
	// OK
}
