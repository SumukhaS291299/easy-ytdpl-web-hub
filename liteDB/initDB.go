package litedb

//  Database to save previous data

import (
	"fmt"

	"github.com/charmbracelet/log"

	"github.com/dgraph-io/badger/v4"
)

var (
	Dbcon *badger.DB
	err   error
)

func DB(path string) {
	log.Debug("Opening database for storage")
	Dbcon, err = badger.Open(badger.DefaultOptions(path))
	if err != nil {
		log.Fatal(err)
	}
}

func WriteData(key, value []byte) {
	txn := Dbcon.NewTransaction(true)
	defer txn.Discard()

	// Use the transaction...
	err := txn.Set(key, value)
	if err != nil {
		log.Error("Error writing data", "Write data error: ", err)
	}

	// Commit the transaction and check for error.
	if err := txn.Commit(); err != nil {
		log.Error("Error writing data", "Write data error: ", err)
	}
}

func ViewData(key []byte) []byte {
	var value []byte
	Dbcon.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}

		val, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}

		fmt.Printf("Value: %s\n", val)
		value = val
		return nil
	})
	return value
}
