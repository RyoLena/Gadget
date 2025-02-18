# 并发工具包技术文档

## 1. 对象池工具 (Pool)

### 1.1 核心功能
🔄 提供安全的对象复用机制，相比原生`sync.Pool`：
- ✅ 类型安全（无需类型断言）
- 🛡️ 防止意外存储错误类型
- 🏗️ 泛型支持

### 1.2 基本使用
```go
type Conn struct{ /* 数据库连接 */ }

// 创建连接池
pool := syncx.NewPool(func() Conn {
    return NewDBConnection() 
})

// 获取连接
conn := pool.Get()

// 归还连接
defer pool.Put(conn)
```

## 2. 分段锁 (SegmentKeysLock)

### 2.1 设计目标
🔒 解决全局锁性能瓶颈：
- 🧩 将锁按key哈希分片
- ⚡ 减少锁竞争
- 📈 提升并发性能

### 2.2 使用示例
```go
lock := syncx.NewSegmentKeysLock(32) // 32个分片

func updateConfig(key string) {
    lock.Lock(key)
    defer lock.Unlock(key)
    // 安全更新配置
}
```

## 3. 限流对象池 (LimitPool)

### 3.1 工作机制
🪙 令牌控制策略：
- ✅ Get() 消耗令牌
- 🔄 Put() 返还令牌
- 🚫 令牌耗尽时返回新对象

### 3.2 使用场景
```go
// 限制最多缓存100个大对象
bigObjPool := syncx.NewLimitPool(100, func() BigObj {
    return BigObj{Data: make([]byte, 10MB)}
})

obj, fromPool := bigObjPool.Get()
if fromPool {
    defer bigObjPool.Put(obj)
}
```

## 4. 增强型条件变量 (Cond)

### 4.1 核心改进
⏰ 相比标准库`sync.Cond`：
- 🕒 支持上下文超时
- 📡 更安全的广播机制
- 🚦 防止条件变量被复制

### 4.2 典型使用模式
```go
var (
    mu  sync.Mutex
    cond = syncx.NewCond(&mu)
    dataReady bool
)

// 等待方
go func() {
    mu.Lock()
    defer mu.Unlock()
    
    for !dataReady {
        if err := cond.Wait(context.TODO()); err != nil {
            // 处理超时
        }
    }
    // 处理数据
}()

// 通知方
mu.Lock()
dataReady = true
cond.Broadcast()
mu.Unlock()
```

## 5. 重要注意事项

### 5.1 对象池
1. 存储对象应保持"干净"状态
2. 大对象优先考虑LimitPool
3. 不要存储携带goroutine状态的对象

### 5.2 分段锁
1. 分片数应设为2的幂次
2. 不同key可能哈希到同一分片
3. 读多写少场景使用RLock/RUnlock

### 5.3 条件变量
1. 必须在持有锁时调用Wait()
2. 使用模式必须保持：
   ```go
   mu.Lock()
   defer mu.Unlock()
   for condition {
       if err := cond.Wait(ctx); err != nil {
           // ...
       }
   }
   ```
3. Broadcast可能唤醒多个等待者

## 6. 性能优化建议

| 工具         | 适用场景                      | 优化要点                     |
|------------|---------------------------|--------------------------|
| SegmentLock | 高频不同key操作               | 增加分片数（32-256）          |
| LimitPool   | 内存敏感的大对象缓存              | 合理设置maxTokens            |
| Cond        | 复杂条件等待场景                | 优先使用channel替代          |
| Pool        | 高频创建/销毁的小对象             | 确保Put前重置对象状态           |

## 7. 常见问题排查

Q: LimitPool获取不到对象？
A: 检查：
1. maxTokens是否过小
2. Put是否正常调用
3. 对象创建函数是否耗时过长

Q: 条件变量等待超时？
A: 检查：
1. 是否在持有锁时调用Wait()
2. 通知方是否正确调用Signal/Broadcast
3. 上下文超时设置是否合理

Q: 分段锁性能未达预期？
A: 尝试：
1. 更换哈希算法
2. 增加分片数量
3. 检查key分布是否均匀
