package slice

import "github.com/RyoLena/Gadget/internal/errs"

func Add[T any](src []T, element T, index int) ([]T, error) {
	length := len(src)
	if index < 0 || index > length {
		return nil, errs.NewErrIndexOutOfRange(length, index)
	}
	//扩充src
	var zeroValue T
	src = append(src, zeroValue)
	for i := len(src) - 1; i > index; i-- {
		if i >= 0 {
			src[i] = src[i-1]
		}
	}
	src[index] = element
	return src, nil
}
