package tron

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/core"
	"github.com/golang/protobuf/proto"
)

func SignTransaction(transaction *core.Transaction, privateKey string) (*core.Transaction, error) {
	privateBytes, err := hex.DecodeString(privateKey)
	if err != nil {
		return nil, fmt.Errorf("hex decode private key error: %v", err)
	}
	priv := crypto.ToECDSAUnsafe(privateBytes)
	defer zeroKey(priv)
	rawData, err := proto.Marshal(transaction.GetRawData())
	if err != nil {
		return nil, fmt.Errorf("proto marshal tx raw data error: %v", err)
	}
	h256h := sha256.New()
	h256h.Write(rawData)
	hash := h256h.Sum(nil)
	signature, err := crypto.Sign(hash, priv)
	if err != nil {
		return nil, fmt.Errorf("sign error: %v", err)
	}
	transaction.Signature = append(transaction.Signature, signature)
	return transaction, nil
}

// zeroKey zeroes a private key in memory.
func zeroKey(k *ecdsa.PrivateKey) {
	b := k.D.Bits()
	for i := range b {
		b[i] = 0
	}
}
