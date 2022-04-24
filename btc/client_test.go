package btc

import (
	"fmt"
	"github.com/btcsuite/btcutil"
	"testing"
)

func TestClient_ListUnspent(t *testing.T) {

	var (
		//host = "47.243.189.44:30000"
		//user = "BTCuser"
		//pass = "BTCpassword"
		//network = "MainNet"

		host       = "127.0.0.1:18332"
		user       = "admin"
		pass       = "123456"
		network    = "TestNet3"
		defaultNet = getNetwork(network)
	)
	//address := "2N8Ah8xiqrfb37aktkkwYmMqpDGQJMSDAWY"
	address := "2NGU5d9Y3uKLJXZa6GgYm2bf4wNK33GcAwJ"
	client, err := NewClient(host, user, pass)
	if err != nil {
		t.Error("创建client失败：", err)
	}
	fromAddress, err := btcutil.DecodeAddress(address, defaultNet)
	if err != nil {
		t.Error("地址错误：", err)
	}

	list, err := client.ListUnspent(fromAddress)
	if err != nil {
		t.Error("获取未交易的UTXO失败：", err)
	}

	fmt.Println(list)

}

func TestClient_GetBlockHash(t *testing.T) {

	var (
		host = "127.0.0.1:18332"
		user = "admin"
		pass = "123456"
	)
	client, err := NewClient(host, user, pass)
	if err != nil {
		t.Error("创建client失败：", err)
	}

	blockNumber := int64(2000000)
	hash, err := client.GetBlockHash(blockNumber)
	if err != nil {
		t.Error("获取区块哈希失败：", err)
	}

	fmt.Println(hash)

}

func TestClient_ListSinceBlock(t *testing.T) {

	var (
		//host = "47.243.189.44:30000"
		//user = "BTCuser"
		//pass = "BTCpassword"

		host = "127.0.0.1:18332"
		user = "admin"
		pass = "123456"
	)
	client, err := NewClient(host, user, pass)
	if err != nil {
		t.Error("创建client失败：", err)
	}

	blockNumber := int64(600000)
	hash, err := client.GetBlockHash(blockNumber)
	if err != nil {
		t.Error("获取区块哈希失败：", err)
	}

	fmt.Println(hash)

	list, err := client.ListSinceBlock(hash)
	if err != nil {
		t.Error("获取交易列表失败：", err)
	}

	fmt.Println(list)

}

func TestClient_ListTransactions(t *testing.T) {

	var (
		//host = "47.243.189.44:30000"
		//user = "BTCuser"
		//pass = "BTCpassword"
		host = "127.0.0.1:18332"
		user = "admin"
		pass = "123456"
	)
	client, err := NewClient(host, user, pass)
	if err != nil {
		t.Error("创建client失败：", err)
	}

	address := "2N8Ah8xiqrfb37aktkkwYmMqpDGQJMSDAWY"
	list, err := client.ListTransactions(address)
	if err != nil {
		t.Error("获取最近发生的钱包交易失败：", err)
	}

	fmt.Println(list)

}

func TestClient_ImportAddress(t *testing.T) {

	var (
		host = "127.0.0.1:18332"
		user = "admin"
		pass = "123456"
	)
	client, err := NewClient(host, user, pass)
	if err != nil {
		t.Error("创建client失败：", err)
	}

	//address := "2N8Ah8xiqrfb37aktkkwYmMqpDGQJMSDAWY"
	address := "2NGU5d9Y3uKLJXZa6GgYm2bf4wNK33GcAwJ"
	err = client.ImportAddress(address)
	if err != nil {
		t.Error("导入地址失败：", err)
	}

}
