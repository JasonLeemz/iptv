package cron

import (
	"fmt"
	"time"
)

var (
	schedulerRunning bool
	stopChan         chan bool
	currentJob       func()
	currentCronExpr  string
)

// Init 初始化定时任务调度器
func Init() {
	stopChan = make(chan bool)
	schedulerRunning = false
}

// AddJob 添加定时任务（简化实现，使用标准库）
func AddJob(cronExpr string, job func()) error {
	if currentJob != nil {
		return fmt.Errorf("已存在定时任务，请先清除")
	}

	currentJob = job
	currentCronExpr = cronExpr

	return nil
}

// Start 启动定时任务（简化实现）
func Start() {
	if currentJob == nil {
		return
	}

	if schedulerRunning {
		return
	}

	schedulerRunning = true
	go runScheduler()
}

// runScheduler 运行调度器（简化实现，每天执行一次）
func runScheduler() {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	// 立即执行一次
	if currentJob != nil {
		currentJob()
	}

	for {
		select {
		case <-ticker.C:
			if currentJob != nil {
				currentJob()
			}
		case <-stopChan:
			return
		}
	}
}

// Stop 停止定时任务
func Stop() {
	if schedulerRunning {
		schedulerRunning = false
		if stopChan != nil {
			stopChan <- true
		}
	}
}

// Clear 清除所有任务
func Clear() {
	Stop()
	currentJob = nil
	currentCronExpr = ""
}

