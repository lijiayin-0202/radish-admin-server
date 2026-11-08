package config

import "time"

// RideHailing 保存后台管理服务访问模块3微服务的连接配置。
type RideHailing struct {
	BaseURL        string        `mapstructure:"base-url" json:"base-url" yaml:"base-url"`
	RequestTimeout time.Duration `mapstructure:"request-timeout" json:"request-timeout" yaml:"request-timeout"`
}
