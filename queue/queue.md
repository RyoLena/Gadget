# 并发队列技术文档

## 1. 核心队列概览

| 队列类型                     | 特点                          | 适用场景                  |
|----------------------------|-----------------------------|-------------------------|
| ConcurrentLinkedBlockingQueue | 链表实现阻塞队列（有界/无界）         | 通用并发队列场景              |
| ConcurrentArrayBlockingQueue  | 数组实现阻塞队列（有界）            | 固定容量队列场景              |
| PriorityQueue               | 优先级队列                     | 需要按优先级处理元素的场景         |
| DelayQueue                  | 延迟队列                      | 定时任务/延迟处理场景           |
| ConcurrentLinkedQueue       | 无锁链表队列                   | 超高并发非阻塞场景             |

## 2. 核心接口说明

### BlockingQueue 接口
```go
type BlockingQueue[T any] interface {
    Enqueue(ctx context.Context, t T) error  // 阻塞入队
    Dequeue(ctx context.Context) (T, error)  // 阻塞出队
}
```

### 通用方法
```go
Len() int       // 获取队列长度
AsSlice() []T   // 转换为切片
```

## 3. 主要队列实现

### 3.1 ConcurrentLinkedBlockingQueue（链表阻塞队列）
```go
// 创建实例（capacity<=0为无界队列）
q := NewConcurrentLinkedBlockingQueue[string](100)

// 入队
err := q.Enqueue(context.Background(), "task1")

// 出队
task, err := q.Dequeue(context.Background())
```

**特点**：
- 基于链表实现
- 读写使用互斥锁保护
- 条件变量实现阻塞等待

### 3.2 ConcurrentArrayBlockingQueue（数组阻塞队列）
```go
// 创建固定容量队列
q := NewConcurrentArrayBlockingQueue[int](1000)

// 入队（自动处理环形数组）
err := q.Enqueue(ctx, 42)

// 出队
num, err := q.Dequeue(ctx)
```

**特点**：
- 基于环形数组实现
- 使用信号量控制容量
- 自动处理数组越界

### 3.3 PriorityQueue（优先级队列）
```go
// 创建优先级队列（需提供比较器）
comp := func(a, b int) int { return a - b }
pq := NewPriorityQueue[int](100, comp)

// 入队（自动排序）
pq.Enqueue(5)  
pq.Enqueue(3)

// 出队顺序为 3 → 5
val, _ := pq.Dequeue()
```

**特点**：
- 基于堆实现
- 元素按比较器顺序出队
- 非阻塞操作

### 3.4 DelayQueue（延迟队列）
```go
// 定义延迟元素类型
type DelayItem struct {
    time time.Time
}
func (d DelayItem) Delay() time.Duration {
    return time.Until(d.time)
}

// 使用示例
dq := NewDelayQueue[DelayItem](100)
dq.Enqueue(ctx, DelayItem{time.Now().Add(1*time.Hour)})

// 到期自动出队
item, _ := dq.Dequeue(ctx)
```

**特点**：
- 元素需实现Delayable接口
- 使用优先级队列管理到期时间
- 精确到毫秒级的延迟处理

## 4. 实现对比

| 特性               | LinkedBlocking | ArrayBlocking | Priority | Delay   | LinkedQueue |
|--------------------|----------------|---------------|----------|---------|-------------|
| 线程安全            | ✅             | ✅            | ✅       | ✅      | ✅（无锁）    |
| 阻塞操作            | ✅             | ✅            | ❌       | ✅      | ❌          |
| 容量限制            | 可选           | 固定           | 可选      | 可选     | 无界         |
| 排序特性            | FIFO          | FIFO          | 自定义排序 | 按到期时间 | FIFO        |
| 时间复杂度（入队/出队） | O(1)          | O(1)          | O(log n) | O(log n)| O(1)        |

## 5. 使用示例

### 5.1 生产者-消费者模式
```go
q := NewConcurrentLinkedBlockingQueue[Task](100)

// 生产者
go func() {
    for {
        task := generateTask()
        if err := q.Enqueue(ctx, task); err != nil {
            // 处理错误
        }
    }
}()

// 消费者
go func() {
    for {
        task, err := q.Dequeue(ctx)
        if err != nil {
            // 处理错误
        }
        processTask(task)
    }
}()
```

### 5.2 延迟任务调度
```go
dq := NewDelayQueue[DelayTask](100)

// 调度任务
dq.Enqueue(ctx, DelayTask{
    ExecuteTime: time.Now().Add(5*time.Minute),
    Job: func(){ /* ... */ },
})

// 执行线程
go func() {
    for {
        task, _ := dq.Dequeue(ctx)
        task.Job()
    }
}()
```

## 6. 重要注意事项

1. **死锁风险**：
   ```go
   // 错误示例：同一goroutine连续操作
   q.Enqueue(ctx, item)
   item, _ = q.Dequeue(ctx) // 可能导致死锁（取决于队列容量）
   ```

2. **上下文传播**：
   ```go
   // 应该为每个操作创建子上下文
   ctx, cancel := context.WithTimeout(context.Background(), time.Second)
   defer cancel()
   q.Enqueue(ctx, item)
   ```

3. **资源清理**：
   ```go
   // 退出时清空队列
   for q.Len() > 0 {
       q.Dequeue(context.Background())
   }
   ```

4. **性能调优**：
    - 根据场景选择合适的队列类型
    - 监控队列长度指标
    - 设置合理的队列容量

## 7. 常见问题

Q：如何选择阻塞队列类型？  
A：根据场景需求：
- 需要严格FIFO → LinkedBlocking/ArrayBlocking
- 需要优先级处理 → PriorityQueue
- 需要延迟执行 → DelayQueue
- 超高并发非阻塞 → ConcurrentLinkedQueue

Q：为什么Dequeue返回error？  
A：需要处理：
1. 上下文取消/超时
2. 队列已关闭
3. 内部系统错误

Q：无界队列的内存风险？  
A：建议：
1. 监控队列长度
2. 设置合理容量限制
3. 使用背压机制控制生产者速度

Q：如何实现队列监控？  
A：可通过装饰器模式扩展：
```go
type MonitoredQueue[T any] struct {
    q       BlockingQueue[T]
    metrics Metrics
}

func (m *MonitoredQueue[T]) Enqueue(ctx context.Context, t T) error {
    start := time.Now()
    defer func() {
        m.metrics.Record(time.Since(start))
    }()
    return m.q.Enqueue(ctx, t)
}
```