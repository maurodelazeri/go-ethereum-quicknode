package postgres

import (
	"fmt"

	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/jmoiron/sqlx"
)

var _ ethdb.Batch = &Batch{}

// Batch is the type that satisfies the ethdb.Batch interface for PG-IPFS Ethereum data using a direct Postgres connection
type Batch struct {
	db        *sqlx.DB
	tx        *sqlx.Tx
	valueSize int
}

// NewBatch returns a ethdb.Batch interface
func NewBatch(db *sqlx.DB, tx *sqlx.Tx) ethdb.Batch {
	b := &Batch{
		db: db,
		tx: tx,
	}
	if tx == nil {
		b.Reset()
	}
	return b
}

// Put satisfies the ethdb.Batch interface
// Put inserts the given value into the key-value data store
// Key is expected to be the keccak256 hash of value
func (b *Batch) Put(key []byte, value []byte) (err error) {
	if _, err = b.tx.Exec(putPgStr, key, value); err != nil {
		return err
	}
	b.valueSize += len(value)
	return nil
}

// Delete satisfies the ethdb.Batch interface
// Delete removes the key from the key-value data store
func (b *Batch) Delete(key []byte) (err error) {
	mhKey, err := MultihashKeyFromKeccak256(key)
	if err != nil {
		return err
	}
	_, err = b.tx.Exec(deletePgStr, mhKey)
	return err
}

// ValueSize satisfies the ethdb.Batch interface
// ValueSize retrieves the amount of data queued up for writing
// The returned value is the total byte length of all data queued to write
func (b *Batch) ValueSize() int {
	return b.valueSize
}

// Write satisfies the ethdb.Batch interface
// Write flushes any accumulated data to disk
func (b *Batch) Write() error {
	if b.tx == nil {
		return nil
	}
	return b.tx.Commit()
}

// Replay satisfies the ethdb.Batch interface
// Replay replays the batch contents
func (b *Batch) Replay(w ethdb.KeyValueWriter) error {
	fmt.Println("Replay nil")
	return nil
}

// Reset satisfies the ethdb.Batch interface
// Reset resets the batch for reuse
// This should be called after every write
func (b *Batch) Reset() {
	var err error
	b.tx, err = b.db.Beginx()
	if err != nil {
		panic(err)
	}
	b.valueSize = 0
}
