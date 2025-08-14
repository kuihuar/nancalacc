package log

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	kratoslog "github.com/go-kratos/kratos/v2/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger 日志接口
type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
	WithFields(fields map[string]interface{}) Logger
	Log(level kratoslog.Level, keyvals ...interface{}) error
}

// Field 日志字段
type Field struct {
	Key   string
	Value interface{}
}

// NewField 创建字段
func NewField(key string, value interface{}) Field {
	return Field{Key: key, Value: value}
}

// zapLogger zap日志实现
type zapLogger struct {
	logger *zap.Logger
	level  zap.AtomicLevel
}

// NewLogger 创建日志记录器
func NewLogger(config *Config) (Logger, error) {
	// 解析日志级别
	level, err := parseLevel(config.Level)
	if err != nil {
		return nil, fmt.Errorf("parse level failed: %w", err)
	}

	// 创建编码器配置
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	// 选择编码器
	var encoder zapcore.Encoder
	if strings.ToLower(config.Format) == "json" {
		// 使用自定义的JSON编码器配置，避免换行符转义
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// 创建输出
	var outputs []zapcore.WriteSyncer
	switch strings.ToLower(config.Output) {
	case "stdout":
		outputs = append(outputs, zapcore.AddSync(os.Stdout))
	case "file":
		outputs = append(outputs, zapcore.AddSync(newFileWriter(config)))
	case "both":
		outputs = append(outputs, zapcore.AddSync(os.Stdout))
		outputs = append(outputs, zapcore.AddSync(newFileWriter(config)))
	default:
		outputs = append(outputs, zapcore.AddSync(os.Stdout))
	}

	// 如果启用Loki，添加Loki输出
	if config.Loki.Enable {
		lokiClient := NewLokiClient(&config.Loki)
		lokiWriter := NewLokiWriter(lokiClient)
		outputs = append(outputs, zapcore.AddSync(lokiWriter))
	}

	// 创建核心
	core := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(outputs...), level)

	// 创建选项
	opts := []zap.Option{
		zap.AddCaller(),
		zap.AddCallerSkip(4), // 跳过Helper包装层，直接显示业务代码调用位置
	}

	if config.Stacktrace {
		opts = append(opts, zap.AddStacktrace(zapcore.ErrorLevel))
	}

	// 创建logger
	zapLoggerInstance := zap.New(core, opts...)

	return &zapLogger{
		logger: zapLoggerInstance,
		level:  level,
	}, nil
}

// newFileWriter 创建文件写入器
func newFileWriter(config *Config) io.Writer {
	// 确保日志目录存在
	logDir := filepath.Dir(config.FilePath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		panic(fmt.Sprintf("create log directory failed: %v", err))
	}

	return &lumberjack.Logger{
		Filename:   config.FilePath,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
	}
}

// parseLevel 解析日志级别
func parseLevel(level string) (zap.AtomicLevel, error) {
	switch strings.ToLower(level) {
	case "debug":
		return zap.NewAtomicLevelAt(zap.DebugLevel), nil
	case "info":
		return zap.NewAtomicLevelAt(zap.InfoLevel), nil
	case "warn":
		return zap.NewAtomicLevelAt(zap.WarnLevel), nil
	case "error":
		return zap.NewAtomicLevelAt(zap.ErrorLevel), nil
	case "fatal":
		return zap.NewAtomicLevelAt(zap.FatalLevel), nil
	default:
		return zap.NewAtomicLevelAt(zap.InfoLevel), fmt.Errorf("unknown level: %s", level)
	}
}

// Debug 调试日志
func (z *zapLogger) Debug(msg string, fields ...Field) {
	zapFields := convertFields(fields)
	z.logger.Debug(msg, zapFields...)
}

// Info 信息日志
func (z *zapLogger) Info(msg string, fields ...Field) {
	zapFields := convertFields(fields)
	z.logger.Info(msg, zapFields...)
}

// Warn 警告日志
func (z *zapLogger) Warn(msg string, fields ...Field) {
	zapFields := convertFields(fields)
	z.logger.Warn(msg, zapFields...)
}

// Error 错误日志
func (z *zapLogger) Error(msg string, fields ...Field) {
	zapFields := convertFields(fields)
	z.logger.Error(msg, zapFields...)
}

// Fatal 致命错误日志
func (z *zapLogger) Fatal(msg string, fields ...Field) {
	zapFields := convertFields(fields)
	z.logger.Fatal(msg, zapFields...)
}

// WithFields 添加字段
func (z *zapLogger) WithFields(fields map[string]interface{}) Logger {
	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}

	newLogger := z.logger.With(zapFields...)
	return &zapLogger{
		logger: newLogger,
		level:  z.level,
	}
}

// Log 实现log.Logger接口
func (z *zapLogger) Log(level log.Level, keyvals ...interface{}) error {
	// 转换级别
	zapLevel := convertLevel(level)
	if !z.level.Enabled(zapLevel) {
		return nil
	}

	// 构建消息和字段
	var msg string
	var fields []zap.Field

	// 处理 keyvals，每两个一组：key, value
	for i := 0; i < len(keyvals); i += 2 {
		if i+1 < len(keyvals) {
			key, ok := keyvals[i].(string)
			if ok {
				value := keyvals[i+1]
				if key == "msg" || key == "message" {
					msg = fmt.Sprintf("%v", value)
				} else {
					// 处理 log.Valuer 类型
					if valuer, ok := value.(log.Valuer); ok {
						fields = append(fields, zap.Any(key, valuer(context.Background())))
					} else {
						fields = append(fields, zap.Any(key, value))
					}
				}
			}
		}
	}

	if msg == "" {
		msg = "log message"
	}

	// 记录日志
	switch zapLevel {
	case zap.DebugLevel:
		z.logger.Debug(msg, fields...)
	case zap.InfoLevel:
		z.logger.Info(msg, fields...)
	case zap.WarnLevel:
		z.logger.Warn(msg, fields...)
	case zap.ErrorLevel:
		z.logger.Error(msg, fields...)
	case zap.FatalLevel:
		z.logger.Fatal(msg, fields...)
	}

	return nil
}

// convertFields 转换字段
func convertFields(fields []Field) []zap.Field {
	zapFields := make([]zap.Field, len(fields))
	for i, field := range fields {
		zapFields[i] = zap.Any(field.Key, field.Value)
	}
	return zapFields
}

// convertLevel 转换级别
func convertLevel(level log.Level) zapcore.Level {
	switch level {
	case log.LevelDebug:
		return zap.DebugLevel
	case log.LevelInfo:
		return zap.InfoLevel
	case log.LevelWarn:
		return zap.WarnLevel
	case log.LevelError:
		return zap.ErrorLevel
	case log.LevelFatal:
		return zap.FatalLevel
	default:
		return zap.InfoLevel
	}
}

// Helper 日志助手
type Helper struct {
	logger Logger
}

// Debug 调试日志
func (h *Helper) Debug(msg string, fields ...Field) {
	h.logger.Debug(msg, fields...)
}

// Info 信息日志
func (h *Helper) Info(msg string, fields ...Field) {
	h.logger.Info(msg, fields...)
}

// Warn 警告日志
func (h *Helper) Warn(msg string, fields ...Field) {
	h.logger.Warn(msg, fields...)
}

// Error 错误日志
func (h *Helper) Error(msg string, fields ...Field) {
	h.logger.Error(msg, fields...)
}

// Fatal 致命错误日志
func (h *Helper) Fatal(msg string, fields ...Field) {
	h.logger.Fatal(msg, fields...)
}

// Infof 格式化信息日志
func (h *Helper) Infof(format string, args ...interface{}) {
	h.logger.Info(fmt.Sprintf(format, args...))
}

// Warnf 格式化警告日志
func (h *Helper) Warnf(format string, args ...interface{}) {
	h.logger.Warn(fmt.Sprintf(format, args...))
}

// Errorf 格式化错误日志
func (h *Helper) Errorf(format string, args ...interface{}) {
	h.logger.Error(fmt.Sprintf(format, args...))
}

// Fatalf 格式化致命错误日志
func (h *Helper) Fatalf(format string, args ...interface{}) {
	h.logger.Fatal(fmt.Sprintf(format, args...))
}

// WithFields 添加字段
func (h *Helper) WithFields(fields map[string]interface{}) *Helper {
	return &Helper{
		logger: h.logger.WithFields(fields),
	}
}

// WithField 添加单个字段
func (h *Helper) WithField(key string, value interface{}) *Helper {
	return h.WithFields(map[string]interface{}{key: value})
}

// WithError 添加错误字段
func (h *Helper) WithError(err error) *Helper {
	return h.WithField("error", err.Error())
}

// WithContext 添加上下文字段
func (h *Helper) WithContext(ctx interface{}) *Helper {
	return h.WithField("context", ctx)
}

// WithTime 添加时间字段
func (h *Helper) WithTime(t time.Time) *Helper {
	return h.WithField("time", t)
}

// WithDuration 添加持续时间字段
func (h *Helper) WithDuration(d time.Duration) *Helper {
	return h.WithField("duration", d)
}

// WithRequestID 添加请求ID字段
func (h *Helper) WithRequestID(requestID string) *Helper {
	return h.WithField("request_id", requestID)
}

// WithUserID 添加用户ID字段
func (h *Helper) WithUserID(userID string) *Helper {
	return h.WithField("user_id", userID)
}

// WithTraceID 添加追踪ID字段
func (h *Helper) WithTraceID(traceID string) *Helper {
	return h.WithField("trace_id", traceID)
}

// WithSpanID 添加跨度ID字段
func (h *Helper) WithSpanID(spanID string) *Helper {
	return h.WithField("span_id", spanID)
}

// Log 实现log.Logger接口
func (h *Helper) Log(level log.Level, keyvals ...interface{}) error {
	return h.logger.Log(level, keyvals...)
}

// With 实现log.Logger接口
func (h *Helper) With(keyvals ...interface{}) log.Logger {
	if len(keyvals) == 0 {
		return h
	}

	fields := make(map[string]interface{})
	for i := 0; i < len(keyvals); i += 2 {
		if i+1 < len(keyvals) {
			key, ok := keyvals[i].(string)
			if ok {
				fields[key] = keyvals[i+1]
			}
		}
	}

	return &Helper{
		logger: h.logger.WithFields(fields),
	}
}

// // 1. 创建日志记录器
// logger, err := log.NewLogger(config)

// // 2. 基础日志记录
// logger.Info("用户登录", log.NewField("user_id", "123"))

// // 3. 使用日志助手
// helper := log.NewLogHelper(logger)
// helper.Info("业务操作")

// // 4. 在GORM中使用
// db.Logger = log.NewGormLoggerAdapter(logger)

// // 5. 在HTTP中间件中使用
// server.Use(log.HTTPLogMiddleware(logger))

// NewHelper 创建日志助手
func NewHelper(logger kratoslog.Logger) *Helper {
	// 将 kratos Logger 转换为我们的 Logger
	adapter := &kratosLoggerAdapter{logger: logger}
	return &Helper{
		logger: adapter,
	}
}

// kratosLoggerAdapter 适配 kratos Logger 到我们的 Logger 接口
type kratosLoggerAdapter struct {
	logger kratoslog.Logger
	fields map[string]interface{} // 添加字段存储
}

func (a *kratosLoggerAdapter) Debug(msg string, fields ...Field) {
	keyvals := []interface{}{"msg", msg}
	// 添加存储的字段
	if a.fields != nil && len(a.fields) > 0 {
		for k, v := range a.fields {
			keyvals = append(keyvals, k, v)
		}
	}
	// 添加传入的字段
	for _, field := range fields {
		keyvals = append(keyvals, field.Key, field.Value)
	}
	a.logger.Log(kratoslog.LevelDebug, keyvals...)
}

func (a *kratosLoggerAdapter) Info(msg string, fields ...Field) {
	keyvals := []interface{}{"msg", msg}
	// 添加存储的字段
	if a.fields != nil && len(a.fields) > 0 {
		for k, v := range a.fields {
			keyvals = append(keyvals, k, v)
		}
	}
	// 添加传入的字段
	for _, field := range fields {
		keyvals = append(keyvals, field.Key, field.Value)
	}
	a.logger.Log(kratoslog.LevelInfo, keyvals...)
}

func (a *kratosLoggerAdapter) Warn(msg string, fields ...Field) {
	keyvals := []interface{}{"msg", msg}
	// 添加存储的字段
	if a.fields != nil && len(a.fields) > 0 {
		for k, v := range a.fields {
			keyvals = append(keyvals, k, v)
		}
	}
	// 添加传入的字段
	for _, field := range fields {
		keyvals = append(keyvals, field.Key, field.Value)
	}
	a.logger.Log(kratoslog.LevelWarn, keyvals...)
}

func (a *kratosLoggerAdapter) Error(msg string, fields ...Field) {
	keyvals := []interface{}{"msg", msg}
	// 添加存储的字段
	if a.fields != nil && len(a.fields) > 0 {
		for k, v := range a.fields {
			keyvals = append(keyvals, k, v)
		}
	}
	// 添加传入的字段
	for _, field := range fields {
		keyvals = append(keyvals, field.Key, field.Value)
	}
	a.logger.Log(kratoslog.LevelError, keyvals...)
}

func (a *kratosLoggerAdapter) Fatal(msg string, fields ...Field) {
	keyvals := []interface{}{"msg", msg}
	// 添加存储的字段
	if a.fields != nil && len(a.fields) > 0 {
		for k, v := range a.fields {
			keyvals = append(keyvals, k, v)
		}
	}
	// 添加传入的字段
	for _, field := range fields {
		keyvals = append(keyvals, field.Key, field.Value)
	}
	a.logger.Log(kratoslog.LevelFatal, keyvals...)
}

func (a *kratosLoggerAdapter) WithFields(fields map[string]interface{}) Logger {
	// 合并现有字段和新字段
	mergedFields := make(map[string]interface{})
	if a.fields != nil {
		for k, v := range a.fields {
			mergedFields[k] = v
		}
	}
	for k, v := range fields {
		mergedFields[k] = v
	}

	return &kratosLoggerAdapter{
		logger: a.logger,
		fields: mergedFields,
	}
}

func (a *kratosLoggerAdapter) Log(level kratoslog.Level, keyvals ...interface{}) error {
	// 如果有字段，添加到日志中
	if a.fields != nil && len(a.fields) > 0 {
		// 将字段转换为 keyvals 格式
		for k, v := range a.fields {
			keyvals = append(keyvals, k, v)
		}
	}
	return a.logger.Log(level, keyvals...)
}
