# 错误类型 (Error) - Error Types

## 概述
错误类型用于定义和管理应用程序中的各种错误情况，提供结构化的错误信息和错误处理机制。

## 分类

### 5.1 基础错误类型
```go
// 基础错误结构
type BaseError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Details string `json:"details,omitempty"`
    Cause   error  `json:"-"`
}

func (e *BaseError) Error() string {
    if e.Details != "" {
        return fmt.Sprintf("%s: %s (%s)", e.Code, e.Message, e.Details)
    }
    return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *BaseError) Unwrap() error {
    return e.Cause
}

// 创建基础错误
func NewBaseError(code, message string) *BaseError {
    return &BaseError{
        Code:    code,
        Message: message,
    }
}

// 创建带详情的错误
func NewBaseErrorWithDetails(code, message, details string) *BaseError {
    return &BaseError{
        Code:    code,
        Message: message,
        Details: details,
    }
}
```

### 5.2 业务错误类型
```go
// 同步任务错误
type SyncTaskError struct {
    *BaseError
    TaskID     string `json:"task_id"`
    TaskType   string `json:"task_type"`
    PlatformID string `json:"platform_id"`
}

func NewSyncTaskError(taskID, taskType, platformID, message string) *SyncTaskError {
    return &SyncTaskError{
        BaseError: NewBaseError("SYNC_TASK_ERROR", message),
        TaskID:    taskID,
        TaskType:  taskType,
        PlatformID: platformID,
    }
}

// 第三方API错误
type ThirdPartyAPIError struct {
    *BaseError
    ServiceName string `json:"service_name"`
    Endpoint    string `json:"endpoint"`
    StatusCode  int    `json:"status_code"`
    Response    string `json:"response,omitempty"`
}

func NewThirdPartyAPIError(serviceName, endpoint string, statusCode int, response string) *ThirdPartyAPIError {
    return &ThirdPartyAPIError{
        BaseError:  NewBaseError("THIRD_PARTY_API_ERROR", "Third party API call failed"),
        ServiceName: serviceName,
        Endpoint:   endpoint,
        StatusCode: statusCode,
        Response:   response,
    }
}

// 数据验证错误
type ValidationError struct {
    *BaseError
    Field   string `json:"field"`
    Value   interface{} `json:"value"`
    Rule    string `json:"rule"`
}

func NewValidationError(field string, value interface{}, rule string) *ValidationError {
    return &ValidationError{
        BaseError: NewBaseError("VALIDATION_ERROR", "Data validation failed"),
        Field:     field,
        Value:     value,
        Rule:      rule,
    }
}
```

### 5.3 系统错误类型
```go
// 数据库错误
type DatabaseError struct {
    *BaseError
    Operation string `json:"operation"`
    Table     string `json:"table,omitempty"`
    SQL       string `json:"sql,omitempty"`
}

func NewDatabaseError(operation, table, sql string, cause error) *DatabaseError {
    return &DatabaseError{
        BaseError: &BaseError{
            Code:    "DATABASE_ERROR",
            Message: "Database operation failed",
            Cause:   cause,
        },
        Operation: operation,
        Table:     table,
        SQL:       sql,
    }
}

// 网络错误
type NetworkError struct {
    *BaseError
    URL         string `json:"url"`
    Method      string `json:"method"`
    Timeout     bool   `json:"timeout"`
    RetryCount  int    `json:"retry_count"`
}

func NewNetworkError(url, method string, timeout bool, retryCount int, cause error) *NetworkError {
    return &NetworkError{
        BaseError: &BaseError{
            Code:    "NETWORK_ERROR",
            Message: "Network request failed",
            Cause:   cause,
        },
        URL:        url,
        Method:     method,
        Timeout:    timeout,
        RetryCount: retryCount,
    }
}

// 配置错误
type ConfigurationError struct {
    *BaseError
    ConfigKey string `json:"config_key"`
    ConfigValue interface{} `json:"config_value"`
}

func NewConfigurationError(configKey string, configValue interface{}, message string) *ConfigurationError {
    return &ConfigurationError{
        BaseError: NewBaseError("CONFIGURATION_ERROR", message),
        ConfigKey: configKey,
        ConfigValue: configValue,
    }
}
```

### 5.4 错误码定义
```go
// 错误码常量
const (
    // 通用错误码
    ErrCodeSuccess           = "SUCCESS"
    ErrCodeInternalError     = "INTERNAL_ERROR"
    ErrCodeInvalidParameter  = "INVALID_PARAMETER"
    ErrCodeNotFound          = "NOT_FOUND"
    ErrCodeUnauthorized      = "UNAUTHORIZED"
    ErrCodeForbidden         = "FORBIDDEN"
    ErrCodeTimeout           = "TIMEOUT"
    
    // 业务错误码
    ErrCodeTaskNotFound      = "TASK_NOT_FOUND"
    ErrCodeTaskAlreadyExists = "TASK_ALREADY_EXISTS"
    ErrCodeTaskInProgress    = "TASK_IN_PROGRESS"
    ErrCodeTaskFailed        = "TASK_FAILED"
    ErrCodeTaskCancelled     = "TASK_CANCELLED"
    
    // 第三方服务错误码
    ErrCodeThirdPartyUnavailable = "THIRD_PARTY_UNAVAILABLE"
    ErrCodeThirdPartyTimeout     = "THIRD_PARTY_TIMEOUT"
    ErrCodeThirdPartyAuthFailed  = "THIRD_PARTY_AUTH_FAILED"
    ErrCodeThirdPartyRateLimit   = "THIRD_PARTY_RATE_LIMIT"
    
    // 数据错误码
    ErrCodeDataNotFound      = "DATA_NOT_FOUND"
    ErrCodeDataConflict      = "DATA_CONFLICT"
    ErrCodeDataInvalid       = "DATA_INVALID"
    ErrCodeDataCorrupted     = "DATA_CORRUPTED"
)
```

### 5.5 错误处理工具
```go
// 错误包装器
type ErrorWrapper struct {
    errors []error
}

func NewErrorWrapper() *ErrorWrapper {
    return &ErrorWrapper{
        errors: make([]error, 0),
    }
}

func (ew *ErrorWrapper) Add(err error) {
    if err != nil {
        ew.errors = append(ew.errors, err)
    }
}

func (ew *ErrorWrapper) HasErrors() bool {
    return len(ew.errors) > 0
}

func (ew *ErrorWrapper) Errors() []error {
    return ew.errors
}

func (ew *ErrorWrapper) Error() string {
    if len(ew.errors) == 0 {
        return ""
    }
    
    messages := make([]string, len(ew.errors))
    for i, err := range ew.errors {
        messages[i] = err.Error()
    }
    return strings.Join(messages, "; ")
}

// 错误类型检查
func IsBaseError(err error) bool {
    _, ok := err.(*BaseError)
    return ok
}

func IsSyncTaskError(err error) bool {
    _, ok := err.(*SyncTaskError)
    return ok
}

func IsThirdPartyAPIError(err error) bool {
    _, ok := err.(*ThirdPartyAPIError)
    return ok
}

func IsValidationError(err error) bool {
    _, ok := err.(*ValidationError)
    return ok
}
```

## 特点

1. **结构化**: 提供结构化的错误信息
2. **可追踪**: 支持错误链和堆栈跟踪
3. **类型安全**: 使用类型系统区分不同错误
4. **国际化**: 支持多语言错误消息
5. **可序列化**: 支持 JSON 序列化

## 使用场景

- API 错误响应
- 日志记录
- 错误监控
- 调试信息
- 用户提示

## 最佳实践

1. **错误码**: 使用统一的错误码系统
2. **错误消息**: 提供清晰、有意义的错误消息
3. **错误分类**: 按功能模块分类错误类型
4. **错误处理**: 在适当的地方处理错误
5. **错误日志**: 记录详细的错误信息用于调试 