# 虎皮椒支付 SDK 测试说明

## 测试文件结构

- `xunhupay_test.go` - 完整的测试套件
- `xunhupay_example.go` - 使用示例

## 运行测试

### 运行所有测试

```bash
go test -v
```

### 运行特定测试

```bash
# 运行签名测试
go test -v -run TestSign

# 运行响应解析测试
go test -v -run TestResponse

# 运行基准测试
go test -bench=.
```

### 运行基准测试

```bash
go test -bench=BenchmarkSign
```

## 测试覆盖范围

### 单元测试

1. **TestNewHuPi** - 测试客户端初始化

   - 验证客户端创建
   - 验证 appId 和 appSecret 正确设置

2. **TestSign** - 测试签名功能

   - 验证签名不为空
   - 验证签名长度为 32 位（MD5）

3. **TestSignConsistency** - 测试签名一致性

   - 验证相同参数产生相同签名

4. **TestSignWithDifferentParams** - 测试签名差异性

   - 验证不同参数产生不同签名

5. **TestResponseUnmarshal** - 测试响应解析

   - 验证 JSON 响应正确解析为 Response 结构体

6. **TestResponseErrorCase** - 测试错误响应

   - 验证错误响应正确解析

7. **TestExecuteParameterProcessing** - 测试参数处理
   - 验证参数处理逻辑
   - 测试签名功能

### 基准测试

- **BenchmarkSign** - 测试签名性能
  - 当前性能：约 3735 ns/op

### 集成测试（已跳过）

- **TestPayExample** - 支付示例测试
- **TestQueryExample** - 查询示例测试

> 注意：集成测试需要真实的 API 密钥，默认被跳过

## 测试数据

### 测试参数

```go
params := map[string]string{
    "version":        "1.1",
    "trade_order_id": "123456789",
    "total_fee":      "0.1",
    "title":          "测试标题",
}
```

### 测试响应

```json
{
  "openid": 123456789,
  "url_qrcode": "https://example.com/qrcode",
  "url": "https://example.com/pay",
  "errcode": 0,
  "errmsg": "success",
  "hash": "abc123hash",
  "data": {
    "open_order_id": "order123",
    "status": "pending"
  }
}
```

## 测试最佳实践

1. **隔离测试** - 每个测试都是独立的
2. **确定性** - 相同输入产生相同输出
3. **覆盖性** - 覆盖正常和异常情况
4. **性能测试** - 包含基准测试

## 添加新测试

### 添加单元测试

```go
func TestNewFeature(t *testing.T) {
    // 设置测试数据
    // 执行被测试的功能
    // 验证结果
}
```

### 添加基准测试

```go
func BenchmarkNewFeature(b *testing.B) {
    // 设置测试数据
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        // 执行被测试的功能
    }
}
```

## 故障排除

### 常见问题

1. **测试失败** - 检查测试数据和预期结果
2. **性能下降** - 运行基准测试比较性能
3. **网络错误** - 集成测试需要网络连接

### 调试技巧

1. 使用 `-v` 参数查看详细输出
2. 使用 `-run` 参数运行特定测试
3. 使用 `t.Logf()` 输出调试信息
