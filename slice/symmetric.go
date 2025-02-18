package slice

// SymmetricDiffSet 对称差集
// 已去重
// 返回值的元素顺序是不定的
func SymmetricDiffSet[T comparable](src, dst []T) []T {
	srcMap, dstMap := toMap[T](src), toMap[T](dst)
	for k := range dstMap {
		if _, ok := srcMap[k]; ok {
			delete(srcMap, k)
		} else {
			srcMap[k] = struct{}{}
		}
	}

	res := make([]T, 0, len(srcMap))
	for k := range srcMap {
		res = append(res, k)
	}

	return res
}

// SymmetricDiffSetFunc 对称差集
// 你应该优先使用 SymmetricDiffSet
// 已去重
func SymmetricDiffSetFunc[T any](src, dst []T, equal equalFunc[T]) []T {
	res := []T{}

	//找出在src不在dst的元素
	for _, v := range src {
		if !ContainsFunc[T](dst, func(t T) bool {
			return equal(t, v)
		}) {
			res = append(res, v)
		}
	}

	//找出在dst不在src的元素
	for _, v := range dst {
		if !ContainsFunc[T](src, func(t T) bool {
			return equal(t, v)
		}) {
			res = append(res, v)
		}
	}

	return deduplicateFunc[T](res, equal)
}
