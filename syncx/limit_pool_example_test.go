package syncx

import (
	"fmt"
)

func ExampleNewLimitPool() {
	p := NewLimitPool(1, func() int {
		return 123
	})
	val, ok := p.Get()
	fmt.Println("第一次", val, ok)
	val, ok = p.Get()
	fmt.Println("第二次", val, ok)
	p.Put(123)
	val, ok = p.Get()
	fmt.Println("第三次", val, ok)
	// Output:
	// 第一次 123 true
	// 第二次 0 false
	// 第三次 123 true
}
