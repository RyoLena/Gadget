package slice

/*
如果cap过小 <64 直接返回和返回一个没有操作的信息
如果cap大于2048并且容量占比>=2时，容量只缩减0.625
如果cap小于2048并且容量占比过大>4时，空置率过大 容量减少一半
*/
func calCapacity(cap, len int) (int, bool) {
	if cap <= 64 {
		return cap, false
	}
	if cap > 2048 && ((cap / len) >= 2) {
		factor := 0.625
		return int(float32(cap) * float32(factor)), true
	}
	if cap <= 2048 && (cap/len >= 4) {
		return cap / 2, true
	}
	return cap, false
}

func Shrink[T any](src []T) []T {
	c, l := cap(src), len(src)
	size, change := calCapacity(c, l)
	if !change {
		return src
	}
	s := make([]T, 0, size)
	s = append(s, src...)
	return s
}
