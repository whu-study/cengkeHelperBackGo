package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

// Conf 定义配置结构体
var Conf struct {
	Mysql struct {
		Host        string `yaml:"host"`
		Port        string `yaml:"port"`
		User        string `yaml:"user"`
		Password    string `yaml:"password"`
		Database    string `yaml:"database"`
		AutoMigrate bool   `yaml:"auto_migrate"`
	} `yaml:"mysql"`
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
}

// LoadConfig 加载配置文件
func LoadConfig(filePath string) bool {
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("failed to read config file: ", filePath, err)
		return false
	}

	err = yaml.Unmarshal(data, &Conf)
	if err != nil {
		fmt.Println("failed to unmarshal config file: ", err)
		return false
	}
	return true
}

func init() {
	if LoadConfig("config.yaml") ||
		LoadConfig("internal/config/config.yaml") {
		fmt.Println("Conf loaded successfully: \n", Conf)
	} else {
		panic("Failed to load config")
	}
}
