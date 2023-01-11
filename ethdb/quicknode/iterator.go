package quicknode

import (
	"fmt"

	"github.com/ethereum/go-ethereum/ethdb"
	"gorm.io/gorm"
)

var _ ethdb.Iterator = &Iterator{}

// Iterator is the type that satisfies the ethdb.Iterator interface Ethereum data using a direct Postgres connection
// Iteratee interface is used in Geth for various tests, trie/sync_bloom.go (for fast sync),
// rawdb.InspectDatabase, and the new core/state/snapshot features.
// This should not be confused with trie.NodeIterator or state.NodeIteraor (which can be constructed
// from the ethdb.KeyValueStoreand ethdb.Database interfaces)
type Iterator struct {
	db                 *gorm.DB
	currentKey, prefix []byte
	err                error
}

// NewIterator returns an ethdb.Iterator interface
func NewIterator(start, prefix []byte, db *gorm.DB) ethdb.Iterator {
	fmt.Println("start", string(start), "prefix", string(prefix))

	return &Iterator{
		db:         db,
		currentKey: start,
		prefix:     prefix,
	}
}

// Next satisfies the ethdb.Iterator interface
// Next moves the iterator to the next key/value pair
// It returns whether the iterator is exhausted
func (i *Iterator) Next() bool {
	return false
}

// Error satisfies the ethdb.Iterator interface
// Error returns any accumulated error
// Exhausting all the key/value pairs is not considered to be an error
func (i *Iterator) Error() error {
	return i.err
}

// Key satisfies the ethdb.Iterator interface
// Key returns the key of the current key/value pair, or nil if done
// The caller should not modify the contents of the returned slice
// and its contents may change on the next call to Next
func (i *Iterator) Key() []byte {
	return i.currentKey
}

// Value satisfies the ethdb.Iterator interface
// Value returns the value of the current key/value pair, or nil if done
// The caller should not modify the contents of the returned slice
// and its contents may change on the next call to Next
func (i *Iterator) Value() []byte {
	out := &EvmData{}
	err := i.db.Where(&EvmData{
		Key: i.currentKey,
	}).First(out).Error
	if err != nil {
		return nil
	}
	return out.Value
}

// Release satisfies the ethdb.Iterator interface
// Release releases associated resources
// Release should always succeed and can be called multiple times without causing error
func (i *Iterator) Release() {
	i.currentKey = nil
	i.prefix = nil
}
