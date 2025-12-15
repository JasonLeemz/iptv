package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"iptv/pkg/config"
	"iptv/pkg/html"
)

// MulticastSource 组播源信息
type MulticastSource struct {
	IP  string
	URL string
}

// FetchMulticastIPs 从iptvmulticast.php获取组播源IP列表
func FetchMulticastIPs(cookies string) ([]MulticastSource, error) {
	// 从配置获取limit，如果未配置则使用默认值5
	limit := 5
	cfg := config.GetConfig()
	if cfg != nil && cfg.MulticastIP.Limit > 0 {
		limit = cfg.MulticastIP.Limit
	}
	multicastURL := "https://tonkiang.us/iptvmulticast.php"
	referer := "https://tonkiang.us/?"

	doc, err := html.FetchHTML(multicastURL, cookies, referer)
	if err != nil {
		return nil, fmt.Errorf("获取组播源页面失败: %v", err)
	}

	var sources []MulticastSource
	count := 0

	// 查找 class="channel" 下的链接，且只获取 p=2 的链接
	// 格式: <div class="channel"><a href='channellist.html?ip=180.127.29.153&tk=03aae295&p=2' title="Channel List">
	doc.Find("div.channel a[href*='channellist.html?ip=']").Each(func(i int, s *goquery.Selection) {
		if count >= limit {
			return
		}

		href, exists := s.Attr("href")
		if !exists {
			return
		}

		// 只获取 p=2 的链接
		if !strings.Contains(href, "p=2") {
			return
		}

		// 提取IP地址
		ip := extractIPFromHref(href)
		if ip == "" {
			return
		}

		// 构建完整URL
		fullURL := buildFullURL(href)
		if fullURL == "" {
			return
		}

		sources = append(sources, MulticastSource{
			IP:  ip,
			URL: fullURL,
		})

		count++
	})

	return sources, nil
}

// extractIPFromHref 从href中提取IP地址
func extractIPFromHref(href string) string {
	// href格式: channellist.html?ip=221.220.131.129&tk=c188c009&p=2
	if !strings.Contains(href, "ip=") {
		return ""
	}

	parts := strings.Split(href, "ip=")
	if len(parts) < 2 {
		return ""
	}

	ipPart := strings.Split(parts[1], "&")[0]
	return strings.TrimSpace(ipPart)
}

// buildFullURL 构建完整的URL
func buildFullURL(href string) string {
	// 如果已经是完整URL，直接返回
	if strings.HasPrefix(href, "http://") || strings.HasPrefix(href, "https://") {
		return href
	}

	// 构建完整URL
	if strings.HasPrefix(href, "channellist.html") {
		return "https://tonkiang.us/" + href
	}

	return ""
}

// UpdateSourceFile 更新config/source.txt文件
func UpdateSourceFile(sources []MulticastSource, configPath string) error {
	filePath := filepath.Join(configPath, "source.txt")

	// 构建文件内容
	var lines []string
	for _, source := range sources {
		lines = append(lines, source.URL)
	}

	content := strings.Join(lines, "\n") + "\n"

	// 写入文件
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("写入source.txt失败: %v", err)
	}

	return nil
}
