package main

import (
	"github.com/lmxdawn/wallet/db"
	"github.com/lmxdawn/wallet/engine"
)

func Start() error {
	url := "https://mainnet.infura.io/v3/adee4ead47844d238802431fcb7683c6"
	ethDB, err := db.NewKeyDB("path/to/db")
	if err != nil {
		return err
	}

	ethErc20DB, err := db.NewKeyDB("path/to/db")
	if err != nil {
		return err
	}

	eth := engine.NewEthEngine(1, 1, 12, url, ethDB)
	ethErc20 := engine.NewEthEngine(1, 1, 12, url, ethErc20DB)

	num := int64(13969128)

	eth.Run(num)
	ethErc20.Run(num)

	return nil

}
