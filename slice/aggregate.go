package slice

import "github.com/RyoLena/Gadget"

// Max 返回最大值。
// 该方法假设你至少会传入一个值。
// 在使用 float32 或者 float64 的时候要小心精度问题
func Max[T Gadget.RealNumber](ts []T) T {
	res := ts[0]
	for i := 1; i < len(ts); i++ {
		if ts[i] > res {
			res = ts[i]
		}
	}
	return res
}

// Min 返回最小值
// 该方法会假设你至少会传入一个值
// 在使用 float32 或者 float64 的时候要小心精度问题
func Min[T Gadget.RealNumber](ts []T) T {
	res := ts[0]
	for i := 1; i < len(ts); i++ {
		if ts[i] < res {
			res = ts[i]
		}
	}
	return res
}

// Sum 求和
// 在使用 float32 或者 float64 的时候要小心精度问题
func Sum[T Gadget.Number](ts []T) T {
	var res T
	for _, n := range ts {
		res += n
	}
	return res
}
