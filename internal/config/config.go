package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

var globalConfig *Config

type Config struct {
	Port              string `mapstructure:"port"`
	AutoType          bool   `mapstructure:"auto_type"`
	Startup           bool   `mapstructure:"startup"`
	Session           string `mapstructure:"session"`
	ConversationLimit int    `mapstructure:"conversation_limit"`
	IntervalTime      int    `mapstructure:"interval_time"`
	ConversationID    string `mapstructure:"conversation_id"`
}

func setDefaults() {
	viper.SetDefault("port", "2828")          // 默认端口
	viper.SetDefault("auto_type", true)       // 默认启用自动输入
	viper.SetDefault("startup", false)        // 默认不开机自启
	viper.SetDefault("session", "")           // 默认 session 为空
	viper.SetDefault("conversation_limit", 5) // 单次获取对话数量
	viper.SetDefault("interval_time", 1000)   // 默认请求间隔时间（毫秒）
	viper.SetDefault("conversation_id", "")   // 对话 ID，留空自动从 curl 中提取
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()                 // 允许环境变量覆盖
	viper.SetEnvPrefix("DOUBAO_SUSURRO") // 环境变量前缀

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

func InitConfig(configPath string) {
	// 如果没有传入路径，使用默认路径
	if configPath == "" {
		// 默认使用可执行文件所在目录作为配置文件路径
		exePath, err := os.Executable()
		if err != nil {
			log.Fatal("获取执行文件路径失败: ", err)
		}
		exeDir := filepath.Dir(exePath)
		configPath = filepath.Join(exeDir, "Doubao-Susurro-config.yml")
	}

	cfg, err := LoadConfig(configPath)
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
