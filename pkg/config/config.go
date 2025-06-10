package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

var AppSetting *App

type App struct {
	JwtSecret     string `yaml:"jwtSecret"`
	LogSavePath   string `yaml:"log_save_path"`
	FileCachePath string `yaml:"file_cache_path"`
	DomainServer  string `yaml:"domain_server"`
	DomainClient  string `yaml:"domain_client"`
	SuperUsername string `yaml:"super_username"`
	DingWebHook   string `yaml:"ding_web_hook"`
	DingSecret    string `yaml:"ding_secret"`
	MailHost      string `yaml:"mail_host"`
	MailAddress   string `yaml:"mail_address"`
	MailPassword  string `yaml:"mail_password"`
	ApiKey        string `yaml:"api_key"`
	SqlitePath    string `yaml:"sqlite_path"`
}

var RedisSetting *RedisConfig

type RedisConfig struct {
	Addr              string `yaml:"addr"`
	Password          string `yaml:"password"`
	DBIndex           int    `yaml:"db_index"`
	MaxIdle           int    `yaml:"max_idle"`
	MaxActive         int    `yaml:"max_active"`
	IdleTimeoutSecond int    `yaml:"idle_timeout"`
}

var ServerSetting *Server

type Server struct {
	TaskSwitch         string `yaml:"task_switch"`
	RunMode            string `yaml:"run_mode"`
	HttpPort           int    `yaml:"http_port"`
	ReadTimeoutSecond  int64  `yaml:"read_timeout"`
	WriteTimeoutSecond int64  `yaml:"write_timeout"`
}

var GlobalConf = &Config{}

type Config struct {
	App    App         `yaml:"app"`
	Server Server      `yaml:"server" `
	Redis  RedisConfig `yaml:"redis" `
}

func SetupYaml() {
	cfgPath := "conf/app.yaml"
	if f, err := os.Open(cfgPath); err != nil {
		panic(err)
	} else {
		_ = yaml.NewDecoder(f).Decode(GlobalConf)
	}
	if GlobalConf.App.JwtSecret == "" {
		panic("配置文件初始化失败")
	}

	AppSetting = &GlobalConf.App
	ServerSetting = &GlobalConf.Server
	RedisSetting = &GlobalConf.Redis
}
