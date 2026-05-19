package config

import (
	"fmt"
	"log"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

var globalConfig *Config

type Config struct {
	Port     string `mapstructure:"port"`
	AutoType bool   `mapstructure:"auto_type"`
	Session  string `mapstructure:"session"`
}

func setDefaults() {
	viper.SetDefault("port", "2828")     // 默认端口
	viper.SetDefault("auto_type", false) // 默认不启用自动输入
	viper.SetDefault("session", "")      // 默认 session 为空
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path) // 如 "config.yml"
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()               // 允许环境变量覆盖
	viper.SetEnvPrefix("DOUBAO_INPUT") // 环境变量前缀

	// 先设置默认值
	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		err := viper.WriteConfigAs(path)
		if err != nil {
			return nil, fmt.Errorf("创建默认配置文件失败: %w", err)
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func InitConfig() {
	cfg, err := LoadConfig("./doubao-input-config.yml")
	if err != nil {
		log.Fatal("加载配置失败: ", err)
	}
	globalConfig = cfg
}

func GetConfig() *Config {
	return globalConfig
}

func SaveConfig(cfg *Config) error {
	// 将结构体转为 map
	var m map[string]interface{}
	if err := mapstructure.Decode(cfg, &m); err != nil {
		return err
	}
	// 遍历设置
	for k, v := range m {
		viper.Set(k, v)
	}

	return viper.WriteConfig()
}
