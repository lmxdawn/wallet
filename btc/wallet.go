package btc

import (
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
)

type Wallet struct {
	LegacyAddress       string // (P2PKH)格式，普通非隔离验证地址（由1开头）| 普及度：较高 | 矿工费：较低
	NestedSegWitAddress string // (P2SH)格式，隔离验证（兼容）地址（由3开头）| 普及度：较高 | 矿工费：较低
	NativeSegWitAddress string // (Bech32)格式，隔离验证（原生）地址（由bc1开头）| 普及度：较低 | 矿工费：最低
	PublicKey           string
	PrivateKey          string
}

// CreateWallet 创建钱包
func CreateWallet(defaultNet *chaincfg.Params) (*Wallet, error) {

	//1.生成私钥，参数：Secp256k1
	privateKey, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		return nil, err
	}

	//2.转成wif格式
	privateKeyWif, err := btcutil.NewWIF(privateKey, defaultNet, true)
	if err != nil {
		return nil, err
	}

	return getWalletByPrivateKey(defaultNet, privateKeyWif)

}

// GetWalletByPrivateKey 获取钱包通过私钥，返回地址和公钥
func GetWalletByPrivateKey(defaultNet *chaincfg.Params, privateKeyStr string) (*Wallet, error) {

	// 转成wif格式
	privateKeyWif, err := btcutil.DecodeWIF(privateKeyStr)
	if err != nil {
		return nil, err
	}

	return getWalletByPrivateKey(defaultNet, privateKeyWif)
}

// getWalletByPrivateKey 根据网络和 私钥的 wif 获取地址
func getWalletByPrivateKey(defaultNet *chaincfg.Params, wif *btcutil.WIF) (*Wallet, error) {

	// 获取publicKey
	publicKeySerial := wif.PrivKey.PubKey().SerializeCompressed()

	publicKey, err := btcutil.NewAddressPubKey(publicKeySerial, defaultNet)
	if err != nil {
		return nil, err
	}

	pkHash := btcutil.Hash160(publicKeySerial)
	nativeSegWitAddressHash, err := btcutil.NewAddressWitnessPubKeyHash(pkHash, defaultNet)
	if err != nil {
		return nil, err
	}

	nestedSegWitAddressWitnessProg, err := txscript.PayToAddrScript(nativeSegWitAddressHash)
	if err != nil {
		return nil, err
	}
	nestedSegWitAddressHash, err := btcutil.NewAddressScriptHash(nestedSegWitAddressWitnessProg, defaultNet)
	if err != nil {
		return nil, err
	}

	return &Wallet{
		LegacyAddress:       publicKey.EncodeAddress(),
		NestedSegWitAddress: nestedSegWitAddressHash.EncodeAddress(),
		NativeSegWitAddress: nativeSegWitAddressHash.EncodeAddress(),
		PublicKey:           publicKey.String(),
		PrivateKey:          wif.String(),
	}, nil
}
