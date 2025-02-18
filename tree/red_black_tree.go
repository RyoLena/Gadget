package tree

import (
	"errors"
	"github.com/RyoLena/Gadget"
	"github.com/RyoLena/Gadget/internal/tree"
)

var (
	errRBTreeComparatorIsNull = errors.New("gadget:RBTree 的 Comparator 不能为 nil")
)

type RBTree[k any, v any] struct {
	rbTree *tree.RBTree[k, v] //红黑树本体
}

func NewRBTree[k any, v any](compare Gadget.Comparator[k]) (*RBTree[k, v], error) {
	if nil == compare {
		return nil, errRBTreeComparatorIsNull
	}
	return &RBTree[k, v]{
		rbTree: tree.NewRBTree[k, v](compare),
	}, nil
}

// Add 方法封装 增加节点
func (rb *RBTree[k, v]) Add(key k, value v) error {
	return rb.rbTree.Add(key, value)
}

func (rb *RBTree[k, v]) Delete(key k) (v, bool) {
	return rb.rbTree.Delete(key)
}

func (rb *RBTree[k, v]) Set(key k, value v) error {
	return rb.rbTree.Set(key, value)
}

func (rb *RBTree[k, v]) Find(key k) (v, error) {
	return rb.rbTree.Find(key)
}

// Size 返回红黑树结点个数
func (rb *RBTree[K, V]) Size() int {
	return rb.rbTree.Size()
}

// KeyValues 获取红黑树所有节点K,V
func (rb *RBTree[K, V]) KeyValues() ([]K, []V) {
	return rb.rbTree.KeyValues()
}
