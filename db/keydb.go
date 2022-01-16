package db

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type KeyDB struct {
	db *leveldb.DB
}

func NewKeyDB(file string) (*KeyDB, error) {
	db, err := leveldb.OpenFile(file, nil)
	if err != nil {
		return nil, err
	}

	return &KeyDB{
		db: db,
	}, nil
}

func (l *KeyDB) Has(key string) (bool, error) {
	has, err := l.db.Has([]byte(key), nil)
	if err != nil {
		return false, err
	}
	return has, nil
}

func (l *KeyDB) Get(key string) (string, error) {
	data, err := l.db.Get([]byte(key), nil)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (l *KeyDB) ListWallet(prefix string) ([]WalletItem, error) {
	iter := l.db.NewIterator(util.BytesPrefix([]byte(prefix)), nil)
	var list []WalletItem
	for iter.Next() {
		list = append(list, WalletItem{
			Address:    string(iter.Key()),
			PrivateKey: string(iter.Value()),
		})
	}
	return list, nil
}

func (l *KeyDB) Put(key string, value string) error {
	err := l.db.Put([]byte(key), []byte(value), nil)
	if err != nil {
		return err
	}
	return nil
}

func (l *KeyDB) Delete(key string) error {
	err := l.db.Delete([]byte(key), nil)
	if err != nil {
		return err
	}
	return nil
}

func (l *KeyDB) Close() error {
	return l.db.Close()
}
