package btc

import (
	"fmt"
	"testing"
)

func TestWeb_CreateWallet(t *testing.T) {
	network := "TestNet3"
	defaultNet := getNetwork(network)
	w, err := CreateWallet(defaultNet)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("1开头地址：", w.LegacyAddress)
	fmt.Println("3开头地址：", w.NestedSegWitAddress)
	fmt.Println("bc1开头地址：", w.NativeSegWitAddress)
	fmt.Println("公钥：", w.PublicKey)
	fmt.Println("私钥：", w.PrivateKey)
}

/**
=== RUN   TestWeb_CreateWallet
112333
*****************************
1开头地址： mxqPCMkp2zBYktDEEDGRogQN9baSFUMhBp
3开头地址： 2N8Ah8xiqrfb37aktkkwYmMqpDGQJMSDAWY
bc1开头地址： tb1qhhmxngfsg3psgj8crq3c6r2mtlfzjma28zd0v5
公钥： 02b5b41245e18d7edb1fc01fa60804b718603de5bcc090aec2e11f49b55a152625
私钥： cNZiKqhwB5gTccW3MfVFVJ5tnLf5RfjLwcCjXprUsrbZnUtcGvZk

=== RUN   TestWeb_CreateWallet
*****************************
1开头地址： mjtDhZ2c5X8bkAgMMyYYmvSkGz9mZSsWi1
3开头地址： 2NGU5d9Y3uKLJXZa6GgYm2bf4wNK33GcAwJ
bc1开头地址： tb1q9lnpjyyff96w6rhf06jeuft00df62uedgwrfk4
公钥： 03d299f3d4096d07aab3b99d33253c5d20cbf83f3a1da1d36585adc12205ee9640
私钥： cTihhVCb3QMk7ZAzgAmnp3nuiaift1soqgvm3FWe6EGmfUnJvEDs
*/

func TestWeb_GetWalletByPrivateKey(t *testing.T) {

	//network := "TestNet3"
	network := "MainNet"
	defaultNet := getNetwork(network)
	privateKey := "cNZiKqhwB5gTccW3MfVFVJ5tnLf5RfjLwcCjXprUsrbZnUtcGvZk"
	w, err := GetWalletByPrivateKey(defaultNet, privateKey)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("1开头地址：", w.LegacyAddress)
	fmt.Println("3开头地址：", w.NestedSegWitAddress)
	fmt.Println("bc1开头地址：", w.NativeSegWitAddress)
	fmt.Println("公钥：", w.PublicKey)
	fmt.Println("私钥：", w.PrivateKey)

}
