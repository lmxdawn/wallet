package db

import (
	"fmt"
	"testing"
)

func TestList(t *testing.T) {

	db, err := NewKeyDB("data/test")
	if err != nil {

	}

	db.Put("fff", "fff")
	db.Put("wallet-aaa", "aaa")
	db.Put("111", "111")

	list, _ := db.ListWallet("wallet-")

	fmt.Println(list)

}
