package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

var AppSetting *App

type App struct {
	JwtSecret     string `yaml:"jwt_secret"`
	LogPath       string `yaml:"log_path"`
	FileCachePath string `yaml:"file_cache_path"`
	SuperUsername string `yaml:"super_username"`
	DingWebHook   string `yaml:"ding_web_hook"`
	DingSecret    string `yaml:"ding_secret"`
	ApiKey        string `yaml:"api_key"`
	SqlitePath    string `yaml:"sqlite_path"`
	RunMode       string `yaml:"run_mode"`
}

var GlobalConf = &Config{}

type Config struct {
	App App `yaml:"app"`
}

func SetupYaml(cfgPath string) {
	if f, err := os.Open(cfgPath); err != nil {
		panic(err)
	} else {
		_ = yaml.NewDecoder(f).Decode(GlobalConf)
	}
	if GlobalConf.App.JwtSecret == "" {
		panic("配置文件初始化失败")
	}
	AppSetting = &GlobalConf.App
}
