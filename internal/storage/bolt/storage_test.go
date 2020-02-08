/*
Copyright 2019-2020 The Nori Authors.
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

package bolt_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/nori-io/nori/internal/storage/bolt"
	"github.com/stretchr/testify/assert"
)

func TestNewBoltStorage(t *testing.T) {
	a := assert.New(t)

	tmpFile, err := ioutil.TempFile("", "bbolt_db")
	defer os.Remove(tmpFile.Name()) // clean up

	s, err := bolt.NewBoltStorage(tmpFile.Name(), 0666)
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

func TestNewBoltStorage_Error(t *testing.T) {
	a := assert.New(t)

	s, err := bolt.NewBoltStorage("", 0666)
	a.Error(err)
	a.Nil(s)
}

func TestStorage_Bucket(t *testing.T) {
	a := assert.New(t)

	tmpFile, err := ioutil.TempFile("", "bbolt_db")
	defer os.Remove(tmpFile.Name()) // clean up

	s, err := bolt.NewBoltStorage(tmpFile.Name(), 0666)
	defer s.Close()
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

func TestStorage_BucketError(t *testing.T) {
	a := assert.New(t)

	tmpFile, err := ioutil.TempFile("", "bbolt_db")
	defer os.Remove(tmpFile.Name()) // clean up

	s, err := bolt.NewBoltStorage(tmpFile.Name(), 0666)
	defer s.Close()
	a.NoError(err)

	_, err = s.CreateBucketIfNotExists("")
	a.Error(err)
}

func TestStorage_BucketDelete(t *testing.T) {
	a := assert.New(t)

	tmpFile, err := ioutil.TempFile("", "bbolt_db")
	defer os.Remove(tmpFile.Name()) // clean up

	s, err := bolt.NewBoltStorage(tmpFile.Name(), 0666)
	defer s.Close()
	a.NoError(err)

	const (
		name = "foo"
	)

	b, err := s.CreateBucket(name)
	a.NoError(err)
	a.NotNil(b)

	_, err = s.Bucket(name)
	a.NoError(err)

	s.DeleteBucket(name)
	_, err = s.Bucket(name)
	a.Error(err)
}

func TestStorage_DeletedBucketGet(t *testing.T) {
	a := assert.New(t)

	tmpFile, err := ioutil.TempFile("", "bbolt_db")
	defer os.Remove(tmpFile.Name()) // clean up

	s, err := bolt.NewBoltStorage(tmpFile.Name(), 0666)
	defer s.Close()
	a.NoError(err)

	const (
		name = "foo"
	)

	b, err := s.CreateBucket(name)
	a.NoError(err)
	a.NotNil(b)

	s.DeleteBucket(name)

	_, err = b.Get("test")
	a.Error(err)
}

func TestStorage_DeletedBucketSet(t *testing.T) {
	a := assert.New(t)

	tmpFile, err := ioutil.TempFile("", "bbolt_db")
	defer os.Remove(tmpFile.Name()) // clean up

	s, err := bolt.NewBoltStorage(tmpFile.Name(), 0666)
	defer s.Close()
	a.NoError(err)

	const (
		name = "foo"
	)

	b, err := s.CreateBucket(name)
	a.NoError(err)
	a.NotNil(b)

	s.DeleteBucket(name)

	err = b.Set("test", []byte("valye"))
	a.Error(err)
}

func TestStorage_DeletedBucketDelete(t *testing.T) {
	a := assert.New(t)

	tmpFile, err := ioutil.TempFile("", "bbolt_db")
	defer os.Remove(tmpFile.Name()) // clean up

	s, err := bolt.NewBoltStorage(tmpFile.Name(), 0666)
	defer s.Close()
	a.NoError(err)

	const (
		name = "foo"
	)

	b, err := s.CreateBucket(name)
	a.NoError(err)
	a.NotNil(b)

	s.DeleteBucket(name)

	err = b.Delete("test")
	a.Error(err)
}

func TestStorage_DeletedBucketCursor(t *testing.T) {
	a := assert.New(t)

	tmpFile, err := ioutil.TempFile("", "bbolt_db")
	defer os.Remove(tmpFile.Name()) // clean up

	s, err := bolt.NewBoltStorage(tmpFile.Name(), 0666)
	defer s.Close()
	a.NoError(err)

	const (
		name = "foo"
	)

	b, err := s.CreateBucket(name)
	a.NoError(err)
	a.NotNil(b)

	s.DeleteBucket(name)

	c := b.Cursor()
	a.Nil(c)
}

func TestStorage_ForEach(t *testing.T) {
	a := assert.New(t)

	tmpFile, err := ioutil.TempFile("", "bbolt_db")
	defer os.Remove(tmpFile.Name()) // clean up

	s, err := bolt.NewBoltStorage(tmpFile.Name(), 0666)
	defer s.Close()
	a.NoError(err)

	b, err := s.CreateBucketIfNotExists("items")
	a.NoError(err)

	values := map[string][]byte{
		"c": []byte("three"),
		"a": []byte("one"),
		"b": []byte("two"),
	}

	for v, k := range values {
		a.NoError(b.Set(v, k))
	}

	keys := ""
	vals := []byte{}

	err = b.ForEach(func(k string, v []byte) error {
		keys = keys + k
		vals = append(vals, v...)
		return nil
	})
	a.NoError(err)

	a.Equal("abc", keys)
	a.Equal([]byte("onetwothree"), vals)
}

func TestStorage_ForEachError(t *testing.T) {
	a := assert.New(t)

	tmpFile, err := ioutil.TempFile("", "bbolt_db")
	defer os.Remove(tmpFile.Name()) // clean up

	s, err := bolt.NewBoltStorage(tmpFile.Name(), 0666)
	defer s.Close()
	a.NoError(err)

	b, err := s.CreateBucketIfNotExists("items")
	a.NoError(err)
	s.DeleteBucket("items")

	err = b.ForEach(func(k string, v []byte) error {
		return nil
	})
	a.Error(err)
}

func TestStorage_Cursor(t *testing.T) {
	a := assert.New(t)

	tmpFile, err := ioutil.TempFile("", "bbolt_db")
	defer os.Remove(tmpFile.Name()) // clean up

	s, err := bolt.NewBoltStorage(tmpFile.Name(), 0666)
	defer s.Close()
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
	defer c.Close()
	for k, v := c.First(); k != ""; k, v = c.Next() {
		cursor[k] = v
	}

	a.Equal(values, cursor)
}

func TestStorage_CursorEmpty(t *testing.T) {
	a := assert.New(t)

	tmpFile, err := ioutil.TempFile("", "bbolt_db")
	defer os.Remove(tmpFile.Name()) // clean up

	s, err := bolt.NewBoltStorage(tmpFile.Name(), 0666)
	defer s.Close()
	a.NoError(err)

	b, err := s.CreateBucketIfNotExists("items")
	a.NoError(err)

	c := b.Cursor()
	a.Nil(c)
}

func TestStorage_CursorFirst(t *testing.T) {
	a := assert.New(t)

	tmpFile, err := ioutil.TempFile("", "bbolt_db")
	defer os.Remove(tmpFile.Name()) // clean up

	s, err := bolt.NewBoltStorage(tmpFile.Name(), 0666)
	defer s.Close()
	a.NoError(err)

	b, err := s.CreateBucketIfNotExists("items")
	a.NoError(err)

	values := map[string][]byte{
		"b": []byte("two"),
		"a": []byte("one"),
		"c": []byte("three"),
	}

	for v, k := range values {
		a.NoError(b.Set(v, k))
	}

	c := b.Cursor()
	defer c.Close()

	k, v := c.First()
	a.Equal("a", k)
	a.Equal(values["a"], v)
}

func TestStorage_CursorLast(t *testing.T) {
	a := assert.New(t)

	tmpFile, err := ioutil.TempFile("", "bbolt_db")
	defer os.Remove(tmpFile.Name()) // clean up

	s, err := bolt.NewBoltStorage(tmpFile.Name(), 0666)
	defer s.Close()
	a.NoError(err)

	b, err := s.CreateBucketIfNotExists("items")
	a.NoError(err)

	values := map[string][]byte{
		"c": []byte("three"),
		"a": []byte("one"),
		"b": []byte("two"),
	}

	for v, k := range values {
		a.NoError(b.Set(v, k))
	}

	c := b.Cursor()
	defer c.Close()

	k, v := c.Last()
	a.Equal("c", k)
	a.Equal(values["c"], v)
}

func TestStorage_CursorHasNext(t *testing.T) {
	a := assert.New(t)

	tmpFile, err := ioutil.TempFile("", "bbolt_db")
	defer os.Remove(tmpFile.Name()) // clean up

	s, err := bolt.NewBoltStorage(tmpFile.Name(), 0666)
	defer s.Close()
	a.NoError(err)

	b, err := s.CreateBucketIfNotExists("items")
	a.NoError(err)

	values := map[string][]byte{
		"c": []byte("three"),
		"a": []byte("one"),
		"b": []byte("two"),
	}

	for v, k := range values {
		a.NoError(b.Set(v, k))
	}

	c := b.Cursor()
	defer c.Close()

	c.First()
	a.True(c.HasNext())
	c.Last()
	a.False(c.HasNext())
}

func TestStorage_CursorDelete(t *testing.T) {
	a := assert.New(t)

	tmpFile, err := ioutil.TempFile("", "bbolt_db")
	defer os.Remove(tmpFile.Name()) // clean up

	s, err := bolt.NewBoltStorage(tmpFile.Name(), 0666)
	defer s.Close()
	a.NoError(err)

	b, err := s.CreateBucketIfNotExists("items")
	a.NoError(err)

	values := map[string][]byte{
		"c": []byte("three"),
		"a": []byte("one"),
		"b": []byte("two"),
	}

	for v, k := range values {
		a.NoError(b.Set(v, k))
	}

	c := b.Cursor()
	defer c.Close()

	c.First()
	a.NoError(c.Delete())
	a.NoError(c.Delete())
	a.NoError(c.Delete())
	a.NoError(c.Delete())
}

func TestStorage_CursorNextPrev(t *testing.T) {
	a := assert.New(t)

	tmpFile, err := ioutil.TempFile("", "bbolt_db")
	defer os.Remove(tmpFile.Name()) // clean up

	s, err := bolt.NewBoltStorage(tmpFile.Name(), 0666)
	defer s.Close()
	a.NoError(err)

	b, err := s.CreateBucketIfNotExists("items")
	a.NoError(err)

	values := map[string][]byte{
		"c": []byte("three"),
		"a": []byte("one"),
		"b": []byte("two"),
	}

	for v, k := range values {
		a.NoError(b.Set(v, k))
	}

	c := b.Cursor()
	defer c.Close()

	k1, v1 := c.First()
	k2, v2 := c.Next()
	k3, v3 := c.Prev()

	a.Equal("a", k1)
	a.Equal([]byte("one"), v1)
	a.Equal("b", k2)
	a.Equal([]byte("two"), v2)
	a.Equal("a", k3)
	a.Equal([]byte("one"), v3)
}

func TestStorage_CursorSeek(t *testing.T) {
	a := assert.New(t)

	tmpFile, err := ioutil.TempFile("", "bbolt_db")
	defer os.Remove(tmpFile.Name()) // clean up

	s, err := bolt.NewBoltStorage(tmpFile.Name(), 0666)
	defer s.Close()
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
	defer c.Close()
	k, v := c.Seek("bob")

	a.Equal("bob", k)
	a.Equal([]byte("two"), v)
}
