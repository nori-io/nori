// Copyright © 2018 Secure2Work info@secure2work.com
//
// This program is free software: you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation, either version 3
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package plugins

import (
	"strings"

	"context"

	"github.com/hashicorp/go-version"
	"github.com/secure2work/nori/core/entities"
	"github.com/secure2work/nori/core/plugins/interfaces"
	"github.com/sirupsen/logrus"
)

const (
	constraintSeparator = ":"
)

type PluginEntry interface {
	GetInstance() interface{}
	GetMeta() entities.PluginMeta
	Init(ctx context.Context, cm interfaces.ConfigManager) error
	Install(ctx context.Context, registry PluginRegistry) error
	Start(ctx context.Context, registry PluginRegistry) error
	Stop(ctx context.Context, registry PluginRegistry) error
	UnInstall(ctx context.Context, registry PluginRegistry) error

	isDependent(id string, version string) (bool, error)
	checkDependencies(map[string]PluginEntry) []error
	calcWeight(map[string]PluginEntry) int
	getWeight() int
}

func NewPluginEntry(plugin Plugin, filePath string) PluginEntry {
	return &pluginEntry{
		plugin:   plugin,
		filePath: filePath,
		weight:   -1,
	}
}

type pluginEntry struct {
	filePath string
	log      *logrus.Logger
	plugin   Plugin
	weight   int
}

func (p *pluginEntry) GetInstance() interface{} {
	return p.plugin
}

func (p *pluginEntry) GetMeta() entities.PluginMeta {
	return p.plugin.GetMeta()
}

func (p *pluginEntry) Init(ctx context.Context, cm interfaces.ConfigManager) error {
	return p.plugin.Init(ctx, cm)
}

func (p *pluginEntry) Install(ctx context.Context, registry PluginRegistry) error {
	return p.plugin.Install(ctx, registry)
}

func (p *pluginEntry) Start(ctx context.Context, registry PluginRegistry) error {
	return p.plugin.Start(ctx, registry)
}

func (p *pluginEntry) Stop(ctx context.Context, registry PluginRegistry) error {
	return p.plugin.Start(ctx, registry)
}

func (p *pluginEntry) UnInstall(ctx context.Context, registry PluginRegistry) error {
	return p.plugin.UnInstall(ctx, registry)
}

func (p *pluginEntry) getWeight() int {
	return p.weight
}

func (p *pluginEntry) calcWeight(plugEntries map[string]PluginEntry) int {
	if p.weight > -1 {
		return p.weight
	}

	deps := p.GetMeta().GetDependencies()

	if len(deps) == 0 {
		p.weight = 0
		return 0
	}

	p.weight = len(deps)

	for _, dep := range deps {
		dep = strings.Split(dep, ":")[0]
		pe, ok := plugEntries[dep]
		if !ok {
			continue
		}
		p.weight += pe.calcWeight(plugEntries)
	}

	return p.weight
}

func (p *pluginEntry) checkDependencies(plugEntries map[string]PluginEntry) []error {
	errs := make([]error, 0)

	for _, dep := range p.GetMeta().GetDependencies() {
		var constraint string
		dep, constraint = splitConstraint(dep)

		depPlug, ok := plugEntries[dep]

		if !ok {
			errs = append(errs, &DependencyError{
				PlugName:      p.GetMeta().GetId(),
				PlugVer:       p.GetMeta().GetVersion(),
				DepName:       dep,
				DepConstraint: constraint,
			})
			continue
		}

		ver := depPlug.GetMeta().GetVersion()
		check, err := versionCheck(ver, constraint)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		if !check {
			errs = append(errs, &DependencyError{
				PlugName:      p.GetMeta().GetId(),
				PlugVer:       ver,
				DepName:       dep,
				DepConstraint: constraint,
			})
		}
	}

	return errs
}

func (p *pluginEntry) isDependent(id string, version string) (bool, error) {
	ver := p.GetMeta().GetVersion()
	for _, dep := range p.GetMeta().GetDependencies() {
		var constraint string
		dep, constraint = splitConstraint(dep)
		if dep == id {
			check, err := versionCheck(ver, constraint)
			if err != nil {
				return false, err
			}
			if check {
				return true, nil
			}
		}
	}
	return false, nil
}

func versionCheck(ver, constraint string) (bool, error) {
	if len(constraint) == 0 {
		return true, nil
	}

	v, err := version.NewVersion(ver)
	if err != nil {
		return false, err
	}

	c, err := version.NewConstraint(constraint)
	if err != nil {
		return false, err
	}

	return c.Check(v), nil
}

func splitConstraint(name string) (string, string) {
	ss := strings.Split(name, constraintSeparator)
	if len(ss) == 1 {
		return ss[0], ""
	}
	return ss[0], ss[1]
}
