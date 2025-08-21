# API 使用示例

本目录包含了 Nancalacc API 的各种使用示例，帮助您快速上手和测试 API 功能。

## 📁 目录结构

```
examples/
├── curl/                    # curl 命令行示例
│   ├── account-api.sh      # 账户管理API测试脚本
│   └── README.md           # curl使用说明
├── postman/                # Postman 测试集合
│   ├── nancalacc-api.postman_collection.json    # API测试集合
│   ├── nancalacc-local.postman_environment.json # 本地环境配置
│   └── README.md           # Postman使用说明
└── README.md               # 本文件
```

## 🚀 快速开始

### 1. 使用 curl 测试

```bash
# 进入curl示例目录
cd docs/api/examples/curl

# 给脚本添加执行权限
chmod +x account-api.sh

# 检查服务状态
./account-api.sh check

# 创建账户
./account-api.sh create test-account test@example.com

# 获取账户列表
./account-api.sh list
```

### 2. 使用 Postman 测试

1. 下载 Postman 应用
2. 导入测试集合：`docs/api/examples/postman/nancalacc-api.postman_collection.json`
3. 导入环境配置：`docs/api/examples/postman/nancalacc-local.postman_environment.json`
4. 选择 "Nancalacc Local" 环境
5. 开始测试 API

## 📋 环境配置

### 环境变量

在使用示例之前，请确保设置正确的环境变量：

```bash
# 设置基础URL
export BASE_URL="http://localhost:8000"

# 设置认证令牌（可选）
export TOKEN="your-access-token"

# 设置API版本
export API_VERSION="v1"
```

### 服务启动

确保 Nancalacc 服务正在运行：

```bash
# 启动服务
make run

# 或者使用 Docker
docker-compose up -d
```

## 🔧 示例详解

### curl 脚本功能

`account-api.sh` 脚本提供以下功能：

- **check**: 检查服务健康状态
- **create**: 创建新的同步账户
- **list**: 获取账户列表
- **get**: 获取单个账户详情
- **update**: 更新账户信息
- **delete**: 删除账户
- **sync**: 启动同步任务
- **status**: 获取同步状态

### Postman 集合

Postman 集合包含以下测试用例：

1. **Health Check**: 服务健康检查
2. **Account Management**: 账户管理操作
3. **Sync Operations**: 同步任务操作
4. **Authentication**: 认证相关操作

## 🛠 自定义配置

### 修改基础URL

```bash
# 修改脚本中的基础URL
sed -i 's/BASE_URL="http:\/\/localhost:8000"/BASE_URL="https:\/\/api.nancalacc.com"/' account-api.sh
```

### 添加认证

```bash
# 设置认证令牌
export TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# 使用脚本时自动包含认证
./account-api.sh create my-account my@email.com
```

### 自定义请求头

```bash
# 添加自定义请求头
curl -X POST "$BASE_URL/v1/account" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -H "X-Custom-Header: custom-value" \
  -d '{"name": "test", "email": "test@example.com"}'
```

## 📊 测试结果

### 成功响应示例

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "acc_1234567890",
    "name": "test-account",
    "email": "test@example.com",
    "status": "active",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

### 错误响应示例

```json
{
  "code": 20001,
  "message": "账户已存在",
  "details": {
    "field": "name",
    "reason": "账户名称已存在"
  },
  "request_id": "req_1234567890"
}
```

## 🔍 调试技巧

### 1. 启用详细输出

```bash
# curl 详细输出
curl -v -X GET "$BASE_URL/v1/accounts"

# 脚本调试模式
bash -x ./account-api.sh list
```

### 2. 检查网络连接

```bash
# 检查端口是否开放
telnet localhost 8000

# 检查DNS解析
nslookup api.nancalacc.com
```

### 3. 查看日志

```bash
# 查看服务日志
docker logs nancalacc-service

# 查看应用日志
tail -f logs/app.log
```

## 🚨 常见问题

### Q: 连接被拒绝
**A**: 检查服务是否启动，端口是否正确

```bash
# 检查服务状态
./account-api.sh check

# 检查端口
netstat -tlnp | grep 8000
```

### Q: 认证失败
**A**: 检查令牌是否有效

```bash
# 重新获取令牌
curl -X POST "$BASE_URL/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "password"}'
```

### Q: 参数错误
**A**: 检查请求参数格式

```bash
# 查看API文档
open http://localhost:8080/docs

# 检查参数格式
./account-api.sh create "" ""
```

## 📚 更多资源

- [API 文档](../swagger/) - 完整的API文档
- [错误码说明](../ERROR_CODES.md) - 错误码详细说明
- [开发指南](../../development.md) - 开发环境配置
- [部署指南](../../deployment.md) - 生产环境部署

## 🤝 贡献

欢迎提交新的示例和改进建议！

1. Fork 项目
2. 创建功能分支
3. 添加示例代码
4. 提交 Pull Request

---

**最后更新**: 2024年1月1日  
**维护者**: Nancalacc Team 