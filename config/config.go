package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

// Config 配置结构
type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port         int `json:"port"`
	ReadTimeout  int `json:"readTimeout"`
	WriteTimeout int `json:"writeTimeout"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Type     string `json:"type"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
}

// 默认配置
var defaultConfig = Config{
	Server: ServerConfig{
		Port:         8080,
		ReadTimeout:  60,
		WriteTimeout: 60,
	},
	Database: DatabaseConfig{
		Type:     "sqlite3",
		Host:     "localhost",
		Port:     3306,
		User:     "root",
		Password: "password",
		DBName:   "goblog",
	},
}

// LoadConfig 加载配置
func LoadConfig() *Config {
	// 检查配置文件是否存在
	configPath := "config.json"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// 配置文件不存在，创建默认配置
		saveDefaultConfig(configPath)
		return &defaultConfig
	}

	// 读取配置文件
	configFile, err := os.Open(configPath)
	if err != nil {
		log.Printf("无法打开配置文件: %v，使用默认配置", err)
		return &defaultConfig
	}
	defer configFile.Close()

	// 解析配置
	var config Config
	decoder := json.NewDecoder(configFile)
	if err := decoder.Decode(&config); err != nil {
		log.Printf("解析配置文件失败: %v，使用默认配置", err)
		return &defaultConfig
	}

	return &config
}

// 保存默认配置到文件
func saveDefaultConfig(path string) {
	// 确保目录存在
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Printf("创建配置目录失败: %v", err)
		return
	}

	// 创建配置文件
	file, err := os.Create(path)
	if err != nil {
		log.Printf("创建配置文件失败: %v", err)
		return
	}
	defer file.Close()

	// 写入默认配置
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(defaultConfig); err != nil {
		log.Printf("写入默认配置失败: %v", err)
	}
}
