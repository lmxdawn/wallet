package main

import (
	"github.com/lmxdawn/wallet/engine"
	"github.com/rs/zerolog/log"
)

func Start() error {
	rpcUrl := "https://mainnet.infura.io/v3/adee4ead47844d238802431fcb7683c6"
	httpUrl := "https://mainnet.infura.io"
	eth, err := engine.NewEthEngine("eth", rpcUrl, httpUrl, "path/to/db", 1, 1, 12)
	if err != nil {
		log.Error().Msgf("eth run err：%v", err)
	}
	ethErc20, err := engine.NewEthEngine("eth", rpcUrl, httpUrl, "path/to/db", 1, 1, 12)
	if err != nil {
		log.Error().Msgf("eth run err：%v", err)
	}

	num := int64(13969128)

	eth.Run(num)
	ethErc20.Run(num)

	return nil

}
