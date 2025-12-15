package main

import (
	"bufio"
	"io"
	"iptv/dto"
	"iptv/pkg/bark"
	"iptv/pkg/config"
	"iptv/pkg/log"
	"os"
	"path/filepath"
	"strings"
)

// runTask 执行主任务
func runTask(cfg *config.Config) {
	log.Info("============================================================")
	log.Info("开始执行IPTV频道汇总任务")
	_ = bark.Push("IPTV", "开始执行IPTV频道汇总任务")
	log.Info("============================================================")

	// 1. 更新组播源列表（从iptvmulticast.php获取，数量从配置读取）
	log.Info("[步骤1] 更新组播源列表...")
	sources, err := FetchMulticastIPs(cfg.Cookie.Data)
	if err != nil {
		log.Warn("获取组播源失败: %v", err)
		_ = bark.Push("IPTV", "获取组播源失败: %v", err.Error())
		log.Info("将使用config/source.txt中的现有URL")
	} else {
		log.Info("成功获取 %d 个组播源IP", len(sources))
		_ = bark.Push("IPTV", "成功获取 %d 个组播源IP", len(sources))
		err = UpdateSourceFile(sources, "config")
		if err != nil {
			log.Warn("更新source.txt失败: %v", err)
		} else {
			log.Info("已更新config/source.txt")
		}
	}

	// 2. 读取URL列表
	log.Info("[步骤2] 读取URL列表...")
	urls, err := readURLsFromFile("config/source.txt")
	if err != nil {
		log.Error("读取【config/source.txt】配置文件失败: %v", err)
		_ = bark.Push("IPTV", "读取【config/source.txt】配置文件失败: %v", err.Error())
		os.Exit(1)
	}

	if len(urls) == 0 {
		log.Error("配置文件中没有找到URL")
		_ = bark.Push("IPTV", "配置文件中没有找到URL")
		os.Exit(1)
	}

	log.Info("从配置文件读取到 %d 个URL", len(urls))

	// 3. 获取所有频道并汇总
	log.Info("[步骤3] 获取频道数据...")
	var allChannels []dto.Channel
	channelMap := make(map[string]bool) // 用于去重

	for i, pageURL := range urls {
		log.Info("[%d/%d] 正在处理: %s", i+1, len(urls), pageURL)

		channels, err := FetchChannelsFromURL(pageURL, cfg.Cookie.Data)
		if err != nil {
			log.Warn("获取频道数据失败: %s,URL:%s", err.Error(), pageURL)
			_ = bark.Push("IPTV", "获取频道数据失败: %s,URL:%s", err.Error(), pageURL)
			continue
		}

		// 去重并添加到汇总列表
		for _, ch := range channels {
			if !channelMap[ch.URL] {
				channelMap[ch.URL] = true
				allChannels = append(allChannels, ch)
			}
		}

		log.Info("成功获取 %d 个频道（累计: %d 个唯一频道）", len(channels), len(allChannels))
		_ = bark.Push("IPTV", "成功获取 %d 个频道（累计: %d 个唯一频道）", len(channels), len(allChannels))
	}

	if len(allChannels) == 0 {
		log.Error("未找到任何频道数据，请检查cookies是否有效")
		_ = bark.Push("IPTV", "未找到任何频道数据，请检查cookies是否有效")
		os.Exit(1)
	}

	// 4. 输出结果
	log.Info("[步骤4] 输出结果...")

	// 输出M3U格式
	m3uPath := cfg.Output.M3U
	err = AggregateChannelsToM3U(allChannels, m3uPath)
	if err != nil {
		log.Error("输出M3U文件失败: %v", err)
		_ = bark.Push("IPTV", "输出M3U文件失败: %v", err.Error())
	} else {
		log.Info("M3U格式: %s", m3uPath)
	}

	// 输出TXT格式
	txtPath := cfg.Output.Local
	err = AggregateChannelsToTXT(allChannels, txtPath)
	if err != nil {
		log.Error("输出TXT文件失败: %v", err)
		_ = bark.Push("IPTV", "输出TXT文件失败: %v", err.Error())
	} else {
		log.Info("CSV格式: %s", txtPath)
	}

	log.Info("成功汇总 %d 个唯一频道", len(allChannels))
	_ = bark.Push("IPTV", "成功汇总 %d 个唯一频道", len(allChannels))

	// 5. 重定向输出（如果启用）
	if cfg.RedirectOutput.Enable {
		log.Info("[步骤5] 重定向输出文件...")
		err = redirectOutput(cfg)
		if err != nil {
			log.Warn("重定向输出文件失败: %v", err)
			_ = bark.Push("IPTV", "重定向输出文件失败: %v", err.Error())
		} else {
			log.Info("成功重定向输出文件: %s -> %s", cfg.RedirectOutput.Move, cfg.RedirectOutput.To)
			_ = bark.Push("IPTV", "重定向输出文件失败: %v", "成功重定向输出文件: %s -> %s", cfg.RedirectOutput.Move, cfg.RedirectOutput.To)
		}
	}

	log.Info("============================================================")
}

// readURLsFromFile 从文件中读取URL列表
func readURLsFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var urls []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// 跳过空行和注释行
		if line != "" && !strings.HasPrefix(line, "#") {
			urls = append(urls, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return urls, nil
}

// redirectOutput 拷贝解析结果到指定位置
func redirectOutput(cfg *config.Config) error {
	if !cfg.RedirectOutput.Enable {
		return nil
	}

	sourceFile := cfg.RedirectOutput.Move
	targetFile := cfg.RedirectOutput.To

	if sourceFile == "" || targetFile == "" {
		return nil
	}

	// 检查源文件是否存在
	if _, err := os.Stat(sourceFile); os.IsNotExist(err) {
		return nil // 源文件不存在，不报错，静默返回
	}

	// 确保目标目录存在
	targetDir := filepath.Dir(targetFile)
	if targetDir != "." && targetDir != "" {
		err := os.MkdirAll(targetDir, 0755)
		if err != nil {
			return err
		}
	}

	// 打开源文件
	source, err := os.Open(sourceFile)
	if err != nil {
		return err
	}
	defer source.Close()

	// 创建目标文件（覆盖模式）
	target, err := os.Create(targetFile)
	if err != nil {
		return err
	}
	defer target.Close()

	// 拷贝文件内容
	_, err = io.Copy(target, source)
	if err != nil {
		return err
	}

	return nil
}
