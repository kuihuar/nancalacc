# API文档生成指南

## 概述

本项目使用Protocol Buffers和OpenAPI规范来自动生成API文档。通过构建时自动更新，确保文档与代码保持同步。

## 技术栈

- **Protocol Buffers**: 定义API接口和数据结构
- **protoc-gen-openapi**: 从.proto文件生成OpenAPI 3.0规范
- **Swagger UI**: 提供交互式API文档界面
- **Makefile**: 自动化构建和文档生成流程

## 目录结构

```
docs/
├── api/
│   ├── openapi/              # OpenAPI规范文件
│   │   ├── account-v1.yaml   # 账户服务API
│   │   └── index.yaml        # API索引
│   ├── swagger/              # Swagger UI静态文件
│   │   ├── index.html
│   │   └── swagger-ui/
│   └── examples/             # API使用示例
│       ├── curl/
│       └── postman/
├── third_party/
│   └── swagger/              # 第三方API文档
└── README.md                 # 文档索引
```

## 快速开始

### 1. 安装依赖

```bash
# 安装protoc和相关插件
make init

# 或者手动安装
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
```

### 2. 生成API文档

```bash
# 生成所有API文档
make api

# 或者单独生成
make generate-openapi
make generate-swagger-ui
```

### 3. 查看文档

```bash
# 启动文档服务器
make serve-docs

# 访问 http://localhost:8080/docs
```

## 详细配置

### Protocol Buffers配置

在`.proto`文件中添加OpenAPI注解：

```protobuf
syntax = "proto3";

package api.account.v1;

import "google/api/annotations.proto";
import "google/protobuf/timestamp.proto";

option go_package = "nancalacc/api/account/v1;v1";

service Account {
    rpc CreateSyncAccount (CreateSyncAccountRequest) returns (CreateSyncAccountReply) {
        option (google.api.http) = {
            post: "/v1/account"
            body: "*"
        };
    };
}

message CreateSyncAccountRequest {
    string name = 1;
    string email = 2;
}
```

### OpenAPI生成配置

在Makefile中配置生成选项：

```makefile
.PHONY: generate-openapi
generate-openapi:
	protoc --proto_path=./api \
	       --proto_path=./third_party \
	       --openapi_out=fq_schema_naming=true,default_response=false,allow_merge=true:./docs/api/openapi \
	       $(API_PROTO_FILES)
```

### Swagger UI配置

创建自定义的Swagger UI页面：

```html
<!DOCTYPE html>
<html>
<head>
    <title>Nancalacc API Documentation</title>
    <link rel="stylesheet" type="text/css" href="./swagger-ui/swagger-ui.css" />
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="./swagger-ui/swagger-ui-bundle.js"></script>
    <script src="./swagger-ui/swagger-ui-standalone-preset.js"></script>
    <script>
        window.onload = function() {
            const ui = SwaggerUIBundle({
                url: './openapi/account-v1.yaml',
                dom_id: '#swagger-ui',
                presets: [SwaggerUIBundle.presets.apis, SwaggerUIStandalonePreset],
                layout: "StandaloneLayout"
            });
        };
    </script>
</body>
</html>
```

## 构建集成

### CI/CD集成

在GitHub Actions中添加文档生成步骤：

```yaml
name: Generate API Docs
on:
  push:
    branches: [main]
    paths: ['api/**/*.proto']

jobs:
  generate-docs:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v3
        with:
          go-version: '1.21'
      - name: Install protoc
        run: |
          sudo apt-get update
          sudo apt-get install -y protobuf-compiler
      - name: Install Go tools
        run: make init
      - name: Generate API docs
        run: make generate-openapi
      - name: Commit and push docs
        run: |
          git config --local user.email "action@github.com"
          git config --local user.name "GitHub Action"
          git add docs/api/openapi/
          git commit -m "Update API docs" || exit 0
          git push
```

### 本地开发

在开发过程中自动生成文档：

```bash
# 监听.proto文件变化并自动生成文档
make watch-docs

# 或者使用air进行热重载
air -c .air.toml
```

## 文档增强

### 添加API描述

在`.proto`文件中添加详细描述：

```protobuf
service Account {
    // 创建同步账户
    // 用于创建新的账户同步任务
    rpc CreateSyncAccount (CreateSyncAccountRequest) returns (CreateSyncAccountReply) {
        option (google.api.http) = {
            post: "/v1/account"
            body: "*"
        };
    };
}

message CreateSyncAccountRequest {
    // 账户名称，必填字段
    string name = 1;
    
    // 邮箱地址，用于通知
    string email = 2;
}
```

### 添加示例

创建API使用示例：

```bash
# 创建示例目录
mkdir -p docs/api/examples/curl
mkdir -p docs/api/examples/postman

# 生成curl示例
cat > docs/api/examples/curl/create-account.sh << 'EOF'
#!/bin/bash
curl -X POST "http://localhost:8000/v1/account" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "test-account",
    "email": "test@example.com"
  }'
EOF
```

### 添加错误码文档

创建错误码说明文档：

```markdown
# API错误码说明

## 通用错误码

| 错误码 | 说明 | HTTP状态码 |
|--------|------|------------|
| 10001  | 参数错误 | 400 |
| 10002  | 未授权 | 401 |
| 10003  | 禁止访问 | 403 |
| 10004  | 资源不存在 | 404 |
| 10005  | 服务器内部错误 | 500 |

## 业务错误码

| 错误码 | 说明 | 解决方案 |
|--------|------|----------|
| 20001  | 账户已存在 | 检查账户名称是否重复 |
| 20002  | 同步任务进行中 | 等待当前任务完成 |
| 20003  | 文件格式不支持 | 检查文件格式 |
```

## 最佳实践

### 1. 版本管理

- 使用语义化版本控制
- 保持向后兼容性
- 废弃的API要标记并说明迁移方案

### 2. 文档维护

- 每次API变更都要更新文档
- 添加充分的示例和说明
- 定期审查文档的准确性

### 3. 测试集成

- 使用OpenAPI规范生成测试用例
- 集成API文档测试到CI/CD流程
- 确保文档示例可以正常运行

### 4. 性能考虑

- 大文件文档考虑分页加载
- 使用CDN加速文档访问
- 缓存生成的文档文件

## 故障排除

### 常见问题

1. **protoc命令未找到**
   ```bash
   # 安装protobuf编译器
   brew install protobuf  # macOS
   sudo apt-get install protobuf-compiler  # Ubuntu
   ```

2. **OpenAPI生成失败**
   ```bash
   # 检查.proto文件语法
   protoc --proto_path=./api --proto_path=./third_party --validate_out=lang=go:./api api/**/*.proto
   ```

3. **Swagger UI无法加载**
   ```bash
   # 检查文件路径和权限
   ls -la docs/api/swagger/
   chmod +r docs/api/swagger/*
   ```

### 调试技巧

- 使用`protoc --debug_out`查看生成过程
- 检查生成的OpenAPI文件语法
- 使用在线Swagger编辑器验证规范

## 扩展功能

### 1. 多语言支持

```yaml
# 支持多语言文档
info:
  title: Nancalacc API
  description: |
    # 中文
    账户管理API接口
    
    # English
    Account Management API
  version: 1.0.0
```

### 2. 认证文档

```yaml
components:
  securitySchemes:
    BearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT
security:
  - BearerAuth: []
```

### 3. 监控集成

```yaml
# 添加监控端点
paths:
  /metrics:
    get:
      summary: 获取应用指标
      tags: [Monitoring]
      responses:
        '200':
          description: Prometheus格式的指标数据
```

## 总结

通过这套完整的API文档生成方案，你可以：

1. **自动化**: 构建时自动生成最新文档
2. **标准化**: 使用OpenAPI 3.0规范
3. **可视化**: 提供交互式Swagger UI界面
4. **可维护**: 文档与代码保持同步
5. **可扩展**: 支持多语言、认证、监控等功能

这套方案将大大提高API的可用性和开发效率。 