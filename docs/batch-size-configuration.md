# 批量大小配置化改进

## 概述

本次优化将原本硬编码的批量大小（500）改为可配置的方式，提高了系统的灵活性和可维护性。

## 改进内容

### 1. 配置结构扩展

在 `internal/conf/conf.proto` 中扩展了 `App` 配置结构：

```protobuf
message App {
  // ... 原有字段 ...
  
  // 批量处理配置
  int32 batch_size = 13;           // 批量大小，默认500
  int32 max_workers = 14;          // 最大并发工作数，默认3
  string log_level = 15;           // 日志级别
  string logout = 16;              // 日志输出
}
```

### 2. 配置文件更新

#### 生产环境配置 (`configs/config.yaml`)
```yaml
app:
  # ... 原有配置 ...
  # 批量处理配置
  batch_size: 500      # 批量大小，可根据环境调整
  max_workers: 3       # 最大并发工作数
```

#### 开发环境配置 (`configs/config_dev.yaml`)
```yaml
app:
  # ... 原有配置 ...
  # 批量处理配置
  batch_size: 200      # 开发环境使用较小的批量大小
  max_workers: 2       # 开发环境使用较少的并发数
```

### 3. 代码结构优化

#### FullSyncUsecase 结构体扩展
```go
type FullSyncUsecase struct {
    // ... 原有字段 ...
    // 批量处理配置
    batchSize  int
    maxWorkers int
}
```

#### 构造函数优化
```go
func NewFullSyncUsecase(...) *FullSyncUsecase {
    // 获取批量处理配置，设置默认值
    batchSize := int(bizConf.GetBatchSize())
    if batchSize <= 0 {
        batchSize = 500 // 默认批量大小
    }
    
    maxWorkers := int(bizConf.GetMaxWorkers())
    if maxWorkers <= 0 {
        maxWorkers = 3 // 默认并发数
    }
    
    return &FullSyncUsecase{
        // ... 其他字段 ...
        batchSize:  batchSize,
        maxWorkers: maxWorkers,
    }
}
```

### 4. 方法优化

所有相关方法现在使用配置的批量大小：

#### transUser 方法
```go
// 使用配置的批量大小
users := make([]*models.TbLasUser, 0, uc.batchSize)

// 批量保存
if len(users) >= uc.batchSize {
    // 执行批量保存
}
```

#### transDept 方法
```go
// 使用配置的批量大小
depts := make([]*models.TbLasDepartment, 0, uc.batchSize)

// 批量保存
if len(depts) >= uc.batchSize {
    // 执行批量保存
}
```

#### transUserDept 方法
```go
// 使用配置的批量大小
deptusers := make([]*models.TbLasDepartmentUser, 0, uc.batchSize)

// 批量保存
if len(deptusers) >= uc.batchSize {
    // 执行批量保存
}
```

## 优势

### 1. 环境适配性
- **生产环境**：可以使用较大的批量大小（500-1000）以提高性能
- **开发环境**：使用较小的批量大小（200-300）以减少资源消耗
- **测试环境**：可以根据测试需求灵活调整

### 2. 性能调优
- 可以根据数据库性能调整批量大小
- 可以根据内存限制调整批量大小
- 可以根据网络延迟调整批量大小

### 3. 运维便利性
- 无需重新编译代码即可调整批量大小
- 可以通过配置热更新（如果支持）
- 便于 A/B 测试不同批量大小的效果

### 4. 代码维护性
- 消除了硬编码的魔法数字
- 提高了代码的可读性和可维护性
- 便于后续功能扩展

## 配置建议

### 生产环境
```yaml
batch_size: 500-1000    # 根据数据库性能调整
max_workers: 3-5        # 根据 CPU 核心数调整
```

### 开发环境
```yaml
batch_size: 200-300     # 较小的批量大小
max_workers: 2-3        # 较少的并发数
```

### 测试环境
```yaml
batch_size: 100-200     # 更小的批量大小便于调试
max_workers: 1-2        # 单线程或少量并发
```

## 监控建议

建议添加以下监控指标：

1. **批量处理时间**：记录每个批次的处理时间
2. **内存使用量**：监控批量处理时的内存使用情况
3. **数据库连接数**：监控批量处理时的数据库连接使用情况
4. **错误率**：监控批量处理时的错误率

## 后续优化方向

1. **动态调整**：根据系统负载动态调整批量大小
2. **自适应优化**：根据历史数据自动优化批量大小
3. **分片处理**：对于超大文件，支持分片处理
4. **进度恢复**：支持批量处理的中断和恢复 