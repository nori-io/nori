/*
Copyright 2018-2020 The Nori Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package bolt

import (
	"fmt"
	"os"

	"github.com/nori-io/common/v5/pkg/domain/storage"
	err2 "github.com/nori-io/common/v5/pkg/errors"
	bolt "go.etcd.io/bbolt"
)

func New(file string, mode os.FileMode) (storage.Storage, error) {
	db, err := bolt.Open(file, mode, nil)
	if err != nil {
		return nil, err
	}
	return &boltStorage{
		db: db,
	}, nil
}

type (
	boltStorage struct {
		db *bolt.DB
	}

	bucket struct {
		db   *bolt.DB
		name []byte
	}

	cursor struct {
		c  *bolt.Cursor
		tx *bolt.Tx
	}
)

// CreateBucket creates a new bucket with the given name and returns it.
func (s *boltStorage) CreateBucket(name string) (storage.Bucket, error) {
	err := s.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte(name))
		if err != nil {
			return fmt.Errorf("can not create %s bucket: %s", name, err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &bucket{
		db:   s.db,
		name: []byte(name),
	}, nil
}

// CreateBucketIfNotExists creates a new bucket if it doesn't already exist and returns it.
func (s *boltStorage) CreateBucketIfNotExists(name string) (storage.Bucket, error) {
	err := s.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(name))
		if err != nil {
			return fmt.Errorf("can not create %s bucket: %s", name, err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &bucket{
		db:   s.db,
		name: []byte(name),
	}, nil
}

// Bucket returns a bucket by name.
func (s *boltStorage) Bucket(name string) (storage.Bucket, error) {
	if err := s.db.Update(func(tx *bolt.Tx) error {
		if tx.Bucket([]byte(name)) == nil {
			return fmt.Errorf("bucket %s does not exist", name)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return &bucket{
		db:   s.db,
		name: []byte(name),
	}, nil
}

// DeleteBucket deletes a bucket with the given name.
func (s *boltStorage) DeleteBucket(name string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte(name))
	})
}

func (s *boltStorage) Close() error {
	return s.db.Close()
}

// Get retrieves the value for a key.
func (b *bucket) Get(key string) ([]byte, error) {
	var v []byte
	err := b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(b.name)
		if bucket == nil {
			return fmt.Errorf("bucket %s not found", b.name)
		}
		v = bucket.Get([]byte(key))
		return nil
	})
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, err2.EntityNotFound{Entity: key}
	}
	return v, err
}

// Set sets the value for a key.
func (b *bucket) Set(key string, value []byte) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(b.name)
		if bucket == nil {
			return fmt.Errorf("bucket %s not found", b.name)
		}
		return bucket.Put([]byte(key), value)
	})
}

// Delete removes a key
func (b *bucket) Delete(key string) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(b.name)
		if bucket == nil {
			return fmt.Errorf("bucket %s not found", b.name)
		}
		return bucket.Delete([]byte(key))
	})
}

// ForEach executes a function for each key/value pair
func (b *bucket) ForEach(fn func(k string, v []byte) error) error {
	return b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(b.name)
		if bucket == nil {
			return fmt.Errorf("bucket %s not found", b.name)
		}
		bucket.ForEach(func(k, v []byte) error {
			return fn(string(k), v)
		})
		return nil
	})
}

func (b *bucket) Cursor() storage.Cursor {
	tx, err := b.db.Begin(true)
	if err != nil {
		tx.Rollback()
		return nil
	}

	bucket := tx.Bucket(b.name)
	if bucket == nil {
		tx.Rollback()
		return nil
	}
	cur := bucket.Cursor()

	// empty Cursor or not
	k, _ := cur.First()
	if k == nil {
		tx.Rollback()
		return nil
	}

	return &cursor{
		c:  cur,
		tx: tx,
	}
}

func (c *cursor) First() (key string, value []byte) {
	k, v := c.c.First()
	return string(k), v
}

func (c *cursor) Last() (key string, value []byte) {
	k, v := c.c.Last()
	return string(k), v
}

func (c *cursor) Next() (key string, value []byte) {
	k, v := c.c.Next()
	return string(k), v
}

func (c *cursor) Prev() (key string, value []byte) {
	k, v := c.c.Prev()
	return string(k), v
}

func (c *cursor) Close() error {
	return c.tx.Rollback()
}

// Seek moves the cursor to a seek key and returns it,
// If the key does not exist then then next key is used.
// If there are no keys, an empty key is returned
func (c *cursor) Seek(seek string) (key string, value []byte) {
	k, v := c.c.Seek([]byte(seek))
	return string(k), v
}

// Delete removes current key-value
func (c *cursor) Delete() error {
	return c.c.Delete()
}

// HasNext returns true if next element exists
func (c *cursor) HasNext() bool {
	// todo: optimise?
	k, _ := c.c.Next()
	if k == nil {
		return false
	}
	c.c.Prev()
	return true
}
