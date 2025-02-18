package slice

// Index 返回和 dst 相等的第一个元素下标
// -1 表示没找到
func Index[T comparable](src []T, dst T) int {
	return IndexFunc[T](src, func(src T) bool {
		return src == dst
	})
}

// IndexFunc 返回 match 返回 true 的第一个下标
// -1 表示没找到
// 你应该优先使用 Index
func IndexFunc[T any](src []T, match matchFunc[T]) int {
	for k, v := range src {
		if match(v) {
			return k
		}
	}
	return -1
}

// LastIndex 返回和 dst 相等的最后一个元素下标
// -1 表示没找到
func LastIndex[T comparable](src []T, dst T) int {
	return LastIndexFunc[T](src, func(src T) bool {
		return src == dst
	})
}

// LastIndexFunc 返回和 dst 相等的最后一个元素下标
// -1 表示没找到
// 你应该优先使用 LastIndex
func LastIndexFunc[T any](src []T, match matchFunc[T]) int {
	for i := len(src) - 1; i >= 0; i-- {
		if match(src[i]) {
			return i
		}
	}
	return -1
}

// IndexAll 返回和 dst 相等的所有元素的下标
func IndexAll[T comparable](src []T, dst T) []int {
	return IndexAllFunc[T](src, func(src T) bool {
		return src == dst
	})
}

// IndexAllFunc 返回和 match 返回 true 的所有元素的下标
// 你应该优先使用 IndexAll
func IndexAllFunc[T any](src []T, match matchFunc[T]) []int {
	var indexes = make([]int, 0, len(src))
	for k, v := range src {
		if match(v) {
			indexes = append(indexes, k)
		}
	}
	return indexes
}
