package btc

import (
	"github.com/btcsuite/btcd/chaincfg"
	"strings"
)

func getNetwork(network string) *chaincfg.Params {
	var defaultNet *chaincfg.Params

	//指定网络 {MainNet：主网，TestNet：测试网，TestNet3：测试网3，SimNet：测试网}
	switch strings.ToLower(network) {
	case "mainnet":
		defaultNet = &chaincfg.MainNetParams
	case "testnet":
		defaultNet = &chaincfg.RegressionNetParams
	case "testnet3":
		defaultNet = &chaincfg.TestNet3Params
	case "simnet":
		defaultNet = &chaincfg.SimNetParams
	}
	return defaultNet
}
