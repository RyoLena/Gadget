
# mapx 工具包技术文档

## 1. 核心数据结构概览
📦 提供多种增强型Map实现：

| 类型            | 特点                          | 适用场景                  |
|----------------|-----------------------------|-----------------------|
| TreeMap        | 按键排序（红黑树实现）               | 需要有序遍历的场景             |
| HashMap        | 哈希表+链表冲突解决                 | 高频读写无序遍历              |
| MultiMap       | 一键多值（支持Tree/HashMap实现）     | 分组聚合数据存储              |
| LinkedMap      | 维护插入顺序（链表+底层Map）          | 需要保留操作历史的场景           |
| BuiltinMap     | 原生map封装                    | 简单键值存储                |

## 2. 基础使用示例

### 2.1 TreeMap
```go
// 创建比较器
comp := Gadget.IntComparator 

// 初始化
treeMap, _ := mapx.NewTreeMap[int, string](comp)

// 操作
treeMap.Put(100, "百分")
val, _ := treeMap.Get(100)  // 返回"百分"
```

### 2.2 HashMap
```go
type UserID struct {
    ID int
}
// 实现Hashable接口
func (u UserID) Code() uint64 { return uint64(u.ID) }
func (u UserID) Equals(key any) bool { 
    other, ok := key.(UserID)
    return ok && u.ID == other.ID 
}

// 初始化
hashMap := mapx.NewHashMap[UserID, string](100)

// 操作
hashMap.Put(UserID{101}, "用户101")
name, _ := hashMap.Get(UserID{101}) // 返回"用户101"
```

### 2.3 MultiMap
```go
// 创建基于TreeMap的MultiMap
multi, _ := mapx.NewMultiTreeMap[int, string](Gadget.IntComparator)

// 添加数据
multi.Put(1, "A")
multi.Put(1, "B")

vals, _ := multi.Get(1) // 返回["A", "B"]
```

## 3. 高级功能

### 3.1 LinkedMap顺序维护
```go
linked := NewLinkedHashMap[int, string](10)
linked.Put(1, "A")
linked.Put(2, "B")

// 遍历顺序保证插入顺序
keys := linked.Keys() // 总是返回[1, 2]
```

### 3.2 数据转换
```go
// 切片转Map
keys := []int{1,2,3}
vals := []string{"A","B","C"}
m, _ := mapx.ToMap(keys, vals) // 得到map[1:A 2:B 3:C]
```

## 4. 性能注意事项

| 操作        | TreeMap | HashMap | LinkedMap |
|-----------|---------|---------|-----------|
| Put       | O(log n)| O(1)*   | O(1)      |
| Get       | O(log n)| O(1)*   | O(1)      |
| Delete    | O(log n)| O(1)*   | O(1)      |
| Traverse  | 有序     | 无序     | 插入顺序     |

* 哈希冲突时退化为O(n)

## 5. 使用规范

### 5.1 初始化要求
```go
// TreeMap必须提供比较器
mapx.NewTreeMap[string, int](Gadget.StringComparator)

// HashMap的Key必须实现Hashable接口
type MyKey struct{} 
func (k MyKey) Code() uint64 {...}
func (k MyKey) Equals(any) bool {...}
```

### 5.2 并发安全
⚠️ 所有实现均**非线程安全**，建议：
```go
var mu sync.RWMutex

// 写操作
mu.Lock()
defer mu.Unlock()
map.Put(key, val)

// 读操作
mu.RLock()
defer mu.RUnlock()
map.Get(key)
```

## 6. 常见问题

Q：MultiMap的Values()为何返回副本？  
A：防止外部修改影响内部数据，保证数据一致性

Q：如何选择Map实现？
```
需要排序 → TreeMap  
高频读写 → HashMap  
维护插入顺序 → LinkedMap  
分组存储 → MultiMap
```

Q：HashMap出现性能下降？  
A：检查：
1. 哈希函数分布是否均匀
2. 冲突链表是否过长
3. 初始容量是否过小

Q：TreeMap比较器示例？
```go
// 自定义结构比较
type User struct{ Age int }

UserComparator := func(u1, u2 User) int {
    return Gadget.IntComparator(u1.Age, u2.Age)
}
```