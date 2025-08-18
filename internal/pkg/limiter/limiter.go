// internal/pkg/limiter/limiter.go
package limiter

import (
	"context"
	"fmt"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// RateLimiter 限流器
type RateLimiter struct {
	limiters map[string]*limiterEntry
	mu       sync.RWMutex
	config   *Config
	metrics  *Metrics
}

// limiterEntry 限流器条目
type limiterEntry struct {
	limiter    *rate.Limiter
	lastAccess time.Time
	created    time.Time
	key        string
}

// Config 限流器配置
type Config struct {
	CleanupInterval time.Duration // 清理间隔
	MaxIdleTime     time.Duration // 最大空闲时间
	MaxEntries      int           // 最大条目数
	DefaultRate     rate.Limit    // 默认速率
	DefaultBurst    int           // 默认突发量
}

// Metrics 监控指标
type Metrics struct {
	mu sync.RWMutex
	// 总请求数
	TotalRequests int64
	// 被限流的请求数
	LimitedRequests int64
	// 当前活跃的限流器数量
	ActiveLimiters int64
	// 清理的限流器数量
	CleanedLimiters int64
}

// DefaultConfig 默认配置
func DefaultConfig() *Config {
	return &Config{
		CleanupInterval: 10 * time.Minute,
		MaxIdleTime:     30 * time.Minute,
		MaxEntries:      10000,
		DefaultRate:     10, // 每秒10个请求
		DefaultBurst:    20, // 突发20个请求
	}
}

// NewRateLimiter 创建新的限流器
func NewRateLimiter(config *Config) *RateLimiter {
	if config == nil {
		config = DefaultConfig()
	}

	// 验证配置参数
	if config.CleanupInterval <= 0 {
		config.CleanupInterval = 10 * time.Minute
	}
	if config.MaxIdleTime <= 0 {
		config.MaxIdleTime = 30 * time.Minute
	}
	if config.MaxEntries <= 0 {
		config.MaxEntries = 10000
	}
	if config.DefaultRate <= 0 {
		config.DefaultRate = 10
	}
	if config.DefaultBurst <= 0 {
		config.DefaultBurst = 20
	}

	rl := &RateLimiter{
		limiters: make(map[string]*limiterEntry),
		config:   config,
		metrics:  &Metrics{},
	}

	// 启动清理协程
	go rl.startCleanup()

	return rl
}

// GetLimiter 获取或创建限流器
func (rl *RateLimiter) GetLimiter(key string, r rate.Limit, b int) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if entry, exists := rl.limiters[key]; exists {
		entry.lastAccess = time.Now()
		return entry.limiter
	}

	// 检查是否超过最大条目数
	if len(rl.limiters) >= rl.config.MaxEntries {
		// 清理最旧的条目
		rl.cleanupOldest()
	}

	// 使用默认值如果参数为0
	if r == 0 {
		r = rl.config.DefaultRate
	}
	if b == 0 {
		b = rl.config.DefaultBurst
	}

	limiter := rate.NewLimiter(r, b)
	entry := &limiterEntry{
		limiter:    limiter,
		lastAccess: time.Now(),
		created:    time.Now(),
		key:        key,
	}

	rl.limiters[key] = entry
	rl.metrics.incrementActiveLimiters()

	return limiter
}

// Allow 检查是否允许请求（非阻塞）
func (rl *RateLimiter) Allow(key string, r rate.Limit, b int) bool {
	limiter := rl.GetLimiter(key, r, b)
	allowed := limiter.Allow()

	rl.metrics.incrementTotalRequests()
	if !allowed {
		rl.metrics.incrementLimitedRequests()
	}

	return allowed
}

// Wait 等待直到允许请求（阻塞）
func (rl *RateLimiter) Wait(ctx context.Context, key string, r rate.Limit, b int) error {
	limiter := rl.GetLimiter(key, r, b)
	rl.metrics.incrementTotalRequests()
	return limiter.Wait(ctx)
}

// Reserve 预留一个令牌
func (rl *RateLimiter) Reserve(key string, r rate.Limit, b int) *rate.Reservation {
	limiter := rl.GetLimiter(key, r, b)
	rl.metrics.incrementTotalRequests()
	return limiter.Reserve()
}

// cleanupOldest 清理最旧的条目
func (rl *RateLimiter) cleanupOldest() {
	var oldestKey string
	var oldestTime time.Time

	for key, entry := range rl.limiters {
		if oldestKey == "" || entry.lastAccess.Before(oldestTime) {
			oldestKey = key
			oldestTime = entry.lastAccess
		}
	}

	if oldestKey != "" {
		delete(rl.limiters, oldestKey)
		rl.metrics.decrementActiveLimiters()
		rl.metrics.incrementCleanedLimiters()
	}
}

// startCleanup 启动清理协程
func (rl *RateLimiter) startCleanup() {
	ticker := time.NewTicker(rl.config.CleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		rl.cleanup()
	}
}

// cleanup 清理过期的限流器
func (rl *RateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	cleaned := 0

	for key, entry := range rl.limiters {
		if now.Sub(entry.lastAccess) > rl.config.MaxIdleTime {
			delete(rl.limiters, key)
			cleaned++
		}
	}

	if cleaned > 0 {
		rl.metrics.decrementActiveLimitersBy(int64(cleaned))
		rl.metrics.incrementCleanedLimitersBy(int64(cleaned))
	}
}

// GetStats 获取统计信息
func (rl *RateLimiter) GetStats() map[string]interface{} {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	rl.metrics.mu.RLock()
	defer rl.metrics.mu.RUnlock()

	return map[string]interface{}{
		"active_limiters":  rl.metrics.ActiveLimiters,
		"total_requests":   rl.metrics.TotalRequests,
		"limited_requests": rl.metrics.LimitedRequests,
		"cleaned_limiters": rl.metrics.CleanedLimiters,
		"current_entries":  len(rl.limiters),
		"max_entries":      rl.config.MaxEntries,
		"cleanup_interval": rl.config.CleanupInterval.String(),
		"max_idle_time":    rl.config.MaxIdleTime.String(),
		"default_rate":     float64(rl.config.DefaultRate),
		"default_burst":    rl.config.DefaultBurst,
	}
}

// Reset 重置限流器
func (rl *RateLimiter) Reset() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.limiters = make(map[string]*limiterEntry)
	rl.metrics.reset()
}

// Remove 移除指定的限流器
func (rl *RateLimiter) Remove(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if _, exists := rl.limiters[key]; exists {
		delete(rl.limiters, key)
		rl.metrics.decrementActiveLimiters()
		return true
	}
	return false
}

// ListKeys 列出所有限流器键
func (rl *RateLimiter) ListKeys() []string {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	keys := make([]string, 0, len(rl.limiters))
	for key := range rl.limiters {
		keys = append(keys, key)
	}
	return keys
}

// GetEntryInfo 获取条目信息
func (rl *RateLimiter) GetEntryInfo(key string) (map[string]interface{}, bool) {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	if entry, exists := rl.limiters[key]; exists {
		return map[string]interface{}{
			"key":         entry.key,
			"created":     entry.created,
			"last_access": entry.lastAccess,
			"rate":        float64(entry.limiter.Limit()),
			"burst":       entry.limiter.Burst(),
		}, true
	}
	return nil, false
}

// Metrics 方法实现
func (m *Metrics) incrementTotalRequests() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.TotalRequests++
}

func (m *Metrics) incrementLimitedRequests() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.LimitedRequests++
}

func (m *Metrics) incrementActiveLimiters() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ActiveLimiters++
}

func (m *Metrics) decrementActiveLimiters() {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.ActiveLimiters > 0 {
		m.ActiveLimiters--
	}
}

func (m *Metrics) decrementActiveLimitersBy(count int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.ActiveLimiters >= count {
		m.ActiveLimiters -= count
	} else {
		m.ActiveLimiters = 0
	}
}

func (m *Metrics) incrementCleanedLimiters() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.CleanedLimiters++
}

func (m *Metrics) incrementCleanedLimitersBy(count int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.CleanedLimiters += count
}

func (m *Metrics) reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.TotalRequests = 0
	m.LimitedRequests = 0
	m.ActiveLimiters = 0
	m.CleanedLimiters = 0
}

// String 实现 Stringer 接口
func (rl *RateLimiter) String() string {
	stats := rl.GetStats()
	return fmt.Sprintf("RateLimiter{active_limiters: %v, total_requests: %v, limited_requests: %v}",
		stats["active_limiters"], stats["total_requests"], stats["limited_requests"])
}
