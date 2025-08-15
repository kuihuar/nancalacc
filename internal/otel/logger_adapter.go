package otel

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	otellog "go.opentelemetry.io/otel/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

// KratosLoggerAdapter 适配OpenTelemetry Logger到Kratos Logger
type KratosLoggerAdapter struct {
	otellogger otellog.Logger
	config     *Config
	writer     *lumberjack.Logger
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
			// 创建lumberjack日志轮转器
			adapter.writer = &lumberjack.Logger{
				Filename:   config.Logs.FilePath,   // 日志文件路径
				MaxSize:    config.Logs.MaxSize,    // 单个文件最大大小(MB)
				MaxBackups: config.Logs.MaxBackups, // 最大备份文件数
				MaxAge:     config.Logs.MaxAge,     // 最大保留天数
				Compress:   config.Logs.Compress,   // 是否压缩
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

	// 如果配置了stdout输出，输出到控制台
	if a.config != nil && a.config.Logs.Output == "stdout" {
		a.writeToStdout(level, message, keyvals)
	}

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

	// 写入文件（lumberjack会自动处理轮转）
	a.writer.Write([]byte(logLine + "\n"))
}

// formatLogLine 格式化日志行
func (a *KratosLoggerAdapter) formatLogLine(level log.Level, message string, keyvals []any) string {
	// 检查是否配置为JSON格式
	if a.config != nil && a.config.Logs.Format == "json" {
		return a.formatJSON(level, message, keyvals)
	}

	// 默认使用文本格式
	return a.formatText(level, message, keyvals)
}

// formatText 格式化文本日志
func (a *KratosLoggerAdapter) formatText(level log.Level, message string, keyvals []any) string {
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

// formatJSON 格式化JSON日志
func (a *KratosLoggerAdapter) formatJSON(level log.Level, message string, keyvals []any) string {
	// 构建JSON对象
	logEntry := map[string]interface{}{
		"level":   strings.ToUpper(level.String()),
		"message": message,
		"time":    time.Now().Format(time.RFC3339),
	}

	// 添加键值对
	for i := 0; i < len(keyvals); i += 2 {
		if i+1 < len(keyvals) {
			key, ok := keyvals[i].(string)
			if ok {
				value := keyvals[i+1]
				logEntry[key] = value
			}
		}
	}

	// 序列化为JSON
	jsonBytes, err := json.Marshal(logEntry)
	if err != nil {
		// 如果JSON序列化失败，回退到文本格式
		return a.formatText(level, message, keyvals)
	}

	return string(jsonBytes)
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

// Close 关闭日志文件
func (a *KratosLoggerAdapter) Close() error {
	if a.writer != nil {
		return a.writer.Close()
	}
	return nil
}

// writeToStdout 将日志输出到控制台
func (a *KratosLoggerAdapter) writeToStdout(level log.Level, message string, keyvals []any) {
	// 构建日志行
	logLine := a.formatLogLine(level, message, keyvals)

	// 输出到标准输出
	fmt.Println(logLine)
}
