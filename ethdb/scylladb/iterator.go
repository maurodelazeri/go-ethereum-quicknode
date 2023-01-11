package scylladb

import (
	"fmt"

	"github.com/ethereum/go-ethereum/ethdb"
)

// iterator can walk over the (potentially partial) keyspace of a memory key
// value store. Internally it is a deep copy of the entire iterated state,
// sorted by keys.
type iterator struct {
	index  int
	keys   []string
	values [][]byte
}

// Next moves the iterator to the next key/value pair. It returns whether the
// iterator is exhausted.
func (it *iterator) Next() bool {
	// Short circuit if iterator is already exhausted in the forward direction.
	if it.index >= len(it.keys) {
		return false
	}
	it.index += 1
	return it.index < len(it.keys)
}

// Error returns any accumulated error. Exhausting all the key/value pairs
// is not considered to be an error. A memory iterator cannot encounter errors.
func (it *iterator) Error() error {
	return nil
}

// Key returns the key of the current key/value pair, or nil if done. The caller
// should not modify the contents of the returned slice, and its contents may
// change on the next call to Next.
func (it *iterator) Key() []byte {
	// Short circuit if iterator is not in a valid position
	if it.index < 0 || it.index >= len(it.keys) {
		return nil
	}
	return []byte(it.keys[it.index])
}

// Value returns the value of the current key/value pair, or nil if done. The
// caller should not modify the contents of the returned slice, and its contents
// may change on the next call to Next.
func (it *iterator) Value() []byte {
	// Short circuit if iterator is not in a valid position
	if it.index < 0 || it.index >= len(it.keys) {
		return nil
	}
	return it.values[it.index]
}

// Release releases associated resources. Release should always succeed and can
// be called multiple times without causing error.
func (it *iterator) Release() {
	it.index, it.keys, it.values = -1, nil, nil
}

// NewIterator creates a binary-alphabetical iterator over a subset
// of database content with a particular key prefix, starting at a particular
// initial key (or after, if it does not exist).
func (db *database) NewIterator(prefix []byte, start []byte) ethdb.Iterator {
	fmt.Println("prefix", string(prefix), "start", string(start))
	// db.lock.RLock()
	// defer db.lock.RUnlock()

	// var (
	// 	pr     = string(prefix)
	// 	st     = string(append(prefix, start...))
	// 	keys   = make([]string, 0, len(db.db))
	// 	values = make([][]byte, 0, len(db.db))
	// )
	// // Collect the keys from the memory database corresponding to the given prefix
	// // and start
	// for key := range db.db {
	// 	if !strings.HasPrefix(key, pr) {
	// 		continue
	// 	}
	// 	if key >= st {
	// 		keys = append(keys, key)
	// 	}
	// }
	// // Sort the items and retrieve the associated values
	// sort.Strings(keys)
	// for _, key := range keys {
	// 	values = append(values, db.db[key])
	// }
	// return &iterator{
	// 	index:  -1,
	// 	keys:   keys,
	// 	values: values,
	// }
	return &iterator{
		index:  -1,
		keys:   []string{},
		values: [][]byte{},
	}
}
