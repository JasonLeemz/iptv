package html

import (
	"bytes"
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"
	httppkg "iptv/pkg/http"
)

// FetchHTML 获取HTML内容
func FetchHTML(url string, cookies string, referer string) (*goquery.Document, error) {
	headers := httppkg.GetHTMLHeaders(referer)
	body, err := httppkg.GetBody(url, headers, cookies)
	if err != nil {
		return nil, err
	}

	return goquery.NewDocumentFromReader(bytes.NewReader(body))
}

// FetchHTMLForAPI 获取API的HTML内容（用于getall.php等API请求）
func FetchHTMLForAPI(url string, cookies string, referer string) (*goquery.Document, error) {
	headers := httppkg.GetAPIHeaders(referer)
	body, err := httppkg.GetBody(url, headers, cookies)
	if err != nil {
		return nil, err
	}

	return goquery.NewDocumentFromReader(bytes.NewReader(body))
}

// FetchHTMLRaw 获取原始HTML内容（返回io.Reader）
func FetchHTMLRaw(url string, cookies string, referer string) (io.Reader, error) {
	headers := httppkg.GetHTMLHeaders(referer)
	return httppkg.GetReader(url, headers, cookies)
}

// ExtractLinks 从文档中提取链接
func ExtractLinks(doc *goquery.Document, selector string) []string {
	var links []string
	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			links = append(links, href)
		}
	})
	return links
}

// ExtractText 从文档中提取文本
func ExtractText(doc *goquery.Document, selector string) []string {
	var texts []string
	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if text != "" {
			texts = append(texts, text)
		}
	})
	return texts
}

