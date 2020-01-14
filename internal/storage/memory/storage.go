package memory

import (
	"github.com/nori-io/nori-common/storage"
)

func NewStorage() (storage.Storage, error) {
	return &memStorage{
		buckets: map[string]storage.Bucket{},
	}, nil
}

type (
	memStorage struct {
		buckets map[string]storage.Bucket
	}

	bucket struct {
		kv map[string][]byte
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
		return nil, storage.AlreadyExists{Entity: name}
	}
	b := &bucket{
		kv: map[string][]byte{},
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
		kv: map[string][]byte{},
	}
	s.buckets[name] = b
	return b, nil
}

// Bucket returns a bucket by name.
func (s *memStorage) Bucket(name string) (storage.Bucket, error) {
	if b, ok := s.buckets[name]; ok {
		return b, nil
	}
	return nil, storage.NotFound{Entity: name}
}

// DeleteBucket deletes a bucket with the given name.
func (s *memStorage) DeleteBucket(name string) error {
	delete(s.buckets, name)
	return nil
}

// Get retrieves the value for a key.
func (b *bucket) Get(key string) ([]byte, error) {
	if v, ok := b.kv[key]; ok {
		return v, nil
	}
	return nil, storage.NotFound{Entity: key}
}

// Set sets the value for a key.
func (b *bucket) Set(key string, value []byte) error {
	b.kv[key] = value
	return nil
}

// Delete removes a key
func (b *bucket) Delete(key string) error {
	delete(b.kv, key)
	return nil
}

// ForEach executes a function for each key/value pair
func (b *bucket) ForEach(fn func(k string, v []byte) error) error {
	for k, v := range b.kv {
		if err := fn(k, v); err != nil {
			return err
		}
	}
	return nil
}

func (b *bucket) Cursor() storage.Cursor {
	if len(b.kv) == 0 {
		return nil
	}

	items := make([]string, len(b.kv))
	i := 0
	for k, _ := range b.kv {
		items[i] = k
		i++
	}
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
		return storage.NotFound{Entity: ""}
	}
	delete(c.bucket.kv, c.items[c.index])
	c.items = append(c.items[:c.index], c.items[c.index+1:]...)
	return nil
}

// HasNext returns true if next element exists
func (c *cursor) HasNext() bool {
	return (c.index + 1) < uint(len(c.items))
}
