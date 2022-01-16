package db

import "io"

type WalletItem struct {
	Address    string // 地址
	PrivateKey string // 私钥
}

type Reader interface {
	// Has retrieves if a key is present in the key-value data store.
	Has(key string) (bool, error)

	// Get retrieves the given key if it's present in the key-value data store.
	Get(key string) (string, error)

	// ListWallet retrieves the given key if it's present in the key-value data store.
	ListWallet(prefix string) ([]WalletItem, error)
}

type Writer interface {
	// Put inserts the given value into the key-value data store.
	Put(key string, value string) error

	// Delete removes the key from the key-value data store.
	Delete(key string) error
}

type Database interface {
	Reader
	Writer
	io.Closer
}
