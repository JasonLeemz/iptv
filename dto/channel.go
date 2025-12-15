package dto

import (
	"fmt"
	"strings"
)

// Channel 频道信息
type Channel struct {
	Name string
	URL  string
}

// IsValidURL 检查字符串是否是有效的URL
func IsValidURL(s string) bool {
	return strings.HasPrefix(s, "http://") ||
		strings.HasPrefix(s, "https://") ||
		strings.HasPrefix(s, "rtsp://") ||
		strings.HasPrefix(s, "rtmp://") ||
		strings.HasPrefix(s, "udp://") ||
		strings.HasPrefix(s, "rtp://")
}

// ConvertToM3U 转换为M3U格式
func ConvertToM3U(channels []Channel) string {
	var builder strings.Builder

	// M3U文件头
	builder.WriteString("#EXTM3U\n")

	// 添加每个频道
	for _, channel := range channels {
		// EXTINF格式: #EXTINF:-1,Channel Name
		builder.WriteString(fmt.Sprintf("#EXTINF:-1,%s\n", channel.Name))
		builder.WriteString(fmt.Sprintf("%s\n", channel.URL))
	}

	return builder.String()
}

// ConvertToCSV 转换为CSV格式：频道名称,接口地址
func ConvertToCSV(channels []Channel) string {
	var builder strings.Builder

	// 添加每个频道，格式：频道名称,接口地址
	for _, channel := range channels {
		builder.WriteString(fmt.Sprintf("%s,%s\n", channel.Name, channel.URL))
	}

	return builder.String()
}

