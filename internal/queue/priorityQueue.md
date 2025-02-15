
# PriorityQueue 技术文档

## 结构概述
```go
type PriorityQueue[T any] struct {
    compare  Gadget.Comparator[T]  // 元素比较器
    capacity int                   // 容量限制（<=0 表示无界队列）
    data     []T                   // 堆存储（索引0为占位）
}
```

## 核心特性
| 特性 | 有界队列 | 无界队列 |
|------|---------|---------|
| 容量限制 | ✅ 创建时固定 | ❌ 动态扩容 |
| 内存管理 | 固定分配 | 自动缩容 |
| 初始化容量 | capacity+1 | 默认64 |

## 方法详解

### 1. 构造方法
```go
func NewPriorityQueue[T any](capacity int, compare Comparator[T]) *PriorityQueue[T]
```
- 容量策略：
    - `capacity > 0`：预分配 `capacity+1` 容量
    - `capacity ≤ 0`：初始容量64，自动扩缩容

### 2. 入队操作
```go
func (p *PriorityQueue[T]) Enqueue(t T) error
```
**执行流程：**
1. 容量检查 → 2. 追加元素 → 3. 堆上浮调整

**堆调整伪代码：**
```python
node = last_index
while node > 1:
    parent = node // 2
    if heap[node] < heap[parent]:
        swap(node, parent)
        node = parent
    else:
        break
```

### 3. 出队操作
```go
func (p *PriorityQueue[T]) Dequeue() (T, error)
```
**执行流程：**
1. 弹出根节点 → 2. 末尾元素补位 → 3. 堆下沉调整 → 4. 缩容检查

**堆下沉示意图：**
```
[1]        // 根节点
[2] [3]    // 子节点
↓ 检查子节点
若子节点更小则交换
循环直到叶子节点
```

### 4. 缩容机制
```go
func (p *PriorityQueue[T]) shrinkIfNecessary()
```
**触发条件：**
- 仅限无界队列
- 当切片容量使用率 < 50% 时触发
- 缩容后容量不低于默认64

**内存变化示例：**
```
原数据: len=32, cap=128 → 新容量=64
原数据: len=100, cap=200 → 保持容量
```

## 复杂度分析
| 操作 | 时间复杂度 | 空间复杂度 |
|------|-----------|-----------|
| Enqueue | O(log n) | O(1) 或 O(n)扩容 |
| Dequeue | O(log n) | O(1) 可能缩容 |
| Peek    | O(1)     | O(1)       |

## 使用示例
```go
// 创建无界队列
pq := NewPriorityQueue[int](0, func(a, b int) int {
    return a - b
})

pq.Enqueue(3)  // 插入元素
val, _ := pq.Dequeue()  // 取出最小元素
```

## 注意事项
1. **非线程安全**：需自行加锁处理并发
2. **比较器要求**：必须实现全序关系
3. **内存波动**：无界队列在突发流量下可能引起GC压力
4. **索引基准**：data[0]为占位元素，实际数据从data[1]开始

# PriorityQueue 技术文档（补充更新）

## 缩容机制深度解析
```go
// 调用 Gadget 内部切片缩容方法
p.data = slice.Shrink[T](p.data)
```
**缩容策略细节：**
- 触发阈值：当切片容量 > 64 且 长度 < 容量/2 时
- 新容量计算：取 max(长度*2, 64)
- 内存回收：创建新切片并拷贝数据，触发原内存GC

**缩容执行示例：**
```
删除前: len=32, cap=100 → 删除后 len=31 → 触发缩容 → 新cap=64
持续删除至 len=30 → cap保持64（因 30 > 64/2=32 不成立）
```

## 错误处理规范
### 入队错误
```go
if p.isFull() {
    return ErrOutOfCapacity
}
```
- 有界队列满容量时返回错误代码示例：
```go
boundedQueue := NewPriorityQueue(5, comparator)
// 第六次入队时返回 ErrOutOfCapacity
```

### 出队错误
```go
if p.isEmpty() {
    return t, ErrEmptyQueue
}
```
- 空队列保护机制：
```go
if len(p.data) < 2 { // 索引0为占位元素
    return zero, ErrEmptyQueue
}
```

## 堆维护算法详解
### 上浮过程 (Enqueue)
```go
for parent > 0 && p.compare(p.data[node], p.data[parent]) < 0 {
    // 交换父子节点...
}
```
**执行路径示例：**
插入元素 → 索引5 → 父节点2 → 比较 → 交换 → 新父节点1

### 下沉过程 (heapify)
```go
func (p *PriorityQueue[T]) heapify(data []T, n, i int) {
    // 寻找最小子节点...
}
```
**三步比较逻辑：**
1. 比较左子节点 (i*2)
2. 比较右子节点 (i*2+1)
3. 选择三者最小节点交换

## 初始化策略对比表
| 队列类型 | 初始容量 | 底层数组容量 | 内存预分配策略 |
|---------|---------|-------------|---------------|
| 有界队列 | capacity | capacity+1 | 精确分配避免扩容 |
| 无界队列 | 64      | 64          | 平衡内存与性能 |

## 完整使用示例
### 有界队列场景
```go
// 创建容量为3的有界队列
bq := NewPriorityQueue(3, func(a, b int) int {
    return a - b
})

// 填充队列
bq.Enqueue(5)  // [5]
bq.Enqueue(3)  // [3,5]
bq.Enqueue(8)  // [3,5,8]

// 尝试超容插入
err := bq.Enqueue(2)  // 返回 ErrOutOfCapacity

// 正确出队流程
val, _ := bq.Dequeue()  // 3 → 队列变为 [5,8]
```

## 内存布局演进图
```
无界队列初始状态:
[ nil ] → len=1, cap=64

插入50个元素后:
[ nil, 1,2...50 ] → len=51, cap=64

删除40个元素后:
[ nil, 41...50 ] → len=11 → 触发缩容 → cap=64（因 11 > 64/2=32 不成立）
```