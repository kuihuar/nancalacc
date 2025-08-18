// internal/task/cron.go
package task

import (
	"context"
	"fmt"
	"sync"

	"nancalacc/internal/biz"
	"nancalacc/internal/data"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/robfig/cron/v3"
)

// CronService 定时任务服务
type CronService struct {
	cron            *cron.Cron
	log             log.Logger
	fullSyncUsecase *biz.FullSyncUsecase
	data            *data.Data
	// 添加任务名称映射，防止重复注册
	registeredJobs map[string]cron.EntryID
	mu             sync.RWMutex
}

// Start 启动定时任务服务
func (s *CronService) Start() error {
	s.log.Log(log.LevelInfo, "msg", "starting cron service")
	s.cron.Start()
	return nil
}

// Stop 停止定时任务服务
func (s *CronService) Stop() error {
	s.log.Log(log.LevelInfo, "msg", "stopping cron service")
	ctx := s.cron.Stop()
	<-ctx.Done()
	return nil
}

// AddFuncWithContext 添加带上下文的定时任务
func (s *CronService) AddFuncWithContext(name, spec string, cmd func(context.Context) error) (cron.EntryID, error) {
	// 检查任务是否已注册
	s.mu.RLock()
	if existingID, exists := s.registeredJobs[name]; exists {
		s.mu.RUnlock()
		s.log.Log(log.LevelWarn, "msg", "job already registered", "name", name, "existing_id", existingID)
		return existingID, nil
	}
	s.mu.RUnlock()

	entryID, err := s.cron.AddFunc(spec, func() {
		ctx := context.Background()
		if err := cmd(ctx); err != nil {
			s.log.Log(log.LevelError, "msg", "AddFuncWithContext.AddFunc", "name", name, "err", err)
		} else {
			s.log.Log(log.LevelInfo, "msg", "AddFuncWithContext.AddFunc", "name", name)
		}
	})

	if err != nil {
		return 0, fmt.Errorf("failed to add job [%s]: %w", name, err)
	}

	// 注册成功后，记录任务名称
	s.mu.Lock()
	s.registeredJobs[name] = entryID
	s.mu.Unlock()

	return entryID, nil
}

// GetEntries 获取所有任务条目
func (s *CronService) GetEntries() []cron.Entry {
	return s.cron.Entries()
}

// RemoveEntry 移除任务条目
func (s *CronService) RemoveEntry(id cron.EntryID) {
	s.cron.Remove(id)
}

// NewCronService 创建定时任务服务
func NewCronService(fullSyncUsecase *biz.FullSyncUsecase, data *data.Data, log log.Logger) *CronService {
	return &CronService{
		cron:            cron.New(cron.WithSeconds()),
		log:             log,
		fullSyncUsecase: fullSyncUsecase,
		data:            data,
		registeredJobs:  make(map[string]cron.EntryID),
	}
}
