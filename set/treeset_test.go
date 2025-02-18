package set

import (
	"github.com/RyoLena/Gadget"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewTreeSet(t *testing.T) {
	tcase := []struct {
		name       string
		keys       []int
		comparable Gadget.Comparator[int]
		wantKey    []int
		wantErr    string
	}{
		{
			name:       "comparable-nil",
			comparable: nil,
			keys:       nil,
			wantKey:    []int{},
			wantErr:    "gadget: Comparator不能为nil",
		},
		{
			name:       "key-nil",
			keys:       nil,
			comparable: compare(),
			wantKey:    []int{},
		},
		{
			name:       "duplicate-key",
			comparable: compare(),
			keys:       []int{0, 1, 2, 1},
			wantKey:    []int{0, 1, 2},
		},
		{
			name:       "disorder-key",
			comparable: compare(),
			keys:       []int{0, 2, 1, 6, 5},
			wantKey:    []int{0, 1, 2, 5, 6},
		},
	}
	for _, tt := range tcase {
		t.Run(tt.name, func(t *testing.T) {
			treeSet, err := NewTreeSet[int](tt.comparable)
			if err != nil {
				assert.Equal(t, tt.wantErr, err.Error())
				return
			}
			for i := 0; i < len(tt.keys); i++ {
				treeSet.Add(tt.keys[i])
			}
			keys := treeSet.Keys()
			assert.ElementsMatch(t, tt.wantKey, keys)
		})
	}
}

func TestTreeSet_Delete(t *testing.T) {
	tcase := []struct {
		name    string
		keys    []int
		k       int
		wantKey []int
	}{
		{
			name:    "key-nil",
			keys:    nil,
			k:       0,
			wantKey: []int{},
		},
		{
			name:    "find-key",
			keys:    []int{0, 1, 2},
			k:       0,
			wantKey: []int{1, 2},
		},
		{
			name:    "no-find-key",
			keys:    []int{0, 1, 2},
			k:       3,
			wantKey: []int{0, 1, 2},
		},
	}
	for _, tt := range tcase {
		t.Run(tt.name, func(t *testing.T) {
			treeSet, err := NewTreeSet[int](compare())
			require.NoError(t, err)
			for i := 0; i < len(tt.keys); i++ {
				treeSet.Add(tt.keys[i])
			}
			treeSet.Delete(tt.k)
			keys := treeSet.Keys()
			assert.Equal(t, tt.wantKey, keys)
		})
	}
}

func TestTreeSet_Exist(t *testing.T) {
	tcase := []struct {
		name string
		keys []int
		k    int
		want bool
	}{
		{
			name: "key-nil",
			keys: nil,
			k:    0,
			want: false,
		},
		{
			name: "find-key",
			keys: []int{0, 1, 2},
			k:    0,
			want: true,
		},
		{
			name: "no-find-key",
			keys: []int{0, 1, 2},
			k:    3,
			want: false,
		},
	}
	for _, tt := range tcase {
		t.Run(tt.name, func(t *testing.T) {
			treeSet, err := NewTreeSet[int](compare())
			require.NoError(t, err)
			for i := 0; i < len(tt.keys); i++ {
				treeSet.Add(tt.keys[i])
			}
			isExist := treeSet.Exist(tt.k)

			assert.Equal(t, tt.want, isExist)
		})
	}
}

func compare() Gadget.Comparator[int] {
	return Gadget.ComparatorRealNumber[int]
}

// goos: windows
// goarch: amd64
// pkg: github.com/gotomicro/Gadget/set
// cpu: Intel(R) Core(TM) i5-7300HQ CPU @ 2.50GHz
// BenchmarkTreeSet/treeSet_add-4            100000               431.6 ns/op
// BenchmarkTreeSet/set_add-4                100000                90.00 ns/op
// BenchmarkTreeSet/map_add-4                100000               117.5 ns/op
// BenchmarkTreeSet/treeSet_exist-4          100000               180.0 ns/op
// BenchmarkTreeSet/set_exist-4              100000                33.85 ns/op
// BenchmarkTreeSet/map_exist-4              100000                37.87 ns/op
// BenchmarkTreeSet/treeSet_del-4            100000               136.3 ns/op
// BenchmarkTreeSet/set_del-4                100000                58.19 ns/op
// BenchmarkTreeSet/map_del-4                100000                52.26 ns/op
func BenchmarkTreeSet(b *testing.B) {
	treeSet, err := NewTreeSet[uint64](Gadget.ComparatorRealNumber[uint64])
	require.NoError(b, err)
	s := NewMapSet[uint64](100)
	m := make(map[uint64]int, 100)
	b.Run("treeSet_add", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			treeSet.Add(uint64(i))
		}
	})

	b.Run("set_add", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			s.Add(uint64(i))
		}
	})
	b.Run("map_add", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			m[uint64(i)] = i
		}
	})

	b.Run("treeSet_exist", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			treeSet.Exist(uint64(i))
		}
	})
	b.Run("set_exist", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			s.Exist(uint64(i))
		}
	})
	b.Run("map_exist", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = m[uint64(i)]
		}
	})

	b.Run("treeSet_del", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			treeSet.Delete(uint64(i))
		}
	})
	b.Run("set_del", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			s.Delete(uint64(i))
		}
	})
	b.Run("map_del", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			delete(m, uint64(i))
		}
	})

}
