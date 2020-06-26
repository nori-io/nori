package storage

import (
	"github.com/nori-io/nori-common/v2/storage"
	"github.com/nori-io/nori/internal/domain/repository"
	"github.com/nori-io/nori/internal/env"
)

type Service struct {
	FileRepository repository.FileRepository
	Env            *env.Env
	Storage        storage.Storage
}

// CreateBucket creates a new bucket with the given name and returns it.
func (s *Service) CreateBucket(name string) (storage.Bucket, error) {
	return s.Storage.CreateBucket(name)
}

// CreateBucketIfNotExists creates a new bucket if it doesn't already exist and returns it.
func (s *Service) CreateBucketIfNotExists(name string) (storage.Bucket, error) {
	return s.Storage.CreateBucketIfNotExists(name)
}

// Bucket returns a bucket by name.
func (s *Service) Bucket(name string) (storage.Bucket, error) {
	return s.Storage.Bucket(name)
}

// DeleteBucket deletes a bucket with the given name.
func (s *Service) DeleteBucket(name string) error {
	return s.Storage.DeleteBucket(name)
}

// Close current storage: connection to db, file etc
func (s *Service) Close() error {
	return s.Storage.Close()
}
