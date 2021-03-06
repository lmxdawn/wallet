package config

import (
	"github.com/jinzhu/configor"
)

type AppConfig struct {
	Port uint `yaml:"port"`
}

type EngineConfig struct {
	CoinName            string `yaml:"coin_name"`             // 币种名称
	Contract            string `yaml:"contract"`              // 合约地址（为空表示主币）
	ContractType        string `yaml:"contract_type"`         // 合约类型（波场需要区分是TRC20还是TRC10）
	Protocol            string `yaml:"protocol"`              // 协议名称
	Network             string `yaml:"network"`               // 网络名称（暂时BTC协议有用{MainNet：主网，TestNet：测试网，TestNet3：测试网3，SimNet：测试网}）
	Rpc                 string `yaml:"rpc"`                   // rpc配置
	User                string `yaml:"user"`                  // rpc用户名（没有则为空）
	Pass                string `yaml:"pass"`                  // rpc密码（没有则为空）
	File                string `yaml:"file"`                  // db文件路径配置
	WalletPrefix        string `yaml:"wallet_prefix"`         // 钱包的存储前缀
	HashPrefix          string `yaml:"hash_prefix"`           // 交易哈希的存储前缀
	BlockInit           uint64 `yaml:"block_init"`            // 初始块（默认读取最新块）
	BlockAfterTime      uint64 `yaml:"block_after_time"`      // 获取最新块的等待时间
	ReceiptCount        uint64 `yaml:"receipt_count"`         // 交易凭证worker数量
	ReceiptAfterTime    uint64 `yaml:"receipt_after_time"`    // 获取交易信息的等待时间
	CollectionAfterTime uint64 `yaml:"collection_after_time"` // 归集等待时间
	CollectionCount     uint64 `yaml:"collection_count"`      // 归集发送worker数量
	CollectionMax       string `yaml:"collection_max"`        // 最大的归集数量（满足多少才归集，为0表示不自动归集）
	CollectionAddress   string `yaml:"collection_address"`    // 归集地址
	Confirms            uint64 `yaml:"confirms"`              // 确认数量
	RechargeNotifyUrl   string `yaml:"recharge_notify_url"`   // 充值通知回调地址
	WithdrawNotifyUrl   string `yaml:"withdraw_notify_url"`   // 提现通知回调地址
	WithdrawPrivateKey  string `yaml:"withdraw_private_key"`  // 提现的私钥地址
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
			return config, err
		}
	} else {
		err := configor.Load(&config, "config/config-example.yml")
		if err != nil {
			return config, err
		}
	}
	return config, nil
}
