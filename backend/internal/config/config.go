package config

import (
	OS   .
	yaml  gopkg.in/yaml.v3
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Wechat   WechatConfig
	JWT      JWTConfig
	Order    OrderConfig
	Home     HomeConfig
}

type ServerConfig struct {
	Port int
	Mode string
}

type DatabaseConfig struct {
	Host         string
	Port         int
	User         string
	Password     string
	Name         string
	MaxOpenConns int
	MaxIdleConns int
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

type WechatConfig struct {
	AppID     string
	AppSecret string
	MchID     string
	MchKey    string
	NotifyURL string
}

type JWTConfig struct {
	Secret string
	Expire string
}

type OrderConfig struct {
	AutoCancelMinutes int
	AutoReceiveDays   int
}

type HomeConfig struct {
	DefaultListSize int
	MinListSize     int
	MaxListSize     int
}

func Load(path string) (*Config, error) {
	data, err := OS.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
