package quicknode

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/ethdb"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

type database struct {
	db *gorm.DB
}

// iterator can walk over the (potentially partial) keyspace of a memory key
// value store. Internally it is a deep copy of the entire iterated state,
// sorted by keys.
type iterator struct {
	index  int
	keys   []string
	values [][]byte
}

func (db *database) Has(key []byte) (bool, error) {
	// Check if returns RecordNotFound error
	if err := db.db.Where(&EvmData{Key: key}).First(&EvmData{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func (db *database) Get(key []byte) ([]byte, error) {
	out := &EvmData{}
	err := db.db.Where(&EvmData{
		Key: key,
	}).First(out).Error
	if err != nil {
		return nil, err
	}
	return out.Value, nil
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
	err := putOnDatabase(db.db, "evm_data", key, value)
	return err
}

func (db *database) Delete(key []byte) error {
	err := db.db.Where(&EvmData{
		Key: key,
	}).Delete(&EvmData{}).Error
	// Hide not found error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	return err
}

func (db *database) ModifyAncients(f func(ethdb.AncientWriteOp) error) (int64, error) {
	return 0, nil
}

func (db *database) TruncateHead(n uint64) error {
	panic("not supported TruncateHead")
}

func (db *database) TruncateTail(n uint64) error {
	panic("not supported TruncateTail")
}

func (db *database) Sync() error {
	return nil
}

func (db *database) MigrateTable(s string, f func([]byte) ([]byte, error)) error {
	panic("not supported MigrateTable")
}

// NewBatch create a db transaction to batch insert
func (db *database) NewBatch() ethdb.Batch {
	return &batch{
		database:    db,
		transaction: db.db.Begin(),
	}
}

type batch struct {
	*database
	transaction *gorm.DB
	size        int
	finished    bool
}

func (b *batch) Put(key, value []byte) (err error) {
	defer func() {
		// Update size if success, or rollback it
		if err == nil {
			b.size += len(value)
		} else if !b.finished {
			b.finished = true
			b.transaction.Rollback()
		}
	}()
	return putOnDatabase(b.transaction, "evm_data", key, value)
}

func (b *batch) Write() (err error) {
	// This transaction is finished before. There is no data so ignore it.
	if b.finished {
		return nil
	}
	b.finished = true
	return b.transaction.Commit().Error
}

func (b *batch) ValueSize() int {
	return b.size
}

func (b *batch) Reset() {
	// Rollback previous transaction
	if !b.finished {
		b.transaction.Rollback()
	}
	b.transaction = b.db.Begin()
	b.size = 0
	b.finished = false
}

// Replay replays the batch contents.
func (b *batch) Replay(w ethdb.KeyValueWriter) error {
	return nil
}

func (db *database) NewBatchWithSize(size int) ethdb.Batch {
	return &batch{
		database:    db,
		transaction: db.db.Begin(),
	}
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

func (db *database) NewIterator(prefix []byte, start []byte) ethdb.Iterator {
	return &iterator{
		index:  -1,
		keys:   []string{},
		values: [][]byte{},
	}
}

func (db *database) Stat(property string) (string, error) {
	panic("not supported Stat")
}

func (db *database) AncientDatadir() (string, error) {
	panic("not supported AncientDatadir")
}

func (db *database) Compact(start []byte, limit []byte) error {
	return nil
}

func (db *database) NewSnapshot() (ethdb.Snapshot, error) {
	panic("not supported NewSnapshot")
}

func (db *database) Close() error {
	return nil
}

// putOnDatabase replaces the record if exists, or insert a new one
func putOnDatabase(db *gorm.DB, tableName string, key []byte, value []byte) (err error) {
	processedResult := db.Table(tableName).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value"}),
	}).Create(&EvmData{Key: key, Value: value})
	if err := processedResult.Error; err != nil {
		return err
	}
	return nil
}

// NewDatabase returns a MySQL wrapped object.
func NewDatabase() (ethdb.Database, error) {
	// Open db
	// Logger
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,          // Disable color
		},
	)

	// Connects to PostgresDB
	db, err := gorm.Open(
		postgres.Open("host=127.0.0.1 port=5432 user=postgres password=tothemoon342d9dS dbname=eth sslmode=disable"), &gorm.Config{
			Logger: newLogger,
		},
	)
	if err != nil {
		panic(err)
	}

	// Migration
	err = db.AutoMigrate(&EvmData{})

	if err != nil {
		panic(err)
	}

	return &database{
		db: db,
	}, nil
}
