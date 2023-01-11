package scylladb

import (
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/gocql/gocql"
)

type database struct {
	session *gocql.Session
	lock    sync.RWMutex
}

// NewIterator satisfies the ethdb.Iteratee interface
// it creates a binary-alphabetical iterator over a subset
// of database content with a particular key prefix, starting at a particular
// initial key (or after, if it does not exist).
//
// Note: This method assumes that the prefix is NOT part of the start, so there's
// no need for the caller to prepend the prefix to the start
func (d *database) NewIterator(prefix []byte, start []byte) ethdb.Iterator {
	return NewIterator(start, prefix, d.session)
}

func (db *database) Has(key []byte) (bool, error) {
	var value []byte
	if err := db.session.Query(`SELECT value FROM blockchain WHERE key = ?`, key).Consistency(gocql.One).Scan(&value); err != nil {
		return false, err
	}
	if len(value) > 0 {
		return true, nil
	}
	return false, nil
}

func (db *database) Get(key []byte) ([]byte, error) {
	var value []byte
	if err := db.session.Query(`SELECT value FROM blockchain WHERE key = ?`, key).Consistency(gocql.One).Scan(&value); err != nil {
		return nil, err
	}
	if len(value) > 0 {
		return value, nil
	}
	return nil, nil
}

func (db *database) HasAncient(kind string, number uint64) (bool, error) {
	if _, err := db.Ancient(kind, number); err != nil {
		return false, nil
	}
	return true, nil
}

func (db *database) Ancient(kind string, number uint64) ([]byte, error) {
	return nil, nil
}

func (db *database) AncientRange(kind string, start, count, maxBytes uint64) ([][]byte, error) {
	panic("not supported AncientRange")
}

func (db *database) Ancients() (uint64, error) {
	return 0, nil
}

func (db *database) Tail() (uint64, error) {
	return 0, nil
}

func (db *database) AncientSize(kind string) (uint64, error) {
	return 0, nil
}

func (db *database) ReadAncients(fn func(op ethdb.AncientReaderOp) error) (err error) {
	return fn(db)
}

func (db *database) Put(key []byte, value []byte) error {
	if err := db.session.Query(`INSERT INTO blockchain (key,value) VALUES (?, ?)`, key, value).Exec(); err != nil {
		return err
	}
	return nil
}

func (db *database) Delete(key []byte) error {
	if err := db.session.Query(`DELETE from FROM blockchain WHERE key = ?`, key).Exec(); err != nil {
		return err
	}
	return nil
}

func (db *database) ModifyAncients(f func(ethdb.AncientWriteOp) error) (int64, error) {
	return 0, nil
}

func (db *database) TruncateHead(n uint64) error {
	fmt.Println("TruncateHead but not implemented, returning nil")
	return nil
}

func (db *database) TruncateTail(n uint64) error {
	fmt.Println("TruncateTail but not implemented, returning nil")
	return nil
}

func (db *database) Sync() error {
	fmt.Println("Sync but not implemented, returning nil")
	return nil
}

func (db *database) MigrateTable(s string, f func([]byte) ([]byte, error)) error {
	fmt.Println("MigrateTable but not implemented, returning nil")
	return nil
}

// NewBatch creates a write-only key-value store that buffers changes to its host
// database until a final write is called.
func (db *database) NewBatch() ethdb.Batch {
	return &batch{
		db: db,
	}
}

// NewBatchWithSize creates a write-only database batch with pre-allocated buffer.
func (db *database) NewBatchWithSize(size int) ethdb.Batch {
	return &batch{
		db: db,
	}
}

func (db *database) Stat(property string) (string, error) {
	fmt.Println("Stat but not implemented, returning nil", property)
	return "quicknode", nil
}

func (db *database) AncientDatadir() (string, error) {
	panic("not supported AncientDatadir")
}

func (db *database) Compact(start []byte, limit []byte) error {
	return nil
}

func (db *database) NewSnapshot() (ethdb.Snapshot, error) {
	fmt.Println("NewSnapshot but not implemented, returning nil")
	return nil, nil
}

func (db *database) Close() error {
	return nil
}

// NewDatabase returns a MySQL wrapped object.
func NewDatabase() (ethdb.Database, error) {
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "eth"
	cluster.Consistency = gocql.Any
	session, _ := cluster.CreateSession()

	return &database{
		session: session,
	}, nil
}
