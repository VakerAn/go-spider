package config

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"gopkg.in/yaml.v2"
	"os"
	"time"
)

type Config struct {
	App     *App     `yaml:"app" comment:"the app settings"`
	Server  *Server  `yaml:"server" comment:"The server settings"`
	RedisDB *Redis   `yaml:"redis" comment:"The redis settings"`
	MysqlDB *Mysql   `yaml:"mysql" comment:"The  db settings"`
	Logger  *Logger  `yaml:"logger" comment:"The log settings"`
	BaseUrl *BaseURL `yaml:"base_url" comment:"its api domain name"`
}
type App struct {
	SpiderPath     string `yaml:"spider_path"`
	SpiderPathName string `yaml:"spider_path_name"`
	DebugPath      string `yaml:"debug_path"`
	DebugPathName  string `yaml:"debug_path_name"`
	SpiderMod      string `yaml:"spider_mod"`
	DebugMod       string `yaml:"debug_mod"`
}
type Server struct {
	RunMode       string        `yaml:"run_mode"`
	HTTPPort      int           `yaml:"http_port"`
	ReadTimeout   time.Duration `yaml:"read_timeout"`
	WriteTimeout  time.Duration `yaml:"write_timeout"`
	ServerTimeout time.Duration `yaml:"server_timeout"`
}

type Redis struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Scheme   string `yaml:"scheme"`
	Alias    string `yaml:"alias"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Ttl      int    `yaml:"ttl"`
}

type Mysql struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Scheme   string `yaml:"scheme"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type Logger struct {
	DefaultName string `yaml:"default_name"`
	Path        string `yaml:"path"`
}

type BaseURL struct {
	Douban   string `yaml:"douban"`
	Tiankong string `yaml:"tiankong"`
	Tencent  string `yaml:"tencent"`
}

var ConfData *Config

func InitConfig() {
	yamlFile, err := os.ReadFile("config/config.yml")
	if err != nil {
		fmt.Printf("yaml is not found：%s\n", err.Error())
	}
	var result map[string]interface{}
	err = yaml.Unmarshal(yamlFile, &result)
	jsonData, err := jsoniter.Marshal(result)
	if err != nil {
		fmt.Printf("jsoniter is not marshal：%s\n", err.Error())
	}
	ConfData = &Config{}
	err = yaml.Unmarshal([]byte(jsonData), ConfData)
	if err != nil {
		fmt.Printf("yaml is not Unmarshal：%s\n", err.Error())
	}
}
