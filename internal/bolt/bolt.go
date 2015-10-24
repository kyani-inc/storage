package bolt

import (
	boltdb "github.com/boltdb/bolt"
)

type Bolt struct {
	Path []byte
	*boltdb.DB
}

func New(path string) (Bolt, error) {
	db, err := boltdb.Open(path, 0600, nil)
	if err != nil {
		return Bolt{}, err
	}

	db.Update(func(tx *boltdb.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(path))
		return err
	})

	return Bolt{[]byte(path), db}, nil
}

func (b Bolt) Put(key string, data []byte) error {
	return b.DB.Update(func(tx *boltdb.Tx) error {
		return tx.Bucket(b.Path).Put([]byte(key), data)
	})
}

func (b Bolt) Get(key string) []byte {
	var value []byte

	b.DB.View(func(tx *boltdb.Tx) error {
		value = tx.Bucket(b.Path).Get([]byte(key))
		return nil
	})

	return value
}

func (b Bolt) Flush() {
	b.DB.Update(func(tx *boltdb.Tx) error {
		tx.DeleteBucket(b.Path)
		_, err := tx.CreateBucket(b.Path)
		return err
	})
}
