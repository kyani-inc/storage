package bolt

import (
	"github.com/boltdb/bolt"
)

type Bolt struct {
	Path []byte
	*bolt.DB
}

func New(path string) (Bolt, error) {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return Bolt{}, err
	}

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(path))
		return err
	})

	return Bolt{[]byte(path), db}, nil
}

func (b Bolt) Put(key string, data []byte) error {
	return b.DB.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(b.Path).Put([]byte(key), data)
	})
}

func (b Bolt) Get(key string) []byte {
	var value []byte

	b.DB.View(func(tx *bolt.Tx) error {
		value = tx.Bucket(b.Path).Get([]byte(key))
		return nil
	})

	return value
}

func (b Bolt) Delete(key string) {
	b.DB.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(b.Path).Delete([]byte(key))
	})
}

func (b Bolt) Flush() {
	b.DB.Update(func(tx *bolt.Tx) error {
		tx.DeleteBucket(b.Path)
		_, err := tx.CreateBucket(b.Path)
		return err
	})
}
