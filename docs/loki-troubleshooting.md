# Loki 故障排查指南

## 问题：查询不到推送的日志数据

### 1. 检查 Loki 服务状态

```bash
# 检查 Loki 是否就绪
curl -s http://192.168.1.142:3100/ready

# 检查 Loki 配置
curl -s http://192.168.1.142:3100/config
```

### 2. 验证推送是否成功

```bash
# 测试推送
curl -X POST http://192.168.1.142:3100/loki/api/v1/push \
  -H "Content-Type: application/json" \
  -d '{
    "streams": [
      {
        "stream": {
          "service": "nancalacc",
          "level": "info"
        },
        "values": [
          ["'$(date +%s)000000000'", "{\"level\":\"info\",\"msg\":\"测试日志\"}"]
        ]
      }
    ]
  }'
```

### 3. 检查标签

```bash
# 查看所有标签
curl -s "http://192.168.1.142:3100/loki/api/v1/labels" | jq .

# 查看特定标签的值
curl -s "http://192.168.1.142:3100/loki/api/v1/label/service/values" | jq .
```

### 4. 查询日志

```bash
# 使用正确的时间范围查询（注意时区）
curl -s "http://192.168.1.142:3100/loki/api/v1/query_range?query={service=\"nancalacc\"}&start=2025-08-16T07:00:00Z&end=2025-08-16T08:00:00Z" | jq .

# 查询所有日志（使用正则表达式）
curl -s "http://192.168.1.142:3100/loki/api/v1/query_range?query={service=~\".+\"}&start=2025-08-16T07:00:00Z&end=2025-08-16T08:00:00Z" | jq .
```

### 5. 常见问题

#### 5.1 时间范围问题
- Loki 查询时间范围限制：30天1小时
- 确保查询时间范围正确
- 注意时区转换

#### 5.2 数据持久化延迟
- 新推送的数据可能还在内存中
- 等待几分钟让数据被持久化
- 检查 Loki 的存储配置

#### 5.3 标签匹配问题
- 确保查询的标签名称正确
- 使用正确的标签值
- 检查标签是否被正确创建

### 6. 调试步骤

1. **推送测试**：使用测试脚本验证推送功能
2. **标签检查**：确认标签已创建
3. **时间范围**：使用正确的时间范围查询
4. **等待持久化**：给数据一些时间被写入存储
5. **检查配置**：确认 Loki 配置正确

### 7. 监控和告警

```logql
# 监控日志量
sum(rate({service="nancalacc"}[5m]))

# 监控错误率
sum(rate({service="nancalacc", level="error"}[5m])) / sum(rate({service="nancalacc"}[5m]))
```

### 8. 最佳实践

1. **结构化日志**：使用 JSON 格式
2. **合理标签**：避免过多标签
3. **时间戳**：使用正确的时间戳格式
4. **批量推送**：使用批量推送提高性能
5. **错误处理**：处理推送失败的情况 