package storage

import (
	"errors"

	"github.com/secure2work/nori/core/plugins/meta"
)

type none struct {
}

func getNoneStorage() (NoriStorage, error) {
	return none{}, nil
}

func (n none) GetPluginMetas() ([]meta.Meta, error) {
	return []meta.Meta{}, nil
}

func (n none) SavePluginMeta(meta meta.Meta) error {
	return errors.New("Can't save to None storage")
}

func (n none) DeletePluginMeta(id meta.ID) error {
	return errors.New("Can't delete from None storage")
}
