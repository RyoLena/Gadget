package pair

import (
	"github.com/stretchr/testify/assert"
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

func TestPair(t *testing.T) {
	suite.Run(t, new(testPairSuite))
}
