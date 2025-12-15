package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	App struct {
		Debug bool `yaml:"debug"`
	} `yaml:"app"`
	MulticastIP struct {
		Limit int `yaml:"limit"`
	} `yaml:"multicastIP"`
	Cookie struct {
		Data string `yaml:"data"`
	} `yaml:"cookie"`
	Crontab struct {
		Enable bool   `yaml:"enable"`
		Job    string `yaml:"job"`
	} `yaml:"crontab"`
	Output struct {
		M3U   string `yaml:"m3u"`
		Local string `yaml:"local"`
		Debug string `yaml:"debug"`
	} `yaml:"output"`
	Log struct {
		Path string `yaml:"path"`
	} `yaml:"log"`
	Push struct {
		Bark struct {
			Host string `yaml:"host"`
			Key  string `yaml:"key"`
		} `yaml:"bark"`
	} `yaml:"push"`
	RedirectOutput struct {
		Enable bool   `yaml:"enable"`
		Move   string `yaml:"move"`
		To     string `yaml:"to"`
	} `yaml:"redirectOutput"`
}

var globalConfig *Config

// LoadConfig 加载配置文件
func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	globalConfig = &config
	return &config, nil
}

// GetConfig 获取全局配置
func GetConfig() *Config {
	return globalConfig
}

