# 切片删除与缩容实现文档

## 文件结构
```go
// delete.go
package slices

import "github.com/pkg/errors"

// Delete 泛型切片删除函数
func Delete[T any](src []T, index int) ([]T, T, error) {
    // 实现细节...
}
```

## 关键逻辑说明

### 1. 索引验证
```go
if index < 0 || index >= length {
    return nil, zero, errors.New("索引越界")
}
```
- 校验规则：`index ∈ [0, len(src)-1]`
- 错误处理：返回包含原始长度和非法索引的错误

### 2. 元素迁移
```go
for i := index; i < length-1; i++ {
    src[i] = src[i+1]  // 前移覆盖操作
}
```
- 操作范围：`[index, len(src)-2]`
- 时间复杂度：O(n-index)

### 3. 缩容机制（关键部分）
```go
// 原始截断操作
src = src[:length-1]

// 新增缩容判断（示例实现）
const shrinkThreshold = 64  // 缩容阈值

func shouldShrink(oldCap, newLen int) bool {
    return oldCap > shrinkThreshold && newLen < oldCap/2
}

if shouldShrink(cap(src), len(src)) {
    newSlice := make([]T, len(src))
    copy(newSlice, src)
    src = newSlice
}
```

#### 缩容策略说明表
| 条件 | 操作 | 目的 |
|------|------|------|
| 原始容量 ≤ 64 | 保持容量 | 避免小切片频繁分配 |
| 新长度 < 原容量/2 | 创建新切片 | 释放未使用内存 |
| 元素大小 > 1KB | 更激进缩容 | 优化大对象内存 |

### 4. 返回值特征
```go
return src, res, nil
```
| 返回值 | 类型 | 说明 |
|--------|------|------|
| 第一返回值 | []T | **可能指向新内存空间**的切片 |
| 第二返回值 | T | 被删除元素的原始值 |
| 错误对象 | error | 成功时为nil |

## 内存变化示意图
```
原始内存布局
[ ][ ][X][ ][ ]  len=5, cap=5

删除索引2后
[ ][ ][ ][ ]     len=4, cap=5

触发缩容后
[ ][ ][ ][ ]     len=4, cap=4 (新数组)
```