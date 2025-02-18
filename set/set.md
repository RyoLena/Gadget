# TreeSet 与 TreeMap 技术文档

## 1. 基本概念
```go
// TreeSet 是基于 TreeMap 实现的集合
type TreeSet[T any] struct {
    treeMap *mapx.TreeMap[T, any] // 使用 TreeMap 存储数据（值部分为 nil）
}

// TreeMap 是基于红黑树实现的键值对容器
type TreeMap[K any, V any] struct {
    tree *tree.RBTree[K, V] // 底层红黑树结构
}
```

## 2. 核心功能对比

| 功能        | TreeSet                     | TreeMap                   |
|-----------|-----------------------------|---------------------------|
| 存储方式     | 唯一元素集合                   | 键值对存储                  |
| 底层实现     | 基于 TreeMap（值固定为 nil）    | 红黑树实现                  |
| 时间复杂度   | 增删查 O(log n)              | 增删查 O(log n)           |
| 排序特性    | 元素按比较器顺序存储             | 键按比较器顺序存储           |

## 3. 基础使用

### 3.1 创建实例
```go
// 创建整型 TreeSet
set, _ := NewTreeSet[int](Gadget.IntComparator)

// 创建字符串 TreeMap
map, _ := mapx.NewTreeMap[string](Gadget.StringComparator)
```

### 3.2 基本操作
```go
// 添加元素
set.Add(42)
map.Put("key", "value")

// 检查存在
exists := set.Exist(42)  // 返回 true
val, found := map.Get("key")

// 删除元素
set.Delete(42)
map.Delete("key")
```

## 4. 实现原理

### 4.1 红黑树特性
```mermaid
graph TD
    A[红黑树] --> B[自平衡二叉搜索树]
    A --> C[节点有颜色标记]
    A --> D[保证基本操作时间复杂度 O(log n)]
```

### 4.2 TreeSet 结构
```
TreeSet → TreeMap → 红黑树
            │
            └── 只使用键（值固定为 nil）
```

## 5. 重要方法说明

### 5.1 TreeSet 方法
```go
func (s *TreeSet[T]) Add(key T)    // 添加唯一元素
func (s *TreeSet[T]) Keys() []T    // 返回无序元素列表
func (s *TreeSet[T]) Len() int64   // 通过 treeMap.Len() 获取
```

### 5.2 TreeMap 方法
```go
func Put(key K, value V) error    // 插入键值对（自动替换旧值）
func Get(key K) (V, bool)         // 安全获取值
func Delete(k T) (V, bool)        // 返回被删除的值
```

## 6. 使用注意事项
1. 🔑 必须提供 Comparator 比较器
2. ⚠️ 非线程安全（需自行加锁）
3. 📊 元素数量越大性能优势越明显
4. 🔄 修改比较器会导致结构失效
5. 🚫 Keys() 方法返回顺序不可依赖

## 7. 最佳实践建议
```go
// 自定义结构体的比较器示例
type User struct { ID int; Name string }

userComparator := func(u1, u2 User) int {
    return Gadget.IntComparator(u1.ID, u2.ID)
}

// 初始化带自定义比较器的集合
userSet, _ := NewTreeSet[User](userComparator)
```

## 8. 常见问题
Q: 为什么选择红黑树实现？
A: 在频繁插入删除的场景下，红黑树相比 AVL 树有更好的综合性能

Q: 如何保证元素唯一性？
A: 通过比较器判断元素相等性，相同元素会被覆盖

Q: 如何处理并发访问？
A: 建议在使用层通过 sync.Mutex 等机制实现同步控制
```