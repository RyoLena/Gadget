## Pair

## 🔨 核心功能
提供通用的键值对操作工具集，支持泛型类型，主要用于：

1. **键值对结构封装**
2. **切片与键值对的相互转换**
3. **数据扁平化与重组**

## 🧩 核心组件

### 1. 基础结构体
```go
type Pair[K any, V any] struct {
    Key   K
    Value V
}
```
- 泛型支持：适用于任意类型组合
- 配套方法：
  - String(): 格式化输出 Key:Value
  - Split(): 分解返回原始键值

### 2、核心功能函数

🆕 构造器
```go
func NewPairs[K, V any](keys []K, values []V) ([]Pair[K, V], error)
```
- 输入验证：非nil检查 + 等长校验
- 生成等长键值对切片

🔄 转换器
```go
func SplitPairs[K, V any](pairs []Pair[K, V]) (keys []K, values []V)
```
- 逆向转换 []Pair → (keys, values)
- 自动处理nil输入

📦 扁平化
```go
func FlattenPairs[K, V any](pairs []Pair[K, V]) []any
```
- 转换结构：[Pair1, Pair2] → [K1, V1, K2, V2]
- ❗ 当前实现需注意预分配优化

🚚 重组
```go
func PackPairs[K, V any](flatPairs []any) []Pair[K, V]
```
- 要求输入顺序：严格交替Key/Value
- ⚠️ 类型安全依赖运行时断言

⚙️ 设计特点
```graph TD
    A[键值对处理] --> B(构造器)
    A --> C(拆分器)
    A --> D(扁平化)
    A --> E(重组)
    B --> F[输入验证]
    D --> G[内存优化]
    E --> H[类型安全]
```
### 3、完整的工作流

```go
pairs, _ := NewPairs([]int{1,2}, []string{"a","b"}) // 构造
flat := FlattenPairs(pairs)            // 扁平化 [1,a,2,b]
restored := PackPairs[int, string](flat) // 重组还原
```
## 🚀 进阶功能扩展

### 类型安全增强方案
```go
// 带类型校验的增强版 PackPairs
func SafePackPairs[K, V any](flat []any) ([]Pair[K, V], error) {
    if len(flat)%2 != 0 {
        return nil, errors.New("输入数组长度必须为偶数")
    }
    
    pairs := make([]Pair[K, V], len(flat)/2)
    for i := 0; i < len(flat); i += 2 {
        key, ok1 := flat[i].(K)
        val, ok2 := flat[i+1].(V)
        if !ok1 || !ok2 {
            return nil, fmt.Errorf("类型不匹配: 位置%d-%d", i, i+1)
        }
        pairs[i/2] = NewPair(key, val)
    }
    return pairs, nil
}
```
🧪 性能优化建议
FlattenPairs 优化前后对比
```go
// 优化前（存在多次内存分配）
func FlattenPairs(pairs []Pair) []any {
    var result []any
    for _, p := range pairs {
        result = append(result, p.Key)
        result = append(result, p.Value)
    }
    return result
}

// 优化后（预分配内存）
func FlattenPairs[K, V any](pairs []Pair[K, V]) []any {
    if len(pairs) == 0 {
        return nil
    }
    flat := make([]any, 0, len(pairs)*2)
    for _, p := range pairs {
        flat = append(flat, p.Key)
        flat = append(flat, p.Value)
    }
    return flat
}
```
**性能提升点**： 
- ▶️ 减少50%的内存分配次数 
- ▶️ 消除隐式的切片扩容开销
