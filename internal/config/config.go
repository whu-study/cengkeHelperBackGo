package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v3"
	"log"
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
	//if LoadConfig("config.yaml") ||
	//	LoadConfig("internal/config/config.yaml") {
	//	fmt.Println("Conf loaded successfully: \n", Conf)
	//} else {
	//	panic("Failed to load config")
	//}
	if WatchConfigLoop("config.yaml") ||
		WatchConfigLoop("internal/config/config.yaml") {
		fmt.Println("Conf loaded successfully: \n", Conf)
	} else {
		panic("Failed to load config")
	}
}
