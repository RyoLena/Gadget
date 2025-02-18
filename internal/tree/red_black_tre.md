# 红黑树技术文档

## 基本概念
红黑树就像自动保持平衡的"智能字典"：
- 每个节点有红/黑两种颜色（类似交通信号灯）
- 自动调整节点颜色和位置保持平衡
- 查找速度始终很快（O(log n)）

## 数据结构图解
```go
type RBTree[K any, V any] struct {
    root    *rbNode[K, V]  // 根节点
    compare Comparator[K]  // 比较器（决定节点顺序）
    size    int            // 节点总数
}

type rbNode[K any, V any] struct {
    color  color  // 节点颜色
    key    K      // 键（用于查找）
    value  V      // 值（存储数据）
    left   *rbNode[K, V]
    right  *rbNode[K, V]
    parent *rbNode[K, V]
}
```

## 核心方法速查表
| 方法 | 作用 | 使用示例 |
|------|------|---------|
| Add() | 添加键值对 | tree.Add(123, "苹果") |
| Delete() | 删除键 | tree.Delete(123) |
| Find() | 查找值 | val, _ := tree.Find(123) |
| Set() | 修改值 | tree.Set(123, "红苹果") |
| KeyValues() | 获取所有数据 | keys, vals := tree.KeyValues() |

## 重要规则说明
1. **颜色规则**：
    - 根节点必须是黑色
    - 红色节点的子节点必须是黑色
    - 所有路径的黑色节点数量相同

2. **自动平衡机制**：
```go
// 插入后调整（示例）
func (rb *RBTree[K, V]) fixAfterAdd(x *rbNode[K, V]) {
    // 根据叔叔节点颜色调整
    // 可能进行旋转操作
}

// 删除后调整（示例）
func (rb *RBTree[K, V]) fixAfterDelete(x *rbNode[K, V]) {
    // 根据兄弟节点颜色调整
    // 可能进行旋转操作
}
```

## 使用示例
```go
// 创建红黑树（需提供比较函数）
tree := tree.NewRBTree[int, string](func(a, b int) int {
    return a - b
})

// 添加数据
tree.Add(100, "苹果")
tree.Add(200, "香蕉")
tree.Add(150, "橘子")

// 查找数据
value, _ := tree.Find(150) // 得到"橘子"

// 修改数据
tree.Set(100, "红富士苹果")

// 删除数据
deletedValue, _ := tree.Delete(200) // 删除"香蕉"
```

## 常见错误处理
| 错误 | 原因 | 解决方法 |
|------|------|---------|
| ErrRBTreeSameRBNode | 添加重复的键 | 检查键是否已存在 |
| ErrRBTreeNotRBNode | 操作不存在的键 | 先用Find()检查存在性 |

## 性能特点
| 操作 | 时间复杂度 | 说明 |
|------|-----------|------|
| 插入 | O(log n) | 自动平衡保证效率 |
| 删除 | O(log n) | 最多三次旋转调整 |
| 查找 | O(log n) | 比普通二叉树更快 |

## 注意事项
1. **比较器要求**：必须正确实现比较逻辑
2. **非线程安全**：并发操作需要加锁
3. **内存占用**：每个节点存储三个指针（左右父）
4. **适用场景**：适合需要快速查找和范围查询的场景

## 学习小贴士
- 插入平衡的三种情况：
  1️⃣ 叔叔红 → 重新着色  
  2️⃣ 叔叔黑且当前节点是左孩子 → 右旋转  
  3️⃣ 叔叔黑且当前节点是右孩子 → 左旋转
- 删除时优先处理黑色节点的情况
