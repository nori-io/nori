package memory_test

import (
	"testing"

	"github.com/nori-io/nori/internal/storage/memory"
	"github.com/stretchr/testify/assert"
)

func TestNewStorage(t *testing.T) {
	a := assert.New(t)

	s, err := memory.NewStorage()
	a.NoError(err)

	bucket, err := s.CreateBucket("storage")
	a.NoError(err)

	v, err := bucket.Get("empty")
	a.Error(err)
	a.Empty(v)

	err = bucket.Set("foo", []byte("bar"))
	a.NoError(err)

	v, err = bucket.Get("foo")
	a.NoError(err)
	a.Equal([]byte("bar"), v)

	err = bucket.Set("zoo", []byte("bar"))
	a.NoError(err)

	err = bucket.Delete("foo")
	a.NoError(err)

	_, err = bucket.Get("foo")
	a.Error(err)
}

func TestStorage_Bucket(t *testing.T) {
	a := assert.New(t)

	s, err := memory.NewStorage()
	a.NoError(err)

	_, err = s.Bucket("empty")
	a.Error(err)

	b, err := s.CreateBucket("foo")
	a.NoError(err)
	a.NotNil(b)

	b, err = s.CreateBucket("foo")
	a.Error(err)
	a.Nil(b)

	b, err = s.Bucket("foo")
	a.NoError(err)
	a.NotNil(b)

	b, err = s.CreateBucketIfNotExists("foo")
	a.NoError(err)
	a.NotNil(b)

	b, err = s.CreateBucketIfNotExists("bar")
	a.NoError(err)
	a.NotNil(b)

}

func TestStorage_Cursor(t *testing.T) {
	a := assert.New(t)

	s, err := memory.NewStorage()
	a.NoError(err)

	b, err := s.CreateBucketIfNotExists("items")
	a.NoError(err)

	values := map[string][]byte{
		"a": []byte("one"),
		"b": []byte("two"),
	}

	for v, k := range values {
		a.NoError(b.Set(v, k))
	}

	cursor := map[string][]byte{}
	c := b.Cursor()
	for k, v := c.First(); k != ""; k, v = c.Next() {
		cursor[k] = v
	}

	a.Equal(values, cursor)
}

func TestStorage_CursorSeek(t *testing.T) {
	a := assert.New(t)

	s, err := memory.NewStorage()
	a.NoError(err)

	b, err := s.CreateBucketIfNotExists("items")
	a.NoError(err)

	values := map[string][]byte{
		"alice": []byte("one"),
		"bob":   []byte("two"),
		"mary":  []byte("three"),
	}

	for v, k := range values {
		a.NoError(b.Set(v, k))
	}

	c := b.Cursor()
	k, v := c.Seek("bob")

	a.Equal("bob", k)
	a.Equal([]byte("two"), v)
}
