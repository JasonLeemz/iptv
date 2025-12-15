package cron

import (
	"fmt"

	"github.com/robfig/cron/v3"
)

var (
	cronInstance *cron.Cron
	entryID      cron.EntryID
)

// Init 初始化定时任务调度器
func Init() {
	// 创建 cron 实例（标准格式：分 时 日 月 星期）
	cronInstance = cron.New()
}

// AddJob 添加定时任务
func AddJob(cronExpr string, job func()) error {
	if cronInstance == nil {
		return fmt.Errorf("调度器未初始化，请先调用 Init()")
	}

	// 如果已存在任务，先移除
	if entryID != 0 {
		cronInstance.Remove(entryID)
	}

	// 添加新任务
	id, err := cronInstance.AddFunc(cronExpr, job)
	if err != nil {
		return fmt.Errorf("添加定时任务失败: %v", err)
	}

	entryID = id
	return nil
}

// Start 启动定时任务
func Start() {
	if cronInstance == nil {
		return
	}

	// 启动 cron 调度器
	cronInstance.Start()
}

// Stop 停止定时任务
func Stop() {
	if cronInstance != nil {
		cronInstance.Stop()
	}
}

// Clear 清除所有任务
func Clear() {
	Stop()
	if cronInstance != nil && entryID != 0 {
		cronInstance.Remove(entryID)
		entryID = 0
	}
}
