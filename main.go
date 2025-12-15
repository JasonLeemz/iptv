package main

import (
	"iptv/pkg/config"
	"iptv/pkg/cron"
	"iptv/pkg/log"
	"os"
)

func main() {
	// 初始化日志
	err := log.Init()
	if err != nil {
		os.Exit(1)
	}
	defer log.Close()

	// 加载配置文件
	cfg, err := config.LoadConfig("config/app.yml")
	if err != nil {
		log.Error("加载配置失败: %v", err)
		os.Exit(1)
	}

	// 执行主任务
	runTask(cfg)

	// 如果启用定时任务，启动调度器
	if cfg.Crontab.Enable {
		log.Info("定时任务已启用: %s", cfg.Crontab.Job)
		cron.Init()
		err := cron.AddJob(cfg.Crontab.Job, func() {
			runTask(cfg)
		})
		if err != nil {
			log.Error("添加定时任务失败: %v", err)
			return
		}
		cron.Start()

		// 保持程序运行
		select {}
	}
}
