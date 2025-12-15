package bark

import (
	"iptv/pkg/config"
	"os"
	"testing"
)

func TestBark(t *testing.T) {
	// 加载配置文件
	_, err := config.LoadConfig("/Users/limingze/GolandProjects/iptv/config/app.yml")
	if err != nil {
		os.Exit(1)
	}
	Push("hello", "world")
}
