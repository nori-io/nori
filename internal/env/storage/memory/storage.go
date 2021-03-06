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

package memory

import (
	"sort"

	"github.com/nori-io/common/v5/pkg/domain/storage"
	"github.com/nori-io/common/v5/pkg/errors"
)

func New() (storage.Storage, error) {
	return &memStorage{
		buckets: map[string]storage.Bucket{},
	}, nil
}

type (
	memStorage struct {
		buckets map[string]storage.Bucket
	}

	bucket struct {
		storage *memStorage
		name    string
		kv      map[string][]byte
	}

	cursor struct {
		index  uint
		items  []string
		bucket *bucket
	}
)

// CreateBucket creates a new bucket with the given name and returns it.
func (s *memStorage) CreateBucket(name string) (storage.Bucket, error) {
	if _, ok := s.buckets[name]; ok {
		return nil, errors.EntityAlreadyExists{Entity: name}
	}
	b := &bucket{
		storage: s,
		name:    name,
		kv:      map[string][]byte{},
	}
	s.buckets[name] = b
	return b, nil
}

// CreateBucketIfNotExists creates a new bucket if it doesn't already exist and returns it.
func (s *memStorage) CreateBucketIfNotExists(name string) (storage.Bucket, error) {
	if b, ok := s.buckets[name]; ok {
		return b, nil
	}

	b := &bucket{
		storage: s,
		name:    name,
		kv:      map[string][]byte{},
	}
	s.buckets[name] = b
	return b, nil
}

// Bucket returns a bucket by name.
func (s *memStorage) Bucket(name string) (storage.Bucket, error) {
	if b, ok := s.buckets[name]; ok {
		return b, nil
	}
	return nil, errors.EntityNotFound{Entity: name}
}

// DeleteBucket deletes a bucket with the given name.
func (s *memStorage) DeleteBucket(name string) error {
	delete(s.buckets, name)
	return nil
}

func (s *memStorage) Close() error {
	s.buckets = map[string]storage.Bucket{}
	return nil
}

// Get retrieves the value for a key.
func (b *bucket) Get(key string) ([]byte, error) {
	if _, ok := b.storage.buckets[b.name]; !ok {
		return nil, errors.EntityNotFound{Entity: b.name}
	}
	if v, ok := b.kv[key]; ok {
		return v, nil
	}
	return nil, errors.EntityNotFound{Entity: key}
}

// Set sets the value for a key.
func (b *bucket) Set(key string, value []byte) error {
	if _, ok := b.storage.buckets[b.name]; !ok {
		return errors.EntityNotFound{Entity: b.name}
	}
	b.kv[key] = value
	return nil
}

// Delete removes a key
func (b *bucket) Delete(key string) error {
	if _, ok := b.storage.buckets[b.name]; !ok {
		return errors.EntityNotFound{Entity: b.name}
	}
	delete(b.kv, key)
	return nil
}

// ForEach executes a function for each key/value pair
func (b *bucket) ForEach(fn func(k string, v []byte) error) error {
	if _, ok := b.storage.buckets[b.name]; !ok {
		return errors.EntityNotFound{Entity: b.name}
	}
	c := b.Cursor()
	if c == nil {
		return nil
	}

	for key, val := c.First(); key != ""; key, val = c.Next() {
		if err := fn(key, val); err != nil {
			return err
		}
	}
	return nil
}

func (b *bucket) Cursor() storage.Cursor {
	if _, ok := b.storage.buckets[b.name]; !ok {
		return nil
	}

	if len(b.kv) == 0 {
		return nil
	}

	items := make([]string, len(b.kv))
	i := 0
	for k, _ := range b.kv {
		items[i] = k
		i++
	}
	sort.Strings(items)
	return &cursor{items: items, bucket: b}
}

func (c *cursor) First() (key string, value []byte) {
	c.index = 0
	if len(c.items) == 0 {
		return "", nil
	}
	return c.items[c.index], c.bucket.kv[c.items[c.index]]
}

func (c *cursor) Last() (key string, value []byte) {
	c.index = uint(len(c.items) - 1)
	if len(c.items) == 0 {
		return "", nil
	}
	return c.items[c.index], c.bucket.kv[c.items[c.index]]
}

func (c *cursor) Next() (key string, value []byte) {
	if c.index+1 >= uint(len(c.items)) {
		return "", nil
	}
	c.index++
	return c.items[c.index], c.bucket.kv[c.items[c.index]]
}

func (c *cursor) Prev() (key string, value []byte) {
	if c.index == 0 {
		return "", nil
	}
	c.index--
	return c.items[c.index], c.bucket.kv[c.items[c.index]]
}

func (c *cursor) Close() error {
	c.items = []string{}
	return nil
}

// Seek moves the cursor to a seek key and returns it,
// If the key does not exist then then next key is used.
// If there are no keys, an empty key is returned
func (c *cursor) Seek(seek string) (key string, value []byte) {
	if len(c.items) == 0 {
		return "", nil
	}
	var idx int = -1
	for i, v := range c.items {
		if v == seek {
			idx = i
			break
		}
	}

	if idx >= 0 {
		c.index = uint(idx)
	} else {
		if c.index+1 < uint(len(c.items)) {
			c.index++
		}
	}

	return c.items[c.index], c.bucket.kv[c.items[c.index]]
}

// Delete removes current key-value
func (c *cursor) Delete() error {
	if len(c.items) == 0 {
		return nil
	}
	delete(c.bucket.kv, c.items[c.index])
	c.items = append(c.items[:c.index], c.items[c.index+1:]...)
	return nil
}

// HasNext returns true if next element exists
func (c *cursor) HasNext() bool {
	return (c.index + 1) < uint(len(c.items))
}
