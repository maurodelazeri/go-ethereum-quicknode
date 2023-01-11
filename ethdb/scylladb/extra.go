package scylladb

import (
	"fmt"

	"github.com/ethereum/go-ethereum/ethdb"
)

func (db *Database) HasAncient(kind string, number uint64) (bool, error) {
	if _, err := db.Ancient(kind, number); err != nil {
		return false, nil
	}
	return true, nil
}

func (db *Database) Ancient(kind string, number uint64) ([]byte, error) {
	return nil, nil
}

func (db *Database) AncientRange(kind string, start, count, maxBytes uint64) ([][]byte, error) {
	panic("not supported AncientRange")
}

func (db *Database) Ancients() (uint64, error) {
	return 0, nil
}

func (db *Database) Tail() (uint64, error) {
	return 0, nil
}

func (db *Database) AncientSize(kind string) (uint64, error) {
	return 0, nil
}

func (db *Database) ReadAncients(fn func(op ethdb.AncientReaderOp) error) (err error) {
	return fn(db)
}

func (db *Database) AncientDatadir() (string, error) {
	panic("not supported AncientDatadir")
}

func (db *Database) ModifyAncients(f func(ethdb.AncientWriteOp) error) (int64, error) {
	return 0, nil
}

func (db *Database) TruncateHead(n uint64) error {
	fmt.Println("TruncateHead but not implemented, returning nil")
	return nil
}

func (db *Database) TruncateTail(n uint64) error {
	fmt.Println("TruncateTail but not implemented, returning nil")
	return nil
}

func (db *Database) Sync() error {
	fmt.Println("Sync but not implemented, returning nil")
	return nil
}

func (db *Database) MigrateTable(s string, f func([]byte) ([]byte, error)) error {
	fmt.Println("MigrateTable but not implemented, returning nil")
	return nil
}
