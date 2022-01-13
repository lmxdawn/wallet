package db

import "github.com/syndtr/goleveldb/leveldb"

type KeyDB struct {
	Protocol string // 协议名称
	db       *leveldb.DB
}

func NewKeyDB(protocol string, file string) (*KeyDB, error) {
	db, err := leveldb.OpenFile(file, nil)
	if err != nil {
		return nil, err
	}

	return &KeyDB{
		Protocol: protocol,
		db:       db,
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
