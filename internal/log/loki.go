package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

// LokiClient Loki客户端
type LokiClient struct {
	url        string
	username   string
	password   string
	tenantID   string
	labels     map[string]string
	batchSize  int
	batchWait  time.Duration
	timeout    time.Duration
	client     *http.Client
	batch      []LokiEntry
	batchMutex sync.Mutex
	ticker     *time.Ticker
	done       chan bool
}

// LokiEntry Loki日志条目
type LokiEntry struct {
	Timestamp time.Time         `json:"ts"`
	Line      string            `json:"line"`
	Labels    map[string]string `json:"labels"`
}

// LokiStream Loki流
type LokiStream struct {
	Stream map[string]string `json:"stream"`
	Values [][]string        `json:"values"`
}

// LokiPayload Loki负载
type LokiPayload struct {
	Streams []LokiStream `json:"streams"`
}

// NewLokiClient 创建Loki客户端
func NewLokiClient(config *LokiConfig) *LokiClient {
	client := &LokiClient{
		url:       config.URL,
		username:  config.Username,
		password:  config.Password,
		tenantID:  config.TenantID,
		labels:    config.Labels,
		batchSize: config.BatchSize,
		batchWait: config.BatchWait,
		timeout:   config.Timeout,
		client: &http.Client{
			Timeout: config.Timeout,
		},
		batch: make([]LokiEntry, 0, config.BatchSize),
		done:  make(chan bool),
	}

	// 启动批处理定时器
	client.ticker = time.NewTicker(config.BatchWait)
	go client.batchProcessor()

	return client
}

// Send 发送日志条目
func (l *LokiClient) Send(entry LokiEntry) error {
	l.batchMutex.Lock()
	defer l.batchMutex.Unlock()

	l.batch = append(l.batch, entry)

	// 如果达到批处理大小，立即发送
	if len(l.batch) >= l.batchSize {
		return l.flush()
	}

	return nil
}

// flush 刷新批处理
func (l *LokiClient) flush() error {
	if len(l.batch) == 0 {
		return nil
	}

	// 按标签分组
	streams := make(map[string]*LokiStream)

	for _, entry := range l.batch {
		// 合并标签
		labels := make(map[string]string)
		for k, v := range l.labels {
			labels[k] = v
		}
		for k, v := range entry.Labels {
			labels[k] = v
		}

		// 创建标签字符串
		labelStr := l.formatLabels(labels)
		stream, exists := streams[labelStr]
		if !exists {
			stream = &LokiStream{
				Stream: labels,
				Values: make([][]string, 0),
			}
			streams[labelStr] = stream
		}

		// 添加日志条目
		timestamp := entry.Timestamp.UnixNano()
		stream.Values = append(stream.Values, []string{
			fmt.Sprintf("%d", timestamp),
			entry.Line,
		})
	}

	// 构建负载
	payload := LokiPayload{
		Streams: make([]LokiStream, 0, len(streams)),
	}
	for _, stream := range streams {
		payload.Streams = append(payload.Streams, *stream)
	}

	// 发送到Loki
	return l.sendToLoki(payload)
}

// sendToLoki 发送到Loki
func (l *LokiClient) sendToLoki(payload LokiPayload) error {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload failed: %w", err)
	}

	req, err := http.NewRequest("POST", l.url+"/loki/api/v1/push", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("create request failed: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if l.tenantID != "" {
		req.Header.Set("X-Scope-OrgID", l.tenantID)
	}
	if l.username != "" && l.password != "" {
		req.SetBasicAuth(l.username, l.password)
	}

	resp, err := l.client.Do(req)
	if err != nil {
		return fmt.Errorf("send request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("loki returned status %d", resp.StatusCode)
	}

	// 清空批处理
	l.batch = l.batch[:0]
	return nil
}

// formatLabels 格式化标签
func (l *LokiClient) formatLabels(labels map[string]string) string {
	// 简单的标签格式化，实际使用时可能需要更复杂的逻辑
	result := ""
	for k, v := range labels {
		if result != "" {
			result += ","
		}
		result += fmt.Sprintf("%s=%s", k, v)
	}
	return result
}

// batchProcessor 批处理处理器
func (l *LokiClient) batchProcessor() {
	for {
		select {
		case <-l.ticker.C:
			l.batchMutex.Lock()
			l.flush()
			l.batchMutex.Unlock()
		case <-l.done:
			return
		}
	}
}

// Close 关闭客户端
func (l *LokiClient) Close() error {
	l.ticker.Stop()
	l.done <- true
	close(l.done)

	// 最后一次刷新
	l.batchMutex.Lock()
	defer l.batchMutex.Unlock()
	return l.flush()
}

// LokiWriter Loki日志写入器
type LokiWriter struct {
	client *LokiClient
}

// NewLokiWriter 创建Loki写入器
func NewLokiWriter(client *LokiClient) *LokiWriter {
	return &LokiWriter{client: client}
}

// Write 实现io.Writer接口
func (l *LokiWriter) Write(p []byte) (n int, err error) {
	entry := LokiEntry{
		Timestamp: time.Now(),
		Line:      string(p),
		Labels:    make(map[string]string),
	}

	err = l.client.Send(entry)
	if err != nil {
		return 0, err
	}

	return len(p), nil
}

// LokiLogger Loki日志记录器
type LokiLogger struct {
	client *LokiClient
	level  log.Level
}

// NewLokiLogger 创建Loki日志记录器
func NewLokiLogger(client *LokiClient, level log.Level) *LokiLogger {
	return &LokiLogger{
		client: client,
		level:  level,
	}
}

// Log 实现log.Logger接口
func (l *LokiLogger) Log(level log.Level, keyvals ...interface{}) error {
	if level < l.level {
		return nil
	}

	// 构建日志行
	var buf bytes.Buffer
	for i, kv := range keyvals {
		if i > 0 {
			buf.WriteString(" ")
		}
		buf.WriteString(fmt.Sprintf("%v", kv))
	}

	// 提取标签
	labels := make(map[string]string)
	for i := 0; i < len(keyvals); i += 2 {
		if i+1 < len(keyvals) {
			key, ok := keyvals[i].(string)
			if ok {
				labels[key] = fmt.Sprintf("%v", keyvals[i+1])
			}
		}
	}

	entry := LokiEntry{
		Timestamp: time.Now(),
		Line:      buf.String(),
		Labels:    labels,
	}

	return l.client.Send(entry)
}

// With 添加字段
func (l *LokiLogger) With(keyvals ...interface{}) log.Logger {
	// 对于Loki，我们直接在Log方法中处理字段
	return l
}
