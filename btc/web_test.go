package btc

import (
	"fmt"
	"testing"
)

func TestWeb_Transfer(t *testing.T) {

	host := "127.0.0.1:18332"
	user := "admin"
	pass := "123456"
	network := "TestNet3"
	web, err := NewWeb(host, user, pass, network)
	if err != nil {
		t.Error(err)
	}

	address := "2N8Ah8xiqrfb37aktkkwYmMqpDGQJMSDAWY"
	to := "2MwjzD8z9NT72MvbLh7K6Wovk5VgNtSW8p6"
	privateKey := "cNZiKqhwB5gTccW3MfVFVJ5tnLf5RfjLwcCjXprUsrbZnUtcGvZk"
	hash, err := web.Transfer(address, to, privateKey, 1000, 0.001*1e8)
	if err != nil {
		t.Error("转账失败：", err)
	}

	fmt.Println("哈希值：", hash)

}
