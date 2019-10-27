package storage

import (
	"fmt"

	"github.com/nori-io/nori-common/meta"
)

type UnknownStorageType struct {
	storageType string
}

func (e UnknownStorageType) Error() string {
	return fmt.Sprintf("unknown storage type %s", e.storageType)
}

type UndefinedStorageSource struct {
	path string
}

func (e UndefinedStorageSource) Error() string {
	return fmt.Sprintf("undefined storage source key [%s] in config file", e.path)
}

type NotFound struct {
	id meta.ID
}

func (e NotFound) Error() string {
	return fmt.Sprintf("not found %s", e.id.String())
}
