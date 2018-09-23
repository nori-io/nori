package storage

import (
	"errors"

	"github.com/secure2work/nori/core/entities"
)

type none struct {
}

func getNoneStorage() (NoriCoreStorage, error) {
	return none{}, nil
}

func (n none) GetInstallations() ([]entities.PluginMeta, error) {
	return []entities.PluginMeta{}, nil
}

func (n none) SaveInstallation(meta entities.PluginMeta) error {
	return errors.New("Can't save to None storage")
}

func (n none) RemoveInstallation(id string) error {
	return errors.New("Can't delete from None storage")
}
