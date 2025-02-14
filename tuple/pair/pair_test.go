package pair

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

type testPairSuite struct{ suite.Suite }

func (s *testPairSuite) TestString() {
	{
		p := NewPair(100, "23333")
		s.Assert().Equal("<100,\"23333\">", p.String())
	}
	{
		p := NewPair("testStruct", map[int]int{
			11: 1,
			22: 2,
			33: 3,
		})
		s.Assert().Equal("<\"testStruct\", map[int]int{11:1, 22:2, 33:3}>", p.String())
	}
}

func (s *testPairSuite) TestNewPairs(t *testing.T) {
	tests := []struct {
		name        string
		keys        []int // 使用int类型测试泛型
		values      []string
		wantPairs   []Pair[int, string]
		wantErr     bool
		errContains string
	}{
		{
			name:   "正常情况",
			keys:   []int{1, 2},
			values: []string{"a", "b"},
			wantPairs: []Pair[int, string]{
				NewPair(1, "a"),
				NewPair(2, "b"),
			},
		},
		{
			name:      "空切片",
			keys:      []int{},
			values:    []string{},
			wantPairs: []Pair[int, string]{},
		},
		{
			name:        "keys是nil",
			keys:        nil,
			values:      []string{"a"},
			wantErr:     true,
			errContains: "均不为nil",
		},
		{
			name:        "values是nil",
			keys:        []int{1},
			values:      nil,
			wantErr:     true,
			errContains: "均不为nil",
		},
		{
			name:        "长度不匹配",
			keys:        []int{1, 2},
			values:      []string{"a"},
			wantErr:     true,
			errContains: "长度不同",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPairs(tt.keys, tt.values)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.wantPairs, got)
		})
	}
}

// 边界条件测试
func (s *testPairSuite) TestNewPairs_Boundary(t *testing.T) {
	t.Run("大尺寸切片", func(t *testing.T) {
		size := 10000
		keys := make([]int, size)
		values := make([]string, size)
		for i := 0; i < size; i++ {
			keys[i] = i
			values[i] = string(rune(i%26 + 97))
		}

		pairs, err := NewPairs(keys, values)
		assert.NoError(t, err)
		assert.Len(t, pairs, size)

		// 验证最后一个元素
		lastKey, lastVal := pairs[size-1].Split()
		assert.Equal(t, size-1, lastKey)
		assert.Equal(t, string(rune((size-1)%26+97)), lastVal)
	})

	t.Run("不同基础类型", func(t *testing.T) {
		keys := []float64{1.1, 2.2}
		values := []bool{true, false}

		pairs, err := NewPairs(keys, values)
		assert.NoError(t, err)
		assert.Equal(t, 1.1, pairs[0].Key)
		assert.True(t, pairs[0].Value)
	})
}

func (s *testPairSuite) TestSplitPairs() {
	type caseType struct {
		pairs []Pair[int, string]

		keys   []int
		values []string
	}
	for _, c := range []caseType{
		{
			pairs: []Pair[int, string]{
				NewPair(1, "1"),
				NewPair(2, "2"),
				NewPair(3, "3"),
				NewPair(4, "4"),
				NewPair(5, "5"),
			},
			keys:   []int{1, 2, 3, 4, 5},
			values: []string{"1", "2", "3", "4", "5"},
		},
		{
			pairs: nil,

			keys:   nil,
			values: nil,
		},
		{
			pairs:  []Pair[int, string]{},
			keys:   []int{},
			values: []string{},
		},
	} {
		keys, values := SplitPairs(c.pairs)
		if c.pairs == nil {
			s.Assert().Nil(keys)
			s.Assert().Nil(values)
		} else {
			s.Assert().Len(keys, len(c.pairs))
			s.Assert().Len(values, len(c.pairs))
			for i, pair := range c.pairs {
				s.Assert().Equal(pair.Key, keys[i])
				s.Assert().Equal(pair.Value, values[i])
			}
		}
	}
}

func (s *testPairSuite) TestFlattenPairs(t *testing.T) {
	type caseType struct {
		pairs []Pair[int, string]

		flattPairs []any
	}
	for _, c := range []caseType{
		{
			pairs: []Pair[int, string]{
				NewPair(1, "1"),
				NewPair(2, "2"),
				NewPair(3, "3"),
				NewPair(4, "4"),
				NewPair(5, "5"),
			},
			flattPairs: []any{1, "1", 2, "2", 3, "3", 4, "4", 5, "5"},
		},
		{
			pairs:      nil,
			flattPairs: nil,
		},
		{
			pairs:      []Pair[int, string]{},
			flattPairs: []any{},
		},
	} {
		flairs := FlattenPairs(c.pairs)
		s.Assert().EqualValues(c.flattPairs, flairs)
	}
}

func (s *testPairSuite) TestPackPairs(t *testing.T) {
	t.Run("正常情况-基础类型", func(t *testing.T) {
		flat := []any{1, "a", 2, "b", 3, "c"}
		pairs := PackPairs[int, string](flat)

		expected := []Pair[int, string]{
			NewPair(1, "a"),
			NewPair(2, "b"),
			NewPair(3, "c"),
		}
		assert.Equal(t, expected, pairs)
	})

	t.Run("nil输入", func(t *testing.T) {
		var flat []any = nil
		pairs := PackPairs[int, string](flat)
		assert.Nil(t, pairs)
	})

	t.Run("空切片", func(t *testing.T) {
		var flat []any
		pairs := PackPairs[string, float64](flat)
		assert.Empty(t, pairs)
	})

	t.Run("奇数长度切片", func(t *testing.T) {
		flat := []any{"a", 1, "b"} // 最后一个元素会被忽略
		pairs := PackPairs[string, int](flat)
		assert.Len(t, pairs, 1)
		assert.Equal(t, "a", pairs[0].Key)
		assert.Equal(t, 1, pairs[0].Value)
	})

	t.Run("类型不匹配引发panic", func(t *testing.T) {
		flat := []any{1, "a", "b", 2} // 第三个元素应该是int

		defer func() {
			if r := recover(); r == nil {
				t.Error("预期会发生panic但没有触发")
			}
		}()

		PackPairs[int, string](flat)
	})

	t.Run("混合类型验证", func(t *testing.T) {
		flat := []any{
			1.1, true,
			2.5, false,
			"3.0", 1, // 最后一个pair会panic
		}
		t.Run("正确类型部分", func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Log("捕获到预期的panic:", r)
				}
			}()

			pairs := PackPairs[float64, bool](flat[:4])
			expected := []Pair[float64, bool]{
				NewPair(1.1, true),
				NewPair(2.5, false),
			}
			assert.Equal(t, expected, pairs)
		})

		t.Run("错误类型部分", func(t *testing.T) {
			assert.PanicsWithValue(t, "interface conversion", func() {
				PackPairs[float64, int](flat[4:])
			})
		})
	})

	t.Run("复杂结构体类型", func(t *testing.T) {
		type Custom struct{ name string }

		flat := []any{
			"id1", Custom{name: "Alice"},
			100, Custom{name: "Bob"},
		}

		pairs := PackPairs[string, Custom](flat)
		require.Len(t, pairs, 2)
		assert.Equal(t, "id1", pairs[0].Key)
		assert.Equal(t, "Alice", pairs[0].Value.name)
		assert.Equal(t, 100, pairs[1].Key)
		assert.Equal(t, "Bob", pairs[1].Value.name)
	})
}

func TestPackPairs_EdgeCases(t *testing.T) {
	t.Run("超大切片", func(t *testing.T) {
		const size = 1e6 // 100万元素
		flat := make([]any, size*2)
		for i := 0; i < size; i++ {
			flat[i*2] = i
			flat[i*2+1] = i * 2
		}

		pairs := PackPairs[int, int](flat)
		assert.Len(t, pairs, size)
		assert.Equal(t, 999999, pairs[size-1].Key)
		assert.Equal(t, 1999998, pairs[size-1].Value)
	})

	t.Run("零值处理", func(t *testing.T) {
		flat := []any{0, "", nil, false}
		pairs := PackPairs[int, interface{}](flat)
		expected := []Pair[int, interface{}]{
			{Key: 0, Value: ""},
			{Key: -1, Value: false}, // 演示潜在类型问题
		}
		assert.Equal(t, expected, pairs) // 这个测试预期会失败，用于演示类型安全问题
	})
}

func TestPair(t *testing.T) {
	suite.Run(t, new(testPairSuite))
}
