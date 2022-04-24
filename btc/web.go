package btc

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/shopspring/decimal"
)

type Web struct {
	client     *Client
	defaultNet *chaincfg.Params
}

func NewWeb(host, user, pass string, network string) (*Web, error) {

	client, err := NewClient(host, user, pass)

	if err != nil {
		return nil, err
	}

	defaultNet := getNetwork(network)

	return &Web{
		client:     client,
		defaultNet: defaultNet,
	}, nil

}

// CreateWallet 创建地址
func (w *Web) CreateWallet() (*Wallet, error) {
	return CreateWallet(w.defaultNet)
}

// GetBlockCount 获取最高区块数量
func (w *Web) GetBlockCount() (int64, error) {
	return w.client.GetBlockCount()
}

// GetWalletByPrivateKey 根据私钥获取地址
func (w *Web) GetWalletByPrivateKey(privateKeyStr string) (string, error) {
	btcW, err := GetWalletByPrivateKey(w.defaultNet, privateKeyStr)
	if err != nil {
		return "", err
	}
	return btcW.NestedSegWitAddress, nil
}

// GetBalance 获取余额
func (w *Web) GetBalance(address string) (int64, error) {
	fromAddress, err := btcutil.DecodeAddress(address, w.defaultNet)
	if err != nil {
		return 0, err
	}

	balance := decimal.NewFromInt(0)
	listUnSpent, err := w.client.ListUnspent(fromAddress)
	if err != nil {
		fmt.Println("报错")
		return 0, err
	}
	powDecimal := decimal.NewFromInt(1e8)
	for _, result := range listUnSpent {
		if result.Amount == 0 {
			continue
		}
		resultAmount := decimal.NewFromFloat(result.Amount)
		balance.Add(resultAmount.Mul(powDecimal))
	}
	return balance.IntPart(), nil
}

// Transfer 转账
func (w *Web) Transfer(from, to, privateKey string, amount, fee int64) (string, error) {

	fromAddress, err := btcutil.DecodeAddress(from, w.defaultNet)
	if err != nil {
		return "", err
	}
	toAddress, err := btcutil.DecodeAddress(to, w.defaultNet)
	if err != nil {
		return "", err
	}

	// 记录累加金额
	outAmount := decimal.NewFromInt(0)
	listUnSpent, err := w.client.ListUnspent(fromAddress)
	if err != nil {
		return "", err
	}

	amountDecimal := decimal.NewFromInt(amount)
	feeDecimal := decimal.NewFromInt(fee)
	powDecimal := decimal.NewFromInt(1e8)

	// 构造输出
	var outputs []*wire.TxOut
	// 构造输入
	var inputs []*wire.TxIn
	var pkScripts [][]byte // txin 签名用script
	for _, result := range listUnSpent {
		if result.Amount == 0 {
			fmt.Println("sssssss")
			continue
		}
		if outAmount.Cmp(amountDecimal) >= 0 {
			fmt.Println("sssssss")
			// 已经累加到需要转账的金额
			break
		}

		resultAmount := decimal.NewFromFloat(result.Amount)

		outAmount.Add(resultAmount.Mul(powDecimal))
		fmt.Println("余额：", outAmount)

		// 构造输入
		hash, _ := chainhash.NewHashFromStr(result.TxID) // tx hash
		outPoint := wire.NewOutPoint(hash, result.Vout)  // 第几个输出
		txIn := wire.NewTxIn(outPoint, nil, nil)
		inputs = append(inputs, txIn)

		//设置签名用script
		txInPkScript, err := hex.DecodeString(result.ScriptPubKey)
		if err != nil {
			return "", err
		}
		pkScripts = append(pkScripts, txInPkScript)
	}

	fmt.Println(outAmount, amountDecimal)

	// 余额不足
	if outAmount.Cmp(amountDecimal) == -1 {
		return "", errors.New("余额不足")
	}

	// 输出给转账者自己
	leftToMe := outAmount.Sub(amountDecimal).Sub(feeDecimal) // 累加值-转账值-交易费就是剩下再给我的
	pkScript, err := txscript.PayToAddrScript(fromAddress)
	if err != nil {
		return "", err
	}
	outputs = append(outputs, wire.NewTxOut(leftToMe.IntPart(), pkScript))

	// 输出给接收者
	pkScript, err = txscript.PayToAddrScript(toAddress)
	if err != nil {
		return "", err
	}
	outputs = append(outputs, wire.NewTxOut(amount, pkScript))

	tx := &wire.MsgTx{
		Version:  wire.TxVersion,
		TxIn:     inputs,
		TxOut:    outputs,
		LockTime: 0,
	}

	// 签名
	err = sign(tx, privateKey, pkScripts)
	if err != nil {
		return "", err
	}

	// 广播裸交易
	hash, err := w.client.SendRawTransaction(tx)
	if err != nil {
		return "", err
	}

	return hash.String(), nil
}

// 签名
func sign(tx *wire.MsgTx, privKeyStr string, prevPkScripts [][]byte) error {
	inputs := tx.TxIn
	wif, err := btcutil.DecodeWIF(privKeyStr)
	if err != nil {
		return err
	}

	privateKey := wif.PrivKey

	for i := range inputs {
		pkScript := prevPkScripts[i]
		var script []byte
		script, err = txscript.SignatureScript(tx, i, pkScript, txscript.SigHashAll,
			privateKey, false)
		inputs[i].SignatureScript = script
	}
	return nil
}
