package config

import (
	"encoding/json"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

// Conf 定义配置结构体
var Conf struct {
	Mysql struct {
		Host        string `yaml:"host" json:"host"`
		Port        string `yaml:"port" json:"port"`
		User        string `yaml:"user" json:"user"`
		Password    string `yaml:"password" json:"password"`
		Database    string `yaml:"database" json:"database"`
		AutoMigrate bool   `yaml:"auto_migrate" json:"autoMigrate"`
	} `yaml:"mysql" json:"mysql"`
	Server struct {
		Port string `yaml:"port" json:"port"`
	} `yaml:"server" json:"server"`
	JwtSecurityKey string `yaml:"jwt_security_key" json:"jwtSecurityKey"`
}

// LoadConfig 加载配置文件
func LoadConfig(filePath string) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal("failed to read config file: ", filePath, err)
	}

	err = yaml.Unmarshal(data, &Conf)
	if err != nil {
		log.Fatal("failed to unmarshal config file: ", err)
	}
}

func WatchConfigLoop(filePath string) bool {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}

	LoadConfig(filePath)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer func(watcher *fsnotify.Watcher) {
		err := watcher.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(watcher)

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					fmt.Println("Config file modified, reloading...")
					LoadConfig(filePath)

				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Printf("Watcher error: %v", err)
			}
		}
	}()

	err = watcher.Add(filePath)
	if err != nil {
		log.Fatal(err)
	}

	return true
}

func init() {
	if WatchConfigLoop("config.yaml") ||
		WatchConfigLoop("internal/config/config.yaml") ||
		WatchConfigLoop("~/config.yaml") {
		confBytes, _ := json.MarshalIndent(Conf, " ", "    ")
		fmt.Println("Conf loaded successfully: \n", string(confBytes))
	} else {
		panic("Failed to load config")
	}
}
