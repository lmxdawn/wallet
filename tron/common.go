package tron

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/fbsobreira/gotron-sdk/pkg/common"
	"math/big"
	"strings"
	"unicode/utf8"
)

func Private2TronAddress(privateKeyStr string) (string, error) {

	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		return "", err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", errors.New("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := address.PubkeyToAddress(*publicKeyECDSA)

	return fromAddress.String(), nil
}

func EthAddress2TronAddress(address string) (string, error) {
	ownerByte, err := common.Hex2Bytes(strings.Replace(address, "0x", "41", 1))
	if err != nil {
		return "", err
	}
	owner := common.EncodeCheck(ownerByte)
	return owner, nil
}

func ToString(data string) (string, error) {
	if common.Has0xPrefix(data) {
		data = data[2:]
	}
	if len(data) > 128 {
		n, _ := ToNumber(data[64:128])
		if n != nil {
			l := n.Uint64()
			if 2*int(l) <= len(data)-128 {
				b, err := hex.DecodeString(data[128 : 128+2*l])
				if err == nil {
					return string(b), nil
				}
			}
		}
	} else if len(data) == 64 {
		// allow string properties as 32 bytes of UTF-8 data
		b, err := hex.DecodeString(data)
		if err == nil {
			i := bytes.Index(b, []byte{0})
			if i > 0 {
				b = b[:i]
			}
			if utf8.Valid(b) {
				return string(b), nil
			}
		}
	}
	return "", fmt.Errorf("cannot parse %s", data)
}

func ToNumber(data string) (*big.Int, error) {
	if common.Has0xPrefix(data) {
		data = data[2:]
	}
	if len(data) == 64 {
		var n big.Int
		_, ok := n.SetString(data, 16)
		if ok {
			return &n, nil
		}
	}
	return nil, fmt.Errorf("cannot parse %s", data)
}
