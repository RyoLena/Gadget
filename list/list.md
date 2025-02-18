# 数据结构工具包技术文档

## 1. 核心数据结构概览

| 数据结构          | 特点                          | 适用场景                  |
|------------------|-----------------------------|-------------------------|
| SkipList         | 跳表结构，O(log n) 复杂度       | 高频查询和有序数据访问场景       |
| ConcurrentList   | 线程安全链表（读写锁封装）           | 多线程环境下的数据操作           |
| List             | 通用列表接口                   | 统一数据结构操作规范            |

## 2. 核心接口定义 (List)

```go
type List[T any] interface {
    Get(index int) (T, error)      // 按索引获取元素
    Append(ts ...T) error          // 尾部追加元素
    Add(index int, t T) error      // 指定位置插入元素
    Set(index int, t T) error      // 修改指定位置元素
    Delete(index int) (T, error)   // 删除元素并返回值
    Len() int                      // 获取当前长度
    Cap() int                      // 获取容量
    Range(func(int, T) error) error // 遍历元素
    AsSlice() []T                  // 转换为切片
}
```

## 3. SkipList 跳表实现

### 3.1 创建实例
```go
// 需要提供元素比较器
comp := Gadget.IntComparator // 示例使用整型比较器
sl := list.NewSkipList[int](comp)
```

### 3.2 核心方法
```go
sl.Insert(42)                  // 插入元素
exists := sl.Search(42)        // 查询元素存在性
sl.DeleteElement(42)           // 删除元素
slice := sl.AsSlice()          // 获取有序切片
length := sl.Len()             // 获取元素数量
```

### 3.3 性能特点
- 平均时间复杂度：
    - 插入：O(log n)
    - 查询：O(log n)
    - 删除：O(log n)
- 空间复杂度：O(n log n)

## 4. ConcurrentList 并发链表

### 4.1 线程安全实现
```go
// 包装任意 List 实现
baseList := NewLinkedList[string]() 
cl := &ConcurrentList[string]{
    List: baseList,
}

// 并发安全操作
go func() {
    cl.Append("data1") 
}()
go func() {
    cl.Delete(0)
}()
```

### 4.2 锁机制说明
| 操作类型 | 锁类型      | 说明                |
|--------|------------|--------------------|
| 读操作   | 读锁（RLock）| 允许多个读操作并行       |
| 写操作   | 写锁（Lock） | 互斥访问，保证原子性      |

## 5. 使用示例对比

### 5.1 SkipList 基础使用
```go
// 创建城市温度记录跳表
cityComp := func(a, b string) int {
    return strings.Compare(a, b)
}
tempList := list.NewSkipList[string](cityComp)

// 插入数据
tempList.Insert("Beijing")
tempList.Insert("Shanghai")

// 查询数据
if tempList.Search("Beijing") {
    fmt.Println("Beijing exists")
}
```

### 5.2 ConcurrentList 并发场景
```go
// 创建线程安全的计数器列表
counterList := &list.ConcurrentList[int]{
    List: list.NewArrayList[int](),
}

// 启动10个协程并发增加计数器
for i := 0; i < 10; i++ {
    go func() {
        counterList.Append(1)
    }()
}
```

## 6. 实现对比指南

| 特性           | SkipList       | ConcurrentList | 标准List       |
|---------------|----------------|----------------|---------------|
| 线程安全        | ❌             | ✅              | ❌            |
| 排序支持        | ✅              | 依赖底层实现       | 依赖底层实现      |
| 随机访问性能     | O(log n)       | O(n)           | 依赖底层实现      |
| 内存占用        | 较高            | 低              | 低             |
| 适用场景        | 大数据量快速查找   | 多线程共享数据     | 单线程常规操作    |

## 7. 重要注意事项

### 7.1 SkipList 使用规范
1. 比较器必须实现全序关系
2. 元素类型必须实现比较器支持的类型
3. 批量操作时建议先排序后插入

### 7.2 ConcurrentList 并发控制
1. 避免在遍历过程中修改列表
```go
// 危险操作示例
cl.Range(func(i int, v int) error {
    cl.Delete(i) // 可能导致死锁
    return nil
})
```

2. 长时间持有锁的风险
```go
// 错误示例
cl.lock.Lock()
result := http.Get() // 网络请求期间持续持有锁
defer cl.lock.Unlock()
```

## 8. 性能优化建议

1. **SkipList调优**：
    - 调整跳表索引概率（需修改底层实现）
    - 批量插入数据时预排序

2. **ConcurrentList优化**：
    - 减少锁粒度（使用分段锁）
    - 使用 `sync.RWMutex` 的 `TryLock` 避免阻塞

3. **通用建议**：
    - 预估容量时初始化合适大小
    - 频繁操作时复用迭代器
    - 监控 `Len()` 与 `Cap()` 比值调整扩容策略

## 9. 常见问题排查

Q: SkipList插入重复元素如何处理？  
A: 根据比较器返回值决定，当比较器返回0时视为相同元素，覆盖原有值

Q: ConcurrentList出现死锁怎么办？  
A: 检查是否存在：
1. 锁的嵌套使用
2. 未释放锁的代码路径
3. 跨函数调用的锁传递

Q: AsSlice() 的性能影响？  
A: 每次调用都会生成新切片，大数据量时建议直接使用 Range 遍历

Q: 如何选择底层List实现？
```go
// 根据场景选择不同实现
list.NewSkipList()    // 需要排序和快速查询
list.NewArrayList()   // 随机访问频繁
list.NewLinkedList()  // 频繁插入删除
```