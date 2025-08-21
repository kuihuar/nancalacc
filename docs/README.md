# Nancalacc 项目文档

欢迎来到 Nancalacc 项目文档中心！这里包含了项目的完整文档，帮助您快速了解和使用我们的系统。

## 📚 文档目录

### 🚀 快速开始
- [项目概述](./overview.md) - 项目介绍和核心功能
- [安装指南](./installation.md) - 环境搭建和安装步骤
- [快速上手](./quickstart.md) - 5分钟快速体验

### 📖 API 文档
- [API 生成指南](./api/API_DOCUMENTATION_GUIDE.md) - API文档生成方案
- [错误码说明](./api/ERROR_CODES.md) - 完整的错误码列表
- [交互式API文档](./api/swagger/) - Swagger UI界面

### 🔧 开发指南
- [架构设计](./architecture.md) - 系统架构和设计理念
- [开发环境](./development.md) - 开发环境配置
- [代码规范](./coding-standards.md) - 编码规范和最佳实践
- [测试指南](./testing.md) - 单元测试和集成测试

### 📋 使用示例
- [curl 示例](./api/examples/curl/) - 命令行API测试
- [Postman 集合](./api/examples/postman/) - Postman测试集合

### 🔍 运维指南
- [部署指南](./deployment.md) - 生产环境部署
- [监控告警](./monitoring.md) - 系统监控和告警配置
- [故障排查](./troubleshooting.md) - 常见问题和解决方案

### 📈 性能优化
- [性能基准](./performance.md) - 性能测试结果
- [优化建议](./optimization.md) - 性能优化指南
- [缓存策略](./caching.md) - 缓存配置和使用

## 🎯 核心功能

### 账户管理
- 账户创建、更新、删除
- 账户状态管理
- 权限控制

### 数据同步
- 多源数据同步
- 增量同步支持
- 同步状态监控

### 文件处理
- 多种格式支持
- 批量处理
- 进度跟踪

## 🛠 技术栈

- **后端**: Go + Kratos框架
- **数据库**: PostgreSQL + Redis
- **消息队列**: RabbitMQ
- **监控**: Prometheus + Grafana
- **文档**: OpenAPI 3.0 + Swagger UI

## 📊 项目状态

| 模块 | 状态 | 完成度 | 文档 |
|------|------|--------|------|
| 账户管理 | ✅ 完成 | 100% | ✅ |
| 数据同步 | ✅ 完成 | 100% | ✅ |
| 文件处理 | ✅ 完成 | 100% | ✅ |
| API文档 | ✅ 完成 | 100% | ✅ |
| 监控告警 | 🔄 进行中 | 80% | 🔄 |
| 性能优化 | 🔄 进行中 | 70% | 🔄 |

## 🚀 快速体验

### 1. 启动服务

```bash
# 克隆项目
git clone https://github.com/your-org/nancalacc.git
cd nancalacc

# 安装依赖
make init

# 启动服务
make run
```

### 2. 访问API文档

```bash
# 生成API文档
make generate-docs

# 启动文档服务器
make serve-docs
```

访问 http://localhost:8080/docs 查看交互式API文档。

### 3. 测试API

```bash
# 使用curl测试
chmod +x docs/api/examples/curl/account-api.sh
./docs/api/examples/curl/account-api.sh create test-account test@example.com
```

## 📞 获取帮助

### 文档问题
- 📖 查看 [FAQ](./faq.md)
- 🔍 搜索文档
- 📝 提交文档改进建议

### 技术问题
- 🐛 [提交Issue](https://github.com/your-org/nancalacc/issues)
- 💬 [讨论区](https://github.com/your-org/nancalacc/discussions)
- 📧 技术支持: support@nancalacc.com

### 社区
- 🌐 [官方网站](https://nancalacc.com)
- 📱 [微信公众号](https://mp.weixin.qq.com/nancalacc)
- 🐦 [Twitter](https://twitter.com/nancalacc)

## 📄 许可证

本项目采用 [MIT License](../LICENSE) 许可证。

## 🤝 贡献指南

我们欢迎所有形式的贡献！请查看 [贡献指南](./contributing.md) 了解如何参与项目开发。

### 贡献类型
- 🐛 Bug修复
- ✨ 新功能开发
- 📖 文档改进
- 🧪 测试用例
- 🔧 工具改进

### 贡献流程
1. Fork 项目
2. 创建功能分支
3. 提交更改
4. 创建 Pull Request

## 📈 项目统计

![GitHub stars](https://img.shields.io/github/stars/your-org/nancalacc)
![GitHub forks](https://img.shields.io/github/forks/your-org/nancalacc)
![GitHub issues](https://img.shields.io/github/issues/your-org/nancalacc)
![GitHub pull requests](https://img.shields.io/github/issues-pr/your-org/nancalacc)

---

**最后更新**: 2024年1月1日  
**文档版本**: v1.0.0  
**维护者**: Nancalacc Team 