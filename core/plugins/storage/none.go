package storage

import (
	"errors"

	"github.com/secure2work/nori/core/entities"
)

type none struct {
}

func getNoneStorage() (NoriStorage, error) {
	return none{}, nil
}

func (n none) GetPluginMetas() ([]entities.PluginMeta, error) {
	return []entities.PluginMeta{}, nil
}

func (n none) SavePluginMeta(meta entities.PluginMeta) error {
	return errors.New("Can't save to None storage")
}

func (n none) DeletePluginMeta(id string) error {
	return errors.New("Can't delete from None storage")
}
