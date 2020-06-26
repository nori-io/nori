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

package nori

import (
	"context"

	log "github.com/nori-io/logger"
	"github.com/nori-io/nori-common/v2/logger"
	"github.com/nori-io/nori-common/v2/storage"
	"github.com/nori-io/nori/internal/domain/entity"
	"github.com/nori-io/nori/internal/domain/manager"
	"github.com/nori-io/nori/internal/env"
)

type Nori struct {
	log      logger.Logger
	env      *env.Env
	managers struct {
		File   manager.File
		Plugin manager.Plugin
	}
	storage storage.Storage
}

func (n *Nori) Run(ctx context.Context) error {
	err := n.loadHooks()
	if err != nil {
		n.log.Error(err.Error())
		return err
	}

	err = n.loadPlugins()
	if err != nil {
		n.log.Error(err.Error())
		return err
	}

	for _, p := range n.managers.Plugin.GetAll() {
		n.log.Info("plugin loaded %s", p.Meta().Id().String())
	}

	err = n.managers.Plugin.StartAll(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (n *Nori) Stop() error {
	log.L().Info("Stop")
	return n.storage.Close()
}

func (n *Nori) loadHooks() error {
	files, err := n.managers.File.Dirs(n.env.Config.Hooks.Hooks)
	if err != nil {
		return err
	}
	// todo: add hooks to logger
	return n.load(files)
}

func (n *Nori) loadPlugins() error {
	files, err := n.managers.File.Dirs(n.env.Config.Plugins.Dirs)
	if err != nil {
		return err
	}
	return n.load(files)
}

func (n *Nori) load(files []*entity.File) error {
	plugins, err := n.managers.File.GetAll(files)
	if err != nil {
		return err
	}

	for _, plugin := range plugins {
		err := n.managers.Plugin.Register(plugin)
		if err != nil {
			return err
		}
	}

	return nil
}
