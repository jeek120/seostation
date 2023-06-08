package config

import (
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Sites     map[string]Site `yaml:"sites"`
	Listen    string          `yaml:"listen"`
	TlsListen string          `yaml:"tlsListen"`
}

type Site struct {
	ProxyHost []ProxyHost   `yaml:"proxyHosts"`
	Replaces  yaml.MapSlice `yaml:"replaces"`
}

type ProxyHost struct {
	Before string `yaml:"oldHost"`
	After  string `yaml:"newHost"`
	Schema string `yaml:"schema"`
}

func NewLocalSample() *Config {
	return &Config{
		Sites: map[string]Site{
			"4月天": {
				ProxyHost: []ProxyHost{
					{
						Before: "local1.com:8080",
						After:  "4yt.net",
						Schema: "https",
					},
				},
				Replaces: []yaml.MapItem{
					{Key: "local1.com:8080/", Value: "4yt.net/"},
				},
			},
		},
		Listen: ":8080",
	}
}

func Get() *Config {
	var c Config
	bs, err := os.ReadFile("config.yaml")
	if err != nil {
		if os.ErrNotExist == err || strings.Contains(err.Error(), "no such file") {
			return create()
		}
		log.Printf("读取配置文件失败: %s", err)
		panic(err)
	}
	err = yaml.Unmarshal(bs, &c)
	if err != nil {
		log.Printf("解析yaml错误: %s", err)
		panic(err)
	}
	return &c
}

func create() *Config {
	c := NewLocalSample()
	bs, err := yaml.Marshal(c)
	if err != nil {
		log.Printf("序列话出错: %s", err)
		panic(err)
	}
	err = os.WriteFile("config.yaml", bs, 0644)
	if err != nil {
		log.Printf("写入配置文件失败: %s", err)
	}
	return c
}
