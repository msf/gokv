package main

import (
	"fmt"
	"log"

	"github.com/cockroachdb/pebble"
	badger "github.com/dgraph-io/badger/v4"
)

func main() {
	pdb, err := pebble.Open("./tmp/demo", &pebble.Options{
		DisableWAL:                  false,
		MemTableSize:                64 << 20,
		MemTableStopWritesThreshold: 3,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer pdb.Close()

	key := []byte("hello")
	if err := pdb.Set(key, []byte("world"), pebble.NoSync); err != nil {
		log.Fatal(err)
	}
	value, closer, err := pdb.Get(key)
	if err != nil {
		log.Fatal(err)
	}
	closer.Close()
	fmt.Printf("%s %s\n", key, value)

	// Open the Badger database located in the /tmp/badger directory.
	// It will be created if it doesn't exist.
	bdb, err := badger.Open(badger.DefaultOptions("./tmp/badger").
		WithSyncWrites(true).
		WithMemTableSize(64 << 20).
		WithNumMemtables(3),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer bdb.Close()

	bdb.NewWriteBatch()
	// Your code here…
	err = bdb.Update(func(txn *badger.Txn) error {
		err := txn.Set(key, []byte("42"))
		return err
	})
	if err != nil {
		log.Fatal(err)
	}

	err = bdb.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {

			// Accessing val here is valid.
			fmt.Printf("The answer is: %s\n", val)
			return nil
		})
		return err
	})
	if err != nil {
		log.Fatal(err)
	}

}
