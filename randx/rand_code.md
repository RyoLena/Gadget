# 随机字符串生成工具文档（randx包）

## 功能概述
本工具提供安全的随机字符串生成能力，支持：
- 🎲 组合数字/字母/符号多种字符类型
- ⚡ 高性能随机生成（优化随机位利用率）
- 🛡️ 完善的参数校验机制
- 🔧 支持自定义字符集

## 核心类型说明

### 字符类型枚举
```go
TypeDigit       // 数字 0-9
TypeLowerCase   // 小写字母 a-z
TypeUpperCase   // 大写字母 A-Z
TypeSpecial     // 特殊符号 ~!@#$%^等
TypeMixed       // 全部字符类型组合
```

## 主要API接口

### 1. RandCode - 按类型生成
```go
// 参数：
// length: 字符串长度 (必须 >=0)
// typ: 字符类型组合（如 TypeDigit|TypeUpperCase）
func RandCode(length int, typ Type) (string, error)
```

### 2. RandStrByCharset - 自定义字符集
```go
// 参数：
// charset: 自定义字符集合（如 "ABC123"）
func RandStrByCharset(length int, charset string) (string, error)
```

## 使用示例

### 基本用法
```go
// 生成6位数字验证码
code, _ := randx.RandCode(6, randx.TypeDigit) 

// 生成8位包含大小写字母的密码
pwd, _ := randx.RandCode(8, randx.TypeLowerCase|randx.TypeUpperCase)

// 自定义字符集（生成10位十六进制字符串）
hexStr, _ := randx.RandStrByCharset(10, "0123456789abcdef")
```

## 实现原理

### 高效随机生成流程
1. 计算最优位掩码（getFirstMask）
2. 预生成63位随机数（rand.Int63）
3. 分段使用随机位：
   ```text
   |<- 63位随机数 ->|
   [6位][6位][6位]... (共63/6=10次使用)
   ```
4. 每个随机位段选择字符集索引

### 性能优化点
- 减少随机数生成调用次数
- 位运算替代取模运算
- 预计算掩码和重用随机位

## 错误处理

| 错误类型                   | 触发条件                     | 处理建议                 |
|--------------------------|----------------------------|------------------------|
| errLengthLessThanZero    | length < 0                 | 检查输入参数合法性          |
| errTypeNotSupported      | 类型组合无效/字符集为空        | 确保至少选择一种有效字符类型   |

## 注意事项
1. 字符集选择建议：
    - 验证码：仅数字（TypeDigit）
    - 密码：混合类型（TypeMixed）
    - Token：增加特殊符号（TypeSpecial）

2. 性能提示：
    - 批量生成时建议复用RandCode实例
    - 超长字符串（>1000字符）建议分批次生成

3. 安全提醒：
    - 不要用于加密场景（使用crypto/rand替代）
    - 重要凭证建议添加有效期限制

## 常见问题

Q：生成的字符串为什么有重复字符？  
A：这是正常现象，随机算法允许字符重复出现

Q：如何排除容易混淆的字符（如0/O）？  
A：使用自定义字符集：
```go
cleanCharset := "23456789ABCDEFGHJKLMNPQRSTUVWXYZ"
RandStrByCharset(8, cleanCharset)
```

Q：为什么有时返回空字符串？  
A：当length=0时会返回空字符串，这是预期行为

Q：如何确保线程安全？  
A：所有函数均已实现并发安全，可放心在多goroutine中使用
