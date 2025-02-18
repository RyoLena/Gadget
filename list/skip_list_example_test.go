package list

import (
	"fmt"
	"github.com/RyoLena/Gadget"
	"github.com/RyoLena/Gadget/internal/list"
)

func ExampleNewSkipList() {
	l := list.NewSkipList[int](Gadget.ComparatorRealNumber[int])
	l.Insert(123)
	val, _ := l.Get(0)
	fmt.Println(val)
	// Output:
	// 123
}
