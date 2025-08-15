package otel

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	otellog "go.opentelemetry.io/otel/log"
)

// KratosLoggerAdapter 适配OpenTelemetry Logger到Kratos Logger
type KratosLoggerAdapter struct {
	otellogger otellog.Logger
	config     *Config
	writer     *os.File
}

// NewKratosLoggerAdapter 创建Kratos Logger适配器
func NewKratosLoggerAdapter(otellogger otellog.Logger, config *Config) log.Logger {
	adapter := &KratosLoggerAdapter{
		otellogger: otellogger,
		config:     config,
	}

	// 如果配置了文件输出，创建文件写入器
	if config != nil && config.Logs.Enabled && config.Logs.Output == "file" && config.Logs.FilePath != "" {
		// 确保日志目录存在
		logDir := filepath.Dir(config.Logs.FilePath)
		if err := os.MkdirAll(logDir, 0755); err == nil {
			// 尝试打开文件进行写入
			if file, err := os.OpenFile(config.Logs.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644); err == nil {
				adapter.writer = file
			}
		}
	}

	return adapter
}

// Log 实现Kratos Logger接口
func (a *KratosLoggerAdapter) Log(level log.Level, keyvals ...any) error {
	// 检查日志级别过滤
	if a.config != nil && a.config.Logs.Enabled {
		configLevel := strings.ToLower(a.config.Logs.Level)
		if !a.shouldLog(level, configLevel) {
			return nil
		}
	}

	// 转换日志级别
	severity := convertLevel(level)

	// 构建消息和属性
	message := ""
	attrs := make([]otellog.KeyValue, 0, len(keyvals)/2)

	for i := 0; i < len(keyvals); i += 2 {
		if i+1 < len(keyvals) {
			key, ok := keyvals[i].(string)
			if !ok {
				continue
			}
			value := keyvals[i+1]

			// 第一个字符串作为消息
			if message == "" {
				if msg, ok := value.(string); ok {
					message = msg
					continue
				}
			}

			// 其他作为属性
			attr := convertToKeyValue(key, value)
			attrs = append(attrs, attr)
		}
	}

	// 如果没有消息，使用默认消息
	if message == "" {
		message = "log message"
	}

	// 创建日志记录
	record := otellog.Record{}
	record.SetSeverity(severity)
	record.SetBody(otellog.StringValue(message))
	for _, attr := range attrs {
		record.AddAttributes(attr)
	}

	// 记录日志到OpenTelemetry
	a.otellogger.Emit(context.Background(), record)

	// 如果配置了文件输出，同时写入文件
	if a.writer != nil {
		a.writeToFile(level, message, keyvals)
	}

	return nil
}

// shouldLog 检查是否应该记录该级别的日志
func (a *KratosLoggerAdapter) shouldLog(level log.Level, configLevel string) bool {
	// 定义日志级别优先级
	levelPriority := map[log.Level]int{
		log.LevelDebug: 0,
		log.LevelInfo:  1,
		log.LevelWarn:  2,
		log.LevelError: 3,
		log.LevelFatal: 4,
	}

	configPriority := map[string]int{
		"debug": 0,
		"info":  1,
		"warn":  2,
		"error": 3,
		"fatal": 4,
	}

	levelPrio, levelExists := levelPriority[level]
	configPrio, configExists := configPriority[configLevel]

	if !levelExists || !configExists {
		return true // 如果级别未知，默认记录
	}

	return levelPrio >= configPrio
}

// writeToFile 将日志写入文件
func (a *KratosLoggerAdapter) writeToFile(level log.Level, message string, keyvals []any) {
	if a.writer == nil {
		return
	}

	// 构建日志行
	logLine := a.formatLogLine(level, message, keyvals)

	// 写入文件
	a.writer.WriteString(logLine + "\n")
	a.writer.Sync() // 确保立即写入磁盘
}

// formatLogLine 格式化日志行
func (a *KratosLoggerAdapter) formatLogLine(level log.Level, message string, keyvals []any) string {
	// 简单的文本格式
	levelStr := strings.ToUpper(level.String())

	// 构建键值对字符串
	var kvPairs []string
	for i := 0; i < len(keyvals); i += 2 {
		if i+1 < len(keyvals) {
			key, ok := keyvals[i].(string)
			if ok {
				value := keyvals[i+1]
				kvPairs = append(kvPairs, fmt.Sprintf("%s=%v", key, value))
			}
		}
	}

	kvStr := ""
	if len(kvPairs) > 0 {
		kvStr = " " + strings.Join(kvPairs, " ")
	}

	return fmt.Sprintf("[%s] %s%s", levelStr, message, kvStr)
}

// convertLevel 转换Kratos日志级别到OpenTelemetry级别
func convertLevel(level log.Level) otellog.Severity {
	switch level {
	case log.LevelDebug:
		return otellog.SeverityDebug
	case log.LevelInfo:
		return otellog.SeverityInfo
	case log.LevelWarn:
		return otellog.SeverityWarn
	case log.LevelError:
		return otellog.SeverityError
	case log.LevelFatal:
		return otellog.SeverityFatal
	default:
		return otellog.SeverityInfo
	}
}

// convertToKeyValue 转换值到OpenTelemetry KeyValue
func convertToKeyValue(key string, value any) otellog.KeyValue {
	switch v := value.(type) {
	case string:
		return otellog.String(key, v)
	case int:
		return otellog.Int(key, v)
	case int64:
		return otellog.Int64(key, v)
	case float64:
		return otellog.Float64(key, v)
	case bool:
		return otellog.Bool(key, v)
	default:
		return otellog.String(key, toString(v))
	}
}

// toString 将任意值转换为字符串
func toString(v any) string {
	if v == nil {
		return "<nil>"
	}
	return string(toStringBytes(v))
}

// toStringBytes 将任意值转换为字节数组
func toStringBytes(v any) []byte {
	switch val := v.(type) {
	case string:
		return []byte(val)
	case []byte:
		return val
	case error:
		return []byte(val.Error())
	default:
		// 对于其他类型，使用默认的字符串表示
		return []byte("unknown type")
	}
}
