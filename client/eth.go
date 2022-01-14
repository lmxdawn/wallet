package client

import "github.com/ethereum/go-ethereum/ethclient"

func NewEthClient(url string) *ethclient.Client {
	client, _ := ethclient.Dial(url)
	return client
}
