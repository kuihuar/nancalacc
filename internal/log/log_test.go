package log

import (
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

func TestNewLogger(t *testing.T) {
	config := &Config{
		Level:      "debug",
		Format:     "json",
		Output:     "stdout",
		Caller:     true,
		Stacktrace: false,
	}

	logger, err := NewLogger(config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// 测试不同级别的日志
	logger.Debug("Debug message", NewField("test", "debug"))
	logger.Info("Info message", NewField("test", "info"))
	logger.Warn("Warn message", NewField("test", "warn"))
	logger.Error("Error message", NewField("test", "error"))
}

func TestHelper(t *testing.T) {
	config := &Config{
		Level:  "info",
		Format: "console",
		Output: "stdout",
	}

	logger, err := NewLogger(config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	helper := NewLogHelper(logger)

	// 测试基础方法
	helper.Info("Helper info message")
	helper.Warn("Helper warn message")
	helper.Error("Helper error message")

	// 测试添加上下文
	helperWithCtx := helper.WithRequestID("test-request-id")
	helperWithCtx.Info("Message with request ID")

	// 测试添加字段
	helperWithField := helper.WithField("user_id", "123")
	helperWithField.Info("Message with user ID")
}

func TestLokiClient(t *testing.T) {
	config := &LokiConfig{
		URL:       "http://localhost:3100",
		Enable:    false, // 测试时不启用
		BatchSize: 10,
		BatchWait: time.Second,
		Timeout:   5 * time.Second,
		Labels: map[string]string{
			"service": "test",
		},
	}

	client := NewLokiClient(config)
	if client == nil {
		t.Skip("Loki client not created (disabled)")
	}

	// 测试发送日志条目
	entry := LokiEntry{
		Timestamp: time.Now(),
		Line:      "Test log entry",
		Labels: map[string]string{
			"level": "info",
		},
	}

	err := client.Send(entry)
	if err != nil {
		t.Logf("Failed to send log entry: %v", err)
	}
}

func TestKratosLoggerAdapter(t *testing.T) {
	config := &Config{
		Level:  "info",
		Format: "console",
		Output: "stdout",
	}

	logger, err := NewLogger(config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	kratosLogger := NewKratosLoggerAdapter(logger)

	// 测试Kratos日志接口
	err = kratosLogger.Log(log.LevelInfo, "msg", "Test message", "level", "info")
	if err != nil {
		t.Errorf("Failed to log message: %v", err)
	}

	// 测试添加字段
	// loggerWithFields := kratosLogger.With("service", "test", "version", "1.0.0")
	// err = loggerWithFields.Log(log.LevelInfo, "msg", "Message with fields")
	// if err != nil {
	// 	t.Errorf("Failed to log message with fields: %v", err)
	// }
}

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: &Config{
				Level:  "info",
				Format: "json",
				Output: "stdout",
			},
			wantErr: false,
		},
		{
			name: "invalid level",
			config: &Config{
				Level:  "invalid",
				Format: "json",
				Output: "stdout",
			},
			wantErr: true,
		},
		{
			name: "invalid format",
			config: &Config{
				Level:  "info",
				Format: "invalid",
				Output: "stdout",
			},
			wantErr: false, // 格式无效时使用默认值
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewLogger(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewLogger() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func BenchmarkLogger(b *testing.B) {
	config := &Config{
		Level:  "info",
		Format: "json",
		Output: "stdout",
	}

	logger, err := NewLogger(config)
	if err != nil {
		b.Fatalf("Failed to create logger: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("Benchmark message", NewField("iteration", i))
	}
}

func BenchmarkHelper(b *testing.B) {
	config := &Config{
		Level:  "info",
		Format: "json",
		Output: "stdout",
	}

	logger, err := NewLogger(config)
	if err != nil {
		b.Fatalf("Failed to create logger: %v", err)
	}

	helper := NewLogHelper(logger)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		helper.Info("Benchmark helper message", NewField("iteration", i))
	}
}
