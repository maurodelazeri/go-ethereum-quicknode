package scylladb

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
)

// keyvalue is a key-value tuple tagged with a deletion field to allow creating
// memory-database write batches.
type keyvalue struct {
	key    []byte
	value  []byte
	delete bool
}

// batch is a write-only memory batch that commits changes to its host
// database when Write is called. A batch cannot be used concurrently.
type batch struct {
	db     *database
	writes []keyvalue
	size   int
}

// Put inserts the given value into the batch for later committing.
func (b *batch) Put(key, value []byte) error {
	b.writes = append(b.writes, keyvalue{common.CopyBytes(key), common.CopyBytes(value), false})
	b.size += len(key) + len(value)
	return nil
}

// Delete inserts the a key removal into the batch for later committing.
func (b *batch) Delete(key []byte) error {
	b.writes = append(b.writes, keyvalue{common.CopyBytes(key), nil, true})
	b.size += len(key)
	return nil
}

// ValueSize retrieves the amount of data queued up for writing.
func (b *batch) ValueSize() int {
	return b.size
}

// Write flushes any accumulated data to the memory database.
func (b *batch) Write() error {
	for _, keyvalue := range b.writes {
		if keyvalue.delete {
			if err := b.db.session.Query(`DELETE from FROM blockchain WHERE key = ?`, keyvalue.key).Exec(); err != nil {
				return err
			}
			continue
		}
		if err := b.db.session.Query(`INSERT INTO blockchain (key,value) VALUES (?, ?)`, keyvalue.key, keyvalue.value).Exec(); err != nil {
			return err
		}
	}
	return nil
}

// Reset resets the batch for reuse.
func (b *batch) Reset() {
	b.writes = b.writes[:0]
	b.size = 0
}

// Replay replays the batch contents.
func (b *batch) Replay(w ethdb.KeyValueWriter) error {
	for _, keyvalue := range b.writes {
		if keyvalue.delete {
			if err := b.db.session.Query(`DELETE from FROM blockchain WHERE key = ?`, keyvalue.key).Exec(); err != nil {
				return err
			}
			continue
		}
		if err := b.db.session.Query(`INSERT INTO blockchain (key,value) VALUES (?, ?)`, keyvalue.key, keyvalue.value).Exec(); err != nil {
			return err
		}
	}
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
