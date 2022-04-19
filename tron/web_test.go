package tron

import (
	"fmt"
	"github.com/fbsobreira/gotron-sdk/pkg/abi"
	"github.com/fbsobreira/gotron-sdk/pkg/common"
	"github.com/shopspring/decimal"
	"math/big"
	"testing"
)

func TestName(t *testing.T) {

	b := abi.Signature("Transfer(from,to,value)")
	fmt.Println([]byte("transfer(address,uint256)"))
	fmt.Println(common.Bytes2Hex(b))

}

func TestName1(t *testing.T) {

	//s, _ := EthAddress2TronAddress("0x1a5a32bd07c33cd8d9f4bd158f235613480c7eef")
	//s, _ := EthAddress2TronAddress("0x740fbbcf714f0295207adc53f0128c0ff93c16cd")
	//s, _ := EthAddress2TronAddress("0xb9565E907eF7613338A3838a2Cd33D9E71bfFe9A")
	//fmt.Println(s)

	rewardAmount := decimal.NewFromBigInt(big.NewInt(100), 0)
	fmt.Println(rewardAmount)

	ii := big.NewInt(1000000000000000000)
	d := new(big.Int).Mul(big.NewInt(1000489), ii)
	a := big.NewInt(81197043129506)
	b := big.NewInt(1)
	s := new(big.Int).Div(d, a)
	c := new(big.Int).Div(s, b)
	fmt.Println(c)
	fmt.Println(new(big.Int).Div(d, a))
	fmt.Println(1000489 / 123)
}

func TestSlice(t *testing.T) {

	var memberUpdateAddress []string

	memberUpdateAddress = append(memberUpdateAddress, "sss")

	fmt.Println(memberUpdateAddress)

	memberUpdate := make(map[string]interface{})
	fmt.Println(memberUpdate)
	memberUpdate["ss"] = 1

	fmt.Println(memberUpdate)

}
