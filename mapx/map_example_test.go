package mapx

//func ExampleNewHashMap() {
//	m := NewHashMap[MockKey, int](10)
//	_ = m.Put(MockKey{}, 123)
//	val, _ := m.Get(MockKey{})
//	fmt.Println(val)
//	// Output:
//	// 123
//}

type MockKey struct {
	values []int
}

func (m MockKey) Code() uint64 {
	res := 3
	for _, v := range m.values {
		res += v * 7
	}
	return uint64(res)
}

func (m MockKey) Equals(key any) bool {
	k, ok := key.(MockKey)
	if !ok {
		return false
	}
	if len(k.values) != len(m.values) {
		return false
	}
	if k.values == nil && m.values != nil {
		return false
	}
	if k.values != nil && m.values == nil {
		return false
	}
	for i, v := range m.values {
		if v != k.values[i] {
			return false
		}
	}
	return true
}
