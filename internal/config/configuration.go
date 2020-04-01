package config

import (
	"flag"
	"io/ioutil"
	stdlog "log"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v2"
)

var (
	globalConfiguration *GlobalConfiguration
	once                sync.Once
)

// GlobalConfiguration 程式內取用的系統設定參數 只暴露程式內需要用到的參數
type GlobalConfiguration struct {
	config Configuration
}

// Database 用來提供連線的資料庫數據
type Database struct {
	Name     string
	Address  string
	Username string
	Password string
	DBName   string
}

// Configuration 用來代表 config 設定物件
type Configuration struct {
	Databases []Database
}

// New function 創建一個 configuration instance 出來
func New(fileName string) *Configuration {
	flag.Parse()
	c := Configuration{}

	configPath := filepath.Join("configs", fileName)
	_, err := os.Stat(configPath)
	if err != nil {
		stdlog.Fatalf("config: file error: %s", err.Error())
	}

	// config exists
	file, err := ioutil.ReadFile(filepath.Clean(configPath))
	if err != nil {
		stdlog.Fatalf("config: read file error: %s", err.Error())
	}

	err = yaml.Unmarshal(file, &c)
	if err != nil {
		stdlog.Fatal("config: yaml unmarshal error:", err)
	}

	return &c
}

// SetGlobalConfiguration 設定程式內取用的系統設定參數
func SetGlobalConfiguration(cfg *Configuration) {
	once.Do(func() {
		globalConfiguration = &GlobalConfiguration{*cfg}
	})
	return
}

// GetConfiguration 取得目前的Configuration
func GetConfiguration() Configuration {
	return globalConfiguration.config
}
