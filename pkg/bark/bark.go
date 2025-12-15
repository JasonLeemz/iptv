package bark

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"iptv/pkg/config"
)

// Push 推送消息到Bark（支持格式化字符串）
func Push(title string, bodyFormat string, args ...interface{}) error {
	// 格式化body内容
	body := fmt.Sprintf(bodyFormat, args...)
	cfg := config.GetConfig()
	if cfg == nil {
		return fmt.Errorf("配置未加载")
	}

	barkHost := cfg.Push.Bark.Host
	barkKey := cfg.Push.Bark.Key

	if barkHost == "" || barkKey == "" {
		return fmt.Errorf("bark配置不完整")
	}

	// 构建推送URL
	// 格式: https://bark.abcd.xyz/Y2MVwieFC/标题/内容
	pushURL := fmt.Sprintf("%s/%s/%s/%s",
		strings.TrimSuffix(barkHost, "/"),
		url.QueryEscape(barkKey),
		url.QueryEscape(title),
		url.QueryEscape(body))

	// 发送HTTP GET请求
	resp, err := http.Get(pushURL)
	if err != nil {
		return fmt.Errorf("发送Bark推送失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取Bark响应失败: %v", err)
	}

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bark推送失败: HTTP %d, %s", resp.StatusCode, string(responseBody))
	}

	return nil
}
