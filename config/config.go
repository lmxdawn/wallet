package config

import (
	"github.com/jinzhu/configor"
	"github.com/pkg/errors"
)

type AppConfig struct {
	Port uint `yaml:"port"`
}

type EngineConfig struct {
	Name               string `yaml:"name"`                 // 名称
	Protocol           string `yaml:"protocol"`             // 协议名称
	Rpc                string `yaml:"server"`               // rpc配置
	File               string `yaml:"file"`                 // db文件配置
	BlockInit          uint64 `yaml:"block_init"`           // 初始块
	BlockCount         uint64 `yaml:"block_count"`          // 区块worker数量
	BlockAfterTime     uint64 `yaml:"block_after_time"`     // 获取最新块的等待时间
	ReceiptCount       uint64 `yaml:"receipt_count"`        // 交易凭证worker数量
	ReceiptAfterTime   uint64 `yaml:"receipt_after_time"`   // 获取交易信息的等待时间
	Confirms           uint64 `yaml:"confirms"`             // 确认数量
	RechargeNotifyUrl  string `yaml:"recharge_notify_url"`  // 充值通知回调地址
	WithdrawNotifyUrl  string `yaml:"withdraw_notify_url"`  // 提现通知回调地址
	WithdrawPrivateKey string `yaml:"withdraw_private_key"` // 提现的私钥地址
	Decimals           int    `yaml:"decimals"`             // 精度
}

type Config struct {
	App     AppConfig
	Engines []EngineConfig
}

func NewConfig(confPath string) (Config, error) {
	var config Config
	if confPath != "" {
		err := configor.Load(&config, confPath)
		if err != nil {
			return config, errors.WithStack(err)
		}
	} else {
		err := configor.Load(&config, "config/config-example.yml")
		if err != nil {
			return config, errors.WithStack(err)
		}
	}
	return config, nil
}
