# Internal Data Layer

## 职责说明

`internal/data` 层负责数据基础设施，提供数据库连接、缓存连接等基础服务。

### 主要职责

1. **数据库连接管理**
   - 多数据库连接池管理
   - 数据库健康检查
   - 连接生命周期管理

2. **Redis 连接管理**
   - Redis 客户端管理
   - 连接池配置

3. **数据模型定义**
   - 数据库表结构模型 (`models/`)
   - 数据验证规则

4. **数据库迁移**
   - Schema 版本管理 (`migrations/`)
   - 数据库初始化

5. **基础设施服务**
   - 连接池监控
   - 性能指标收集
   - 资源清理

### 文件说明

- `data.go` - 核心数据层结构，管理所有数据库连接
- `database_factory.go` - 数据库工厂，创建不同类型的数据库连接
- `database_init.go` - 数据库初始化逻辑
- `redis.go` - Redis 连接管理
- `models/` - 数据模型定义
- `migrations/` - 数据库迁移文件

### 使用方式

```go
// 获取数据库连接
db, err := data.GetSyncDB()

// 获取 Redis 客户端
redis := data.GetRedis()

// 健康检查
status := data.HealthCheck(ctx)
```

### 与 Repository 层的关系

- `internal/data` 提供基础设施（连接、模型）
- `internal/repository` 提供业务数据访问逻辑
- Repository 层依赖 Data 层获取数据库连接
