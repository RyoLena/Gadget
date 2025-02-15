package slice

import "Gadget/internal/slice"

// Add 在index处添加元素
// index 范围应为[0, len(src)]
// 如果index == len(src) 则表示往末尾添加元素
func Add[Src any](src []Src, element Src, index int) ([]Src, error) {
	return slice.Add(src, element, index)
}
