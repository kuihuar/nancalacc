package otel

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

// LokiExporter Loki日志导出器
type LokiExporter struct {
	endpoint string
	client   *http.Client
	logger   *log.Helper
}

// LokiPushRequest Loki推送请求格式
type LokiPushRequest struct {
	Streams []LokiStream `json:"streams"`
}

type LokiStream struct {
	Stream map[string]string `json:"stream"`
	Values [][]string        `json:"values"`
}

// NewLokiExporter 创建Loki导出器
func NewLokiExporter(endpoint string, logger log.Logger) *LokiExporter {
	return &LokiExporter{
		endpoint: endpoint,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		logger: log.NewHelper(log.With(logger, "module", "loki_exporter")),
	}
}

// PushLog 推送日志到Loki
func (e *LokiExporter) PushLog(ctx context.Context, level, message string, labels map[string]string) error {
	// 构建Loki推送请求
	now := time.Now().UnixNano()

	// 确保基础标签存在
	if labels == nil {
		labels = make(map[string]string)
	}
	labels["service"] = "nancalacc"
	labels["level"] = level

	// 构建日志内容
	logContent := map[string]interface{}{
		"level":     level,
		"msg":       message,
		"service":   "nancalacc",
		"timestamp": time.Now().Format(time.RFC3339),
	}

	logJSON, err := json.Marshal(logContent)
	if err != nil {
		return fmt.Errorf("序列化日志失败: %w", err)
	}

	request := LokiPushRequest{
		Streams: []LokiStream{
			{
				Stream: labels,
				Values: [][]string{
					{
						fmt.Sprintf("%d", now),
						string(logJSON),
					},
				},
			},
		},
	}

	// 序列化请求
	jsonData, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("序列化请求失败: %w", err)
	}

	// 发送HTTP请求
	req, err := http.NewRequestWithContext(ctx, "POST", e.endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := e.client.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %w", err)
	}

	// 检查响应状态
	if resp.StatusCode != 200 && resp.StatusCode != 204 {
		return fmt.Errorf("Loki推送失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	e.logger.WithContext(ctx).Debugf("Loki推送成功: %s", string(body))
	return nil
}

// BatchPushLog 批量推送日志
func (e *LokiExporter) BatchPushLog(ctx context.Context, logs []LogEntry) error {
	if len(logs) == 0 {
		return nil
	}

	// 按标签分组日志
	streams := make(map[string]*LokiStream)

	for _, logEntry := range logs {
		// 生成流标识符
		streamKey := fmt.Sprintf("%s_%s", logEntry.Level, logEntry.Service)

		stream, exists := streams[streamKey]
		if !exists {
			stream = &LokiStream{
				Stream: map[string]string{
					"service": logEntry.Service,
					"level":   logEntry.Level,
				},
				Values: make([][]string, 0),
			}
			streams[streamKey] = stream
		}

		// 添加自定义标签
		for k, v := range logEntry.Labels {
			stream.Stream[k] = v
		}

		// 构建日志内容
		logContent := map[string]interface{}{
			"level":     logEntry.Level,
			"msg":       logEntry.Message,
			"service":   logEntry.Service,
			"timestamp": logEntry.Timestamp.Format(time.RFC3339),
		}

		logJSON, err := json.Marshal(logContent)
		if err != nil {
			return fmt.Errorf("序列化日志失败: %w", err)
		}

		stream.Values = append(stream.Values, []string{
			fmt.Sprintf("%d", logEntry.Timestamp.UnixNano()),
			string(logJSON),
		})
	}

	// 构建请求
	request := LokiPushRequest{
		Streams: make([]LokiStream, 0, len(streams)),
	}

	for _, stream := range streams {
		request.Streams = append(request.Streams, *stream)
	}

	// 序列化请求
	jsonData, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("序列化请求失败: %w", err)
	}

	// 发送HTTP请求
	req, err := http.NewRequestWithContext(ctx, "POST", e.endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := e.client.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %w", err)
	}

	// 检查响应状态
	if resp.StatusCode != 200 && resp.StatusCode != 204 {
		return fmt.Errorf("Loki批量推送失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	e.logger.WithContext(ctx).Debugf("Loki批量推送成功，推送了 %d 条日志", len(logs))
	return nil
}

// LogEntry 日志条目
type LogEntry struct {
	Level     string
	Message   string
	Service   string
	Timestamp time.Time
	Labels    map[string]string
}
