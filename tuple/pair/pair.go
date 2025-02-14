package pair

import "fmt"

type Pair[K any, V any] struct {
	Key   K
	Value V
}

// 拼接字符串
func (pair *Pair[K, V]) String() string {
	return fmt.Sprintf("%v:%v", pair.Key, pair.Value)
}

// Split 方法将Key，Value作为返回参数传出
func (pair *Pair[K, V]) Split() (K, V) {
	return pair.Key, pair.Value
}

func NewPair[K, V any](Key K, Value V) Pair[K, V] {
	return Pair[K, V]{
		Key:   Key,
		Value: Value,
	}
}

// NewPairs 需要传入两个长度相同并且均不为nil的数组 keys 和 values，
/*
设keys长度为n，返回一个长度为n的pair数组。
保证：
	返回的pair数组满足条件（设pair数组为p）:
		对于所有的 0 <= i < n
		p[i].Key == keys[i] 并且 p[i].Value == values[i]
	如果传入的keys或者values为nil，会返回error
	如果传入的keys长度与values长度不同，会返回error
*/
func NewPairs[k any, v any](keys []k, values []v) (
	[]Pair[k, v], error) {
	if keys == nil || values == nil {
		return nil, fmt.Errorf("keys和values均不为nil")
	}
	n := len(keys)
	if n != len(values) {
		return nil, fmt.Errorf("keys与values的长度不同, "+
			"len(keys)=%d, len(values)=%d", n, len(values))
	}
	pairs := make([]Pair[k, v], n)
	for i := 0; i < n; i++ {
		pairs[i] = NewPair(keys[i], values[i])
	}
	return pairs, nil
}

/*
SplitPairs 需要传入一个[]Pair[K, V]，数组可以为nil。
设pairs数组的长度为n，返回两个长度均为n的数组keys, values。
如果pairs数组是nil, 则返回的keys与values也均为nil。
*/
func SplitPairs[k any, v any](pairs []Pair[k, v]) (keys []k,
	values []v) {
	if pairs == nil {
		return nil, nil
	}
	n := len(pairs)
	key := make([]k, n)
	value := make([]v, n)
	for i, pair := range pairs {
		key[i], value[i] = pair.Split()
	}
	return key, value
}

/*
FlattenPairs 需要传入一个[]Pair[K, V]，数组可以为nil
如果pairs数组为nil，则返回的flatPairs数组也为nil

	设pairs数组长度为n，保证返回的flatPairs数组长度为2 * n且满足:
		对于所有的 0 <= i < n
		flatPairs[i * 2] == pairs[i].Key
		flatPairs[i * 2 + 1] == pairs[i].Value
*/
func FlattenPairs[k any, v any](pairs []Pair[k, v]) (flattenPair []any) {
	if pairs == nil {
		return nil
	}
	n := len(pairs)
	flattenPair = make([]any, 2*n)
	for _, pair := range pairs {
		flattenPair = append(flattenPair, pair.Key, pair.Value)
	}
	return flattenPair
}

/*
PackPairs 需要传入一个长度为2 * n的数组flatPairs，数组可以为nil。

	函数将会返回一个长度为n的pairs数组，pairs满足
		对于所有的 0 <= i < n
		pairs[i].Key == flatPairs[i * 2]
		pairs[i].Value == flatPairs[i * 2 + 1]
	如果flatPairs为nil,则返回的pairs也为nil

	入参flatPairs需要满足以下条件：
		对于所有的 0 <= i < n
		flatPairs[i * 2] 的类型为 K
		flatPairs[i * 2 + 1] 的类型为 V
	否则会panic
*/
func PackPairs[k, v any](flatPairs []any) (pairs []Pair[k, v]) {
	if flatPairs == nil {
		return nil
	}
	n := len(flatPairs) / 2
	pairs = make([]Pair[k, v], n)
	for i := 0; i < n; i++ {
		pairs[i] = NewPair(flatPairs[i*2].(k), flatPairs[i*2+1].(v))
	}

	return pairs
}
