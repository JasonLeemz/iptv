package main

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"iptv/dto"
	"iptv/pkg/config"
	"iptv/pkg/html"
)

// FetchChannelsFromURL 从URL获取频道列表
func FetchChannelsFromURL(pageURL string, cookies string) ([]dto.Channel, error) {
	// 解析URL参数
	parsedURL, err := url.Parse(pageURL)
	if err != nil {
		return nil, fmt.Errorf("解析URL失败: %v", err)
	}

	ip := parsedURL.Query().Get("ip")
	tk := parsedURL.Query().Get("tk")
	p := parsedURL.Query().Get("p")
	c := parsedURL.Query().Get("c")
	if c == "" {
		c = ""
	}

	// 构建getall.php的URL
	apiURL := fmt.Sprintf("https://tonkiang.us/getall.php?ip=%s&c=%s&tk=%s&p=%s",
		url.QueryEscape(ip), url.QueryEscape(c), url.QueryEscape(tk), url.QueryEscape(p))

	// 获取数据
	channels, err := fetchChannelsFromAPI(apiURL, pageURL, cookies)
	if err != nil {
		return nil, err
	}

	// 如果API返回空数据，尝试从原始页面获取
	if len(channels) == 0 {
		channels, err = fetchChannelsFromPage(pageURL, cookies)
		if err != nil {
			return nil, fmt.Errorf("从原始页面获取数据失败: %v", err)
		}
	}

	return channels, nil
}

// fetchChannelsFromAPI 从API获取频道数据
func fetchChannelsFromAPI(apiURL string, refererURL string, cookies string) ([]dto.Channel, error) {
	doc, err := html.FetchHTMLForAPI(apiURL, cookies, refererURL)
	if err != nil {
		return nil, err
	}

	// 保存debug HTML（如果启用debug模式）
	cfg := config.GetConfig()
	if cfg != nil && cfg.App.Debug {
		htmlContent, _ := doc.Html()
		debugFile := cfg.Output.Debug
		if debugFile != "" {
			os.WriteFile(debugFile, []byte(htmlContent), 0644)
		}
	}

	return parseChannelsFromDoc(doc), nil
}

// fetchChannelsFromPage 从页面获取频道数据
func fetchChannelsFromPage(pageURL string, cookies string) ([]dto.Channel, error) {
	doc, err := html.FetchHTML(pageURL, cookies, pageURL)
	if err != nil {
		return nil, err
	}

	return parseChannelsFromDoc(doc), nil
}

// parseChannelsFromDoc 从goquery文档中解析频道
func parseChannelsFromDoc(doc *goquery.Document) []dto.Channel {
	var channels []dto.Channel

	// 查找所有result块
	doc.Find("div.result").Each(func(i int, s *goquery.Selection) {
		// 提取频道名称
		channelName := ""
		tipDiv := s.Find("div.channel div.tip").First()
		if tipDiv.Length() > 0 {
			channelName = strings.TrimSpace(tipDiv.Text())
		} else {
			// 备用方法：从整个channel div提取
			channelDiv := s.Find("div.channel").First()
			if channelDiv.Length() > 0 {
				channelName = strings.TrimSpace(channelDiv.Text())
				// 移除HTML标签
				channelName = cleanTextFromHTML(channelName)
			}
		}

		// 跳过无效的频道名称
		if channelName == "" ||
			strings.Contains(channelName, "请使用搜索框") ||
			strings.Contains(channelName, "验证") ||
			strings.Contains(channelName, "来自") ||
			strings.Contains(channelName, "组播源") {
			return
		}

		// 提取URL
		m3u8URL := ""
		
		// 方法1: 从onclick属性中提取
		s.Find("img[onclick*='copyto']").Each(func(i int, img *goquery.Selection) {
			onclick, exists := img.Attr("onclick")
			if exists {
				// 提取 onclick=copyto('URL') 中的URL
				if strings.Contains(onclick, "copyto('") {
					start := strings.Index(onclick, "copyto('") + len("copyto('")
					end := strings.Index(onclick[start:], "'")
					if end > 0 {
						m3u8URL = strings.TrimSpace(onclick[start : start+end])
					}
				}
			}
		})

		// 方法2: 从m3u8 div中的td标签提取
		if m3u8URL == "" {
			s.Find("div.m3u8 td").Each(func(i int, td *goquery.Selection) {
				text := strings.TrimSpace(td.Text())
				if strings.HasPrefix(text, "http://") || strings.HasPrefix(text, "https://") {
					// 提取第一个URL
					parts := strings.Fields(text)
					for _, part := range parts {
						if strings.HasPrefix(part, "http://") || strings.HasPrefix(part, "https://") {
							m3u8URL = part
							break
						}
					}
				}
			})
		}

		// 验证URL格式
		if !dto.IsValidURL(m3u8URL) {
			return
		}

		if channelName != "" && m3u8URL != "" {
			channels = append(channels, dto.Channel{
				Name: channelName,
				URL:  m3u8URL,
			})
		}
	})

	return channels
}

// cleanTextFromHTML 清理HTML文本
func cleanTextFromHTML(text string) string {
	// 移除常见的HTML标签内容
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.ReplaceAll(text, "\t", " ")
	
	// 合并多个空格
	spaceRegex := regexp.MustCompile(`\s+`)
	text = spaceRegex.ReplaceAllString(text, " ")
	
	return strings.TrimSpace(text)
}

// AggregateChannelsToM3U 汇总频道到M3U格式
func AggregateChannelsToM3U(channels []dto.Channel, outputPath string) error {
	// 确保输出目录存在
	dir := filepath.Dir(outputPath)
	if dir != "." && dir != "" {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return fmt.Errorf("创建输出目录失败: %v", err)
		}
	}

	content := dto.ConvertToM3U(channels)
	err := os.WriteFile(outputPath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("写入M3U文件失败: %v", err)
	}

	return nil
}

// AggregateChannelsToTXT 汇总频道到TXT格式
func AggregateChannelsToTXT(channels []dto.Channel, outputPath string) error {
	// 确保输出目录存在
	dir := filepath.Dir(outputPath)
	if dir != "." && dir != "" {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return fmt.Errorf("创建输出目录失败: %v", err)
		}
	}

	content := dto.ConvertToCSV(channels)
	err := os.WriteFile(outputPath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("写入TXT文件失败: %v", err)
	}

	return nil
}
