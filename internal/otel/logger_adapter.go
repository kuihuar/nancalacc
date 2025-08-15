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
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// KratosLoggerAdapter 适配OpenTelemetry Logger到Kratos Logger
type KratosLoggerAdapter struct {
	otellogger otellog.Logger
	config     *Config
	writer     *lumberjack.Logger
	zapLogger  *zap.Logger // 新增zap logger
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

	// 初始化zap logger
	adapter.initZapLogger()

	return adapter
}

// initZapLogger 初始化zap logger
func (a *KratosLoggerAdapter) initZapLogger() {
	if a.config == nil || !a.config.Logs.Enabled {
		return
	}

	// 检查是否启用zap
	if !a.config.Logs.UseZap {
		return
	}

	var core zapcore.Core
	var err error

	// 根据配置创建zap logger
	switch a.config.Logs.Output {
	case "stdout":
		core, err = a.createZapCore("stdout", nil)
	case "file":
		if a.writer != nil {
			core, err = a.createZapCore("file", a.writer)
		}
	case "both":
		// 同时输出到控制台和文件
		consoleCore, _ := a.createZapCore("stdout", nil)
		fileCore, _ := a.createZapCore("file", a.writer)
		if consoleCore != nil && fileCore != nil {
			core = zapcore.NewTee(consoleCore, fileCore)
		} else if consoleCore != nil {
			core = consoleCore
		} else if fileCore != nil {
			core = fileCore
		}
	}

	if err != nil {
		// 如果创建失败，使用默认配置
		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			zapcore.AddSync(os.Stdout),
			zap.NewAtomicLevelAt(zap.InfoLevel),
		)
	}

	if core != nil {
		// 构建zap选项
		options := []zap.Option{}

		// 添加调用者信息
		if !a.config.Logs.ZapDisableCaller {
			options = append(options, zap.AddCaller(), zap.AddCallerSkip(6))
		}

		// 添加堆栈跟踪
		if !a.config.Logs.ZapDisableStacktrace {
			options = append(options, zap.AddStacktrace(zap.ErrorLevel))
		}

		a.zapLogger = zap.New(core, options...)
	}
}

// createZapCore 创建zap core
func (a *KratosLoggerAdapter) createZapCore(output string, writer *lumberjack.Logger) (zapcore.Core, error) {
	// 配置编码器
	var encoderConfig zapcore.EncoderConfig
	if a.config.Logs.ZapDevelopment {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
	} else {
		encoderConfig = zap.NewProductionEncoderConfig()
	}

	// 自定义字段名
	if a.config.Logs.ZapTimeKey != "" {
		encoderConfig.TimeKey = a.config.Logs.ZapTimeKey
	}
	if a.config.Logs.ZapLevelKey != "" {
		encoderConfig.LevelKey = a.config.Logs.ZapLevelKey
	}
	if a.config.Logs.ZapNameKey != "" {
		encoderConfig.NameKey = a.config.Logs.ZapNameKey
	}
	if a.config.Logs.ZapCallerKey != "" {
		encoderConfig.CallerKey = a.config.Logs.ZapCallerKey
	}
	if a.config.Logs.ZapFunctionKey != "" {
		encoderConfig.FunctionKey = a.config.Logs.ZapFunctionKey
	}
	if a.config.Logs.ZapMessageKey != "" {
		encoderConfig.MessageKey = a.config.Logs.ZapMessageKey
	}
	if a.config.Logs.ZapStacktraceKey != "" {
		encoderConfig.StacktraceKey = a.config.Logs.ZapStacktraceKey
	}

	// 设置时间编码
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	var encoder zapcore.Encoder
	encoding := a.config.Logs.ZapEncoding
	if encoding == "" {
		encoding = a.config.Logs.Format
	}

	switch encoding {
	case "json":
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	case "console":
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	default:
		// 默认使用JSON格式
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	// 配置输出
	var writeSyncer zapcore.WriteSyncer
	switch output {
	case "stdout":
		writeSyncer = zapcore.AddSync(os.Stdout)
	case "file":
		if writer != nil {
			writeSyncer = zapcore.AddSync(writer)
		} else {
			return nil, fmt.Errorf("file writer not available")
		}
	case "both":
		if writer != nil {
			// 同时输出到控制台和文件
			writeSyncer = zapcore.NewMultiWriteSyncer(
				zapcore.AddSync(os.Stdout),
				zapcore.AddSync(writer),
			)
		} else {
			// 如果文件writer不可用，只输出到控制台
			writeSyncer = zapcore.AddSync(os.Stdout)
		}
	default:
		// 默认输出到控制台
		writeSyncer = zapcore.AddSync(os.Stdout)
	}

	// 配置日志级别
	level := a.convertToZapLevel(a.config.Logs.Level)

	return zapcore.NewCore(encoder, writeSyncer, level), nil
}

// convertToZapLevel 转换日志级别到zap级别
func (a *KratosLoggerAdapter) convertToZapLevel(level string) zap.AtomicLevel {
	switch strings.ToLower(level) {
	case "debug":
		return zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		return zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn", "warning":
		return zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		return zap.NewAtomicLevelAt(zap.ErrorLevel)
	case "fatal":
		return zap.NewAtomicLevelAt(zap.FatalLevel)
	default:
		return zap.NewAtomicLevelAt(zap.InfoLevel)
	}
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
	zapFields := make([]zap.Field, 0, len(keyvals)/2)

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

			// 为zap准备字段
			zapField := convertToZapField(key, value)
			zapFields = append(zapFields, zapField)
		}
	}

	// 如果没有消息，使用默认消息
	if message == "" {
		message = "log message"
	}

	// 优先使用zap logger（如果可用）
	if a.zapLogger != nil {
		a.logWithZap(level, message, zapFields)
	} else {
		// 回退到原有实现
		a.logWithOtel(severity, message, attrs)
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

// convertToZapField 转换值到zap Field
func convertToZapField(key string, value any) zap.Field {
	switch v := value.(type) {
	case string:
		return zap.String(key, v)
	case int:
		return zap.Int(key, v)
	case int64:
		return zap.Int64(key, v)
	case float64:
		return zap.Float64(key, v)
	case bool:
		return zap.Bool(key, v)
	case error:
		return zap.Error(v)
	default:
		return zap.Any(key, v)
	}
}

// logWithZap 使用zap记录日志
func (a *KratosLoggerAdapter) logWithZap(level log.Level, message string, fields []zap.Field) {
	switch level {
	case log.LevelDebug:
		a.zapLogger.Debug(message, fields...)
	case log.LevelInfo:
		a.zapLogger.Info(message, fields...)
	case log.LevelWarn:
		a.zapLogger.Warn(message, fields...)
	case log.LevelError:
		a.zapLogger.Error(message, fields...)
	case log.LevelFatal:
		a.zapLogger.Fatal(message, fields...)
	default:
		a.zapLogger.Info(message, fields...)
	}
}

// logWithOtel 使用OpenTelemetry记录日志
func (a *KratosLoggerAdapter) logWithOtel(severity otellog.Severity, message string, attrs []otellog.KeyValue) {
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
		a.writeToStdout(convertOtelLevelToKratos(severity), message, nil)
	}

	// 如果配置了文件输出，同时写入文件
	if a.writer != nil {
		a.writeToFile(convertOtelLevelToKratos(severity), message, nil)
	}
}

// convertOtelLevelToKratos 转换OpenTelemetry级别到Kratos级别
func convertOtelLevelToKratos(severity otellog.Severity) log.Level {
	switch severity {
	case otellog.SeverityDebug:
		return log.LevelDebug
	case otellog.SeverityInfo:
		return log.LevelInfo
	case otellog.SeverityWarn:
		return log.LevelWarn
	case otellog.SeverityError:
		return log.LevelError
	case otellog.SeverityFatal:
		return log.LevelFatal
	default:
		return log.LevelInfo
	}
}
