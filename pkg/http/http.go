package http

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/go-resty/resty/v2"
	"iptv/pkg/config"
)

var (
	client *resty.Client
)

// Init 初始化HTTP客户端
func Init() error {
	cfg := config.GetConfig()
	if cfg == nil {
		return fmt.Errorf("配置未加载")
	}

	// 创建resty客户端
	client = resty.New()

	// 设置超时时间（从配置读取，默认30秒）
	timeout := 30
	if cfg.HTTP.Timeout > 0 {
		timeout = cfg.HTTP.Timeout
	}
	client.SetTimeout(time.Duration(timeout) * time.Second)

	// 设置默认请求头
	client.SetHeader("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36")
	client.SetHeader("Accept-Language", "zh-CN,zh;q=0.9")
	client.SetHeader("Dnt", "1")
	client.SetHeader("Sec-Ch-Ua", `"Google Chrome";v="143", "Chromium";v="143", "Not A(Brand";v="24"`)
	client.SetHeader("Sec-Ch-Ua-Mobile", "?0")
	client.SetHeader("Sec-Ch-Ua-Platform", `"macOS"`)

	return nil
}

// GetClient 获取HTTP客户端
func GetClient() *resty.Client {
	if client == nil {
		Init()
	}
	return client
}

// Get 执行GET请求
func Get(url string, headers map[string]string, cookies string) (*resty.Response, error) {
	req := GetClient().R()

	// 设置请求头
	if headers != nil {
		for k, v := range headers {
			req.SetHeader(k, v)
		}
	}

	// 设置Cookie
	if cookies != "" {
		req.SetHeader("Cookie", cookies)
	}

	// 执行请求
	resp, err := req.Get(url)
	if err != nil {
		return nil, fmt.Errorf("GET请求失败: %v", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("HTTP错误: %d", resp.StatusCode())
	}

	return resp, nil
}

// Post 执行POST请求
func Post(url string, body interface{}, headers map[string]string, cookies string) (*resty.Response, error) {
	req := GetClient().R()

	// 设置请求体
	if body != nil {
		req.SetBody(body)
	}

	// 设置请求头
	if headers != nil {
		for k, v := range headers {
			req.SetHeader(k, v)
		}
	}

	// 设置Cookie
	if cookies != "" {
		req.SetHeader("Cookie", cookies)
	}

	// 执行请求
	resp, err := req.Post(url)
	if err != nil {
		return nil, fmt.Errorf("POST请求失败: %v", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("HTTP错误: %d", resp.StatusCode())
	}

	return resp, nil
}

// GetBody 执行GET请求并返回响应体
func GetBody(url string, headers map[string]string, cookies string) ([]byte, error) {
	resp, err := Get(url, headers, cookies)
	if err != nil {
		return nil, err
	}
	return resp.Body(), nil
}

// GetReader 执行GET请求并返回io.Reader
func GetReader(url string, headers map[string]string, cookies string) (io.Reader, error) {
	body, err := GetBody(url, headers, cookies)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(body), nil
}

// GetHTMLHeaders 获取HTML请求的默认请求头
func GetHTMLHeaders(referer string) map[string]string {
	return map[string]string{
		"Accept":                  "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
		"Referer":                 referer,
		"Sec-Fetch-Dest":          "document",
		"Sec-Fetch-Mode":          "navigate",
		"Sec-Fetch-Site":          "same-origin",
		"Sec-Fetch-User":           "?1",
		"Upgrade-Insecure-Requests": "1",
	}
}

// GetAPIHeaders 获取API请求的默认请求头
func GetAPIHeaders(referer string) map[string]string {
	return map[string]string{
		"Accept":          "*/*",
		"Referer":         referer,
		"Priority":        "u=1, i",
		"Sec-Fetch-Dest":  "empty",
		"Sec-Fetch-Mode":  "cors",
		"Sec-Fetch-Site":  "same-origin",
		"X-Requested-With": "XMLHttpRequest",
	}
}
