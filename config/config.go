package config

import (
	"time"

	"github.com/spf13/viper"
)

// Configuration 项目配置
type Configuration struct {
	// gpt apikey
	ApiKey string `json:"api_key"`
	// 自动通过好友
	AutoPass bool `json:"auto_pass"`
	// 会话超时时间
	SessionTimeout time.Duration `json:"session_timeout"`
	// GPT请求最大字符数
	MaxTokens int `json:"max_tokens"`
	// GPT模型
	Model string `json:"model"`
	// 热度
	Temperature float32 `json:"temperature"`
	// 回复前缀
	ReplyPrefix string `json:"reply_prefix"`
	// 清空会话口令
	SessionClearToken string `json:"session_clear_token"`
	// Proxy Forward-proxy
	Proxy string `json:"proxy,omitempty"`

	// PrivateTrigger private trigger words
	PrivateTrigger string `json:"private_trigger"`
}

// LoadConfig 加载配置
func LoadConfig() *Configuration {
	viper.AddConfigPath("conf")
	viper.SetConfigType("json")
	viper.SetConfigName("config")
	viper.ReadInConfig()

	var config *Configuration
	viper.Unmarshal(&config)

	return config
}
