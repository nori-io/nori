package dep

import (
	"github.com/hashicorp/go-version"
	"github.com/secure2work/nori/core/entities"
)

type Manager interface {
	// add interface
	InterfaceAdd(name string) error

	// add plugin
	PluginAdd(meta entities.PluginMeta) error

	// remove plugin
	PluginRemove(meta entities.PluginMeta)

	// Check dependencies
	CheckDependencies() error
}

type manager struct {
	// id, entry
	plugins map[string][]Entry

	// alias, entry
	aliases map[string]Entry
}

// add alias
func (m *manager) Alias(alias string, meta entities.PluginMeta) error {
	// todo
	return nil
}

// add plugin
func (m *manager) PluginAdd(meta entities.PluginMeta) error {
	// todo
	return nil
}

// remove plugin
func (m *manager) PluginRemove(meta entities.PluginMeta) {
	// todo
}

func (m *manager) CheckDependencies() error {
	// todo
	return nil
}

type Entry interface {
	Id() string
	Version() *version.Version
	Dependencies() []Dependency
}

type entry struct {
	id   string
	ver  *version.Version
	deps []Dependency
}

func (e entry) Id() string {
	return e.id
}

func (e entry) Version() *version.Version {
	return e.ver
}

func (e entry) Dependencies() []Dependency {
	return e.deps
}

type Node interface {
	PID() string
	Version() string
}

type node struct {
	id      int
	pid     string
	version string
}

func (d node) ID() int {
	return d.id
}

func (d node) PID() string {
	return d.pid
}

func (d node) Version() string {
	return d.version
}
