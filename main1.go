package main

import (
	"crypto/ecdsa"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
)

const (
	GasLimitErc20 = 76918
	GasPrice      = 500000000000
)

func StringToPrivateKey(privateKeyStr string) (*ecdsa.PrivateKey, error) {
	privateKeyByte, err := hexutil.Decode(privateKeyStr)
	if err != nil {
		return nil, err
	}
	privateKey, err := crypto.ToECDSA(privateKeyByte)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

func SignTransaction(chainID int64, tx *types.Transaction, privateKeyStr string) (string, error) {
	privateKey, err := StringToPrivateKey(privateKeyStr)
	if err != nil {
		return "", err
	}
	signTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(chainID)), privateKey)
	//signTx, err := types.SignTx(tx, types.HomesteadSigner{}, privateKey)
	if err != nil {
		return "", nil
	}

	b, err := rlp.EncodeToBytes(signTx)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func MakeERC20TransferData(toAddress string, amount *big.Int) ([]byte, error) {
	methodId := crypto.Keccak256([]byte("transfer(address,uint256)"))
	var data []byte
	data = append(data, methodId[:4]...)
	paddedAddress := common.LeftPadBytes(common.HexToAddress(toAddress).Bytes(), 32)
	data = append(data, paddedAddress...)
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)
	data = append(data, paddedAmount...)
	return data, nil
}

func OfflineTransferERC20(chainID int64, nonce uint64, toAddress, toContractAddress string, value *big.Int, privk string) (string, error) {
	data, err := MakeERC20TransferData(toAddress, value)
	if err != nil {
		return "", err
	}

	tx := types.NewTransaction(uint64(nonce), common.HexToAddress(toContractAddress), big.NewInt(0), GasLimitErc20, big.NewInt(GasPrice), data)
	return SignTransaction(chainID, tx, privk)
}
