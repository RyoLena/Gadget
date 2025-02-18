package mapx

import "github.com/RyoLena/Gadget"

type LinkedMap[K any, V any] struct {
	m          mapi[K, *linkedKV[K, V]]
	head, tail *linkedKV[K, V]
	length     int
}

type linkedKV[K any, V any] struct {
	key        K
	value      V
	prev, next *linkedKV[K, V]
}

func NewLinkedHashMap[K Hashable, V any](size int) *LinkedMap[K, V] {
	hashmap := NewHashMap[K, *linkedKV[K, V]](size)
	head := &linkedKV[K, V]{}
	tail := &linkedKV[K, V]{next: head, prev: head}
	head.prev, head.next = tail, tail
	return &LinkedMap[K, V]{
		m:    hashmap,
		head: head,
		tail: tail,
	}
}

func NewLinkedTreeMap[K any, V any](comparator Gadget.Comparator[K]) (*LinkedMap[K, V], error) {
	treeMap, err := NewTreeMap[K, *linkedKV[K, V]](comparator)
	if err != nil {
		return nil, err
	}
	head := &linkedKV[K, V]{}
	tail := &linkedKV[K, V]{next: head, prev: head}
	head.prev, head.next = tail, tail
	return &LinkedMap[K, V]{
		m:    treeMap,
		head: head,
		tail: tail,
	}, nil
}

func (l *LinkedMap[K, V]) Put(key K, val V) error {
	if lk, ok := l.m.Get(key); ok {
		lk.value = val
		return nil
	}
	lk := &linkedKV[K, V]{
		key:   key,
		value: val,
		prev:  l.tail.prev,
		next:  l.tail,
	}
	if err := l.m.Put(key, lk); err != nil {
		return err
	}
	lk.prev.next, lk.next.prev = lk, lk
	l.length++
	return nil
}

func (l *LinkedMap[K, V]) Get(key K) (V, bool) {
	if lk, ok := l.m.Get(key); ok {
		return lk.value, ok
	}
	var v V
	return v, false
}

func (l *LinkedMap[K, V]) Delete(key K) (V, bool) {
	if lk, ok := l.m.Delete(key); ok {
		lk.prev.next = lk.next
		lk.next.prev = lk.prev
		l.length--
		return lk.value, ok
	}
	var v V
	return v, false
}

func (l *LinkedMap[K, V]) Keys() []K {
	keys := make([]K, 0, l.length)
	for cur := l.head.next; cur != l.tail; {
		keys = append(keys, cur.key)
		cur = cur.next
	}
	return keys
}

func (l *LinkedMap[K, V]) Values() []V {
	values := make([]V, 0, l.length)
	for cur := l.head.next; cur != l.tail; {
		values = append(values, cur.value)
		cur = cur.next
	}
	return values
}

func (l *LinkedMap[K, V]) Len() int64 {
	return int64(l.length)
}
