package log

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"iptv/pkg/bark"
	"iptv/pkg/config"
)

var (
	logFile   *os.File
	logDir    = "logs"
	logPrefix = "app"
)

// Init 初始化日志系统
func Init() error {
	// 从配置读取日志路径
	cfg := config.GetConfig()
	if cfg != nil && cfg.Log.Path != "" {
		logDir = cfg.Log.Path
	}

	// 确保logs目录存在
	err := os.MkdirAll(logDir, 0755)
	if err != nil {
		return fmt.Errorf("创建logs目录失败: %v", err)
	}

	// 生成日志文件名（按日期）
	today := time.Now().Format("2006-01-02")
	logFileName := filepath.Join(logDir, fmt.Sprintf("%s-%s.log", logPrefix, today))

	// 打开日志文件（追加模式）
	file, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("打开日志文件失败: %v", err)
	}

	logFile = file
	return nil
}

// Close 关闭日志文件
func Close() {
	if logFile != nil {
		logFile.Close()
	}
}

// writeLog 写入日志（内部函数）
func writeLog(level string, format string, args ...interface{}) {
	if logFile == nil {
		// 如果未初始化，尝试初始化
		if err := Init(); err != nil {
			return
		}
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf(format, args...)
	logLine := fmt.Sprintf("[%s] [%s] %s\n", timestamp, level, message)

	logFile.WriteString(logLine)
	logFile.Sync() // 立即刷新到磁盘
}

// Info 记录信息日志
func Info(format string, args ...interface{}) {
	writeLog("INFO", format, args...)
}

// Error 记录错误日志
func Error(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	writeLog("ERROR", format, args...)

	// 推送ERROR日志到Bark
	cfg := config.GetConfig()
	if cfg != nil && cfg.Push.Bark.Host != "" && cfg.Push.Bark.Key != "" {
		// 异步推送，不阻塞日志写入
		go func() {
			bark.Push("IPTV错误", "%s", message)
		}()
	}
}

// Warn 记录警告日志
func Warn(format string, args ...interface{}) {
	writeLog("WARN", format, args...)
}

// Debug 记录调试日志
func Debug(format string, args ...interface{}) {
	writeLog("DEBUG", format, args...)
}
