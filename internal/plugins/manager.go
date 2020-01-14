package plugins

import (
	"context"
	"fmt"

	go_config "github.com/cheebo/go-config"

	"github.com/nori-io/nori-common/config"

	"github.com/nori-io/nori-common/logger"
	"github.com/nori-io/nori-common/meta"
	"github.com/nori-io/nori-common/plugin"
	"github.com/nori-io/nori-common/version"
	configManager "github.com/nori-io/nori/internal/config"
	"github.com/nori-io/nori/internal/dependency"
	"github.com/nori-io/nori/pkg/errors"
	"github.com/nori-io/nori/pkg/files"
	"github.com/nori-io/nori/pkg/registry"
	"github.com/nori-io/nori/pkg/types"
)

type Manager interface {
	AddFile(path string) (meta.Meta, error)
	AddDir(paths []string) ([]meta.Meta, error)

	Install(ctx context.Context, id meta.ID) error

	Meta(id meta.ID) (meta.Meta, error)

	Start(ctx context.Context, id meta.ID) error
	StartAll(ctx context.Context) error

	Stop(ctx context.Context, id meta.ID) error
	StopAll(ctx context.Context) error

	UnInstall(ctx context.Context, id meta.ID) error
}

type manager struct {
	helpers struct {
		fl files.FilesLoader
	}
	data struct {
		installablePlugins types.MetaList
		runningPlugins     types.MetaList
		files              []*types.File
	}
	configManager config.Manager
	dependency    dependency.Manager
	logger        logger.Logger
	registry      registry.Registry
	version       version.Version
}

func NewManager(cfg go_config.Config, log logger.Logger) Manager {
	v1, _ := version.NewVersion("0.2.0") // todo: fixme
	return &manager{
		helpers: struct {
			fl files.FilesLoader
		}{
			fl: files.NewFilesLoader(),
		},
		data: struct {
			installablePlugins types.MetaList
			runningPlugins     types.MetaList
			files              []*types.File
		}{
			installablePlugins: types.MetaList{},
			runningPlugins:     types.MetaList{},
			files:              []*types.File{},
		},
		logger: log.With(logger.Field{
			Key:   "component",
			Value: "plugin_manager",
		}),
		configManager: configManager.NewManager(cfg),
		dependency:    dependency.NewManager(),
		registry:      registry.NewRegistry(log),
		version:       v1,
	}
}

func (m *manager) AddFile(path string) (meta.Meta, error) {
	f, err := m.helpers.fl.Get(path)
	if err != nil {
		return nil, err
	}

	// plugins files
	m.data.files = append(m.data.files, f)

	cons, err := f.Plugin.Meta().GetCore().GetConstraint()
	if err != nil {
		return nil, err
	}

	if !cons.Check(m.version) {
		return nil, errors.IncompatibleCoreVersion{
			ID:                 f.Plugin.Meta().Id(),
			NeededCoreVersion:  f.Plugin.Meta().GetCore().VersionConstraint,
			CurrentCoreVersion: m.version.Original(),
		}
	}

	// todo
	// check installed or not
	// if plugin not installed then added plugin to list of installable plugins
	// and exit function
	//if _, ok := f.Plugin.(plugin.Installable); ok {
	//	installed, err := m.storage.Plugins().IsInstalled(p.Meta().Id())
	//	if err != nil {
	//		return nil, err
	//	}
	//	if !installed {
	//		m.installablePlugins.Add(p)
	//		return p.Meta(), nil
	//	}
	//}

	// add to dependency manager
	err = m.dependency.Add(f.Plugin.Meta())
	if err != nil {
		return nil, err
	}

	// add to registry, on error remove from dependency manager
	err = m.registry.Add(f.Plugin)
	if err != nil {
		m.dependency.Remove(f.Plugin.Meta().Id())
		return nil, err
	}

	return f.Plugin.Meta(), nil
}

func (m *manager) AddDir(paths []string) ([]meta.Meta, error) {
	files, err := m.helpers.fl.Files(paths)
	if err != nil {
		return nil, err
	}

	mts := []meta.Meta{}

	for _, file := range files {
		// load plugin
		mt, err := m.AddFile(file)
		if err != nil {
			return nil, err
		}

		mts = append(mts, mt)

		m.logger.Info(
			"Found '%s' implements '%s' by '%s'",
			mt.Id().ID,
			mt.GetInterface(),
			mt.GetAuthor().Name,
		)
	}
	return mts, nil
}

func (m *manager) Install(ctx context.Context, id meta.ID) error {
	// @todo check depManager for dependencies
	pm, err := m.data.installablePlugins.Find(id)
	if err != nil {
		return err
	}
	p, err := m.registry.Get(pm.Id())
	installable, ok := p.(plugin.Installable)
	if !ok {
		return errors.NonInstallablePlugin{
			ID: id,
			// todo
			//Path: m.data.files.Find(p.Id()),
		}
	}
	return installable.Install(ctx, m.registry)
}

func (m *manager) Meta(id meta.ID) (meta.Meta, error) {
	p, err := m.registry.Get(id)
	if err != nil {
		return nil, err
	}
	return p.Meta(), nil
}

func (m *manager) Start(ctx context.Context, id meta.ID) error {
	p, err := m.registry.Get(id)
	if err != nil {
		return err
	}

	// all dependencies must be resolvable
	// all dependencies must be started, otherwise start dependency
	var depErrs errors.DependenciesNotFound
	for _, dep := range p.Meta().GetDependencies() {
		did, err := m.dependency.Resolve(dep)
		if err != nil {
			depErrs.Add(p.Meta().Id(), dep)
			continue
		}

		if depMeta, err := m.data.runningPlugins.Find(did); depMeta == nil {
			if err != nil {
				return err
			}
			err = m.Start(ctx, did)
			if err != nil {
				return err
			}
		}
	}

	if depErrs.HasErrors() {
		return depErrs
	}

	err = p.Init(ctx, m.configManager)
	if err != nil {
		return err
	}

	var startErr error
	var recovered interface{}
	func() {
		defer func() {
			recovered = recover()
		}()
		startErr = p.Start(ctx, m.registry)
	}()

	if startErr != nil {
		return startErr
	}

	if recovered != nil {
		// todo: custom error
		return fmt.Errorf("%v", recovered)
	}

	m.data.runningPlugins.Add(p.Meta())

	m.logger.With(LogFieldsMeta(p.Meta())...).Info("Plugin successfully started")

	return nil
}

func (m *manager) StartAll(ctx context.Context) error {
	// start plugins in dependency correct order
	pl, err := m.dependency.Sort()
	if err != nil {
		return err
	}

	for _, id := range pl {
		if err := m.Start(ctx, id); err != nil {
			return err
		}
	}

	return nil
}

func (m *manager) Stop(ctx context.Context, id meta.ID) error {
	if running, _ := m.data.runningPlugins.Find(id); running == nil {
		return nil
	}

	p, err := m.registry.Get(id)
	if err != nil {
		return err
	}

	// stop dependent plugins before stop the plugin
	for _, dep := range m.dependency.GetDependent(id) {
		depPlugin, err := m.registry.Get(dep)
		if err != nil {
			return err
		}
		//@todo collect errors
		if err := m.Stop(ctx, depPlugin.Meta().Id()); err != nil {
			m.logger.With(LogFieldsMeta(depPlugin.Meta())...).Error(err.Error())
		}
	}

	var stopErr error
	var recovered interface{}
	func() {
		defer func() {
			recovered = recover()
		}()
		stopErr = p.Stop(ctx, m.registry)
	}()

	m.data.runningPlugins.Remove(id)

	if recovered != nil {
		return fmt.Errorf("%s", recovered)
	}

	if stopErr != nil {
		return stopErr
	}

	m.logger.With(LogFieldsMeta(p.Meta())...).Info("Plugin successfully stopped")

	return nil
}

func (m *manager) StopAll(ctx context.Context) error {
	// todo stop running plugins in reverse order
	plugins := m.registry.Plugins()
	for i := len(plugins) - 1; i >= 0; i-- {
		p := plugins[i]
		err := m.Stop(ctx, p.Meta().Id())
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *manager) UnInstall(ctx context.Context, id meta.ID) error {
	// @todo check depManager for dependent plugins
	p, err := m.registry.Get(id)
	if err != nil {
		return err
	}

	// stop plugin before uninstall it
	err = p.Stop(ctx, m.registry)
	if err != nil {
		return err
	}

	installable, ok := p.(plugin.Installable)
	if !ok {
		return errors.NonInstallablePlugin{
			ID: p.Meta().Id(),
			// todo: add path
			//Path: m.data.files.Find(p.Meta().Id()),
		}
	}
	err = installable.UnInstall(ctx, m.registry)
	if err != nil {
		m.logger.Error(err.Error())
		return err
	}
	// todo: remove from storage
	//err = m.storage.Plugins().Delete(id)
	if err != nil {
		m.logger.Error(err.Error())
		return err
	}

	m.data.installablePlugins.Add(p.Meta())

	return nil
}
