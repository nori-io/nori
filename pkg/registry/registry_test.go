package registry_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/nori-io/logger"
	"github.com/nori-io/nori-common/v2/meta"
	"github.com/nori-io/nori-common/v2/mocks/mock_logger"
	"github.com/nori-io/nori-common/v2/mocks/mock_plugin"
	"github.com/nori-io/nori/pkg/registry"
	"github.com/stretchr/testify/assert"
)

func TestRegistry_ID(t *testing.T) {
	a := assert.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	l := mock_logger.NewMockLogger(ctrl)
	l.EXPECT().With(gomock.Any()).Return(l).AnyTimes()
	r := registry.NewRegistry(logger.New())

	p_ID := meta.ID{ID: "nori/Http", Version: "1.0.0"}
	p_http := mock_plugin.NewMockPlugin(ctrl)
	p_http.EXPECT().Meta().Return(&meta.Data{
		ID:           p_ID,
		Author:       meta.Author{},
		Dependencies: []meta.Dependency{},
		Description:  meta.Description{},
		Core:         meta.Core{},
		Interface:    "nori/Http@1.0.0",
		License:      nil,
		Links:        nil,
		Repository:   meta.Repository{},
		Tags:         nil,
	}).AnyTimes()

	err := r.Add(p_http)
	a.NoError(err)

	p, err := r.Get(p_ID)
	a.NoError(err)
	a.NotNil(p)
}

func TestRegistry_Interface(t *testing.T) {
	a := assert.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	l := mock_logger.NewMockLogger(ctrl)
	l.EXPECT().With(gomock.Any()).Return(l).AnyTimes()
	r := registry.NewRegistry(logger.New())

	p_ID := meta.ID{ID: "nori/Http", Version: "1.0.0"}
	Interface := meta.Interface("nori/Http@1.0.0")
	p_http := mock_plugin.NewMockPlugin(ctrl)
	p_http.EXPECT().Meta().Return(&meta.Data{
		ID:           p_ID,
		Author:       meta.Author{},
		Dependencies: []meta.Dependency{},
		Description:  meta.Description{},
		Core:         meta.Core{},
		Interface:    Interface,
		License:      nil,
		Links:        nil,
		Repository:   meta.Repository{},
		Tags:         nil,
	}).AnyTimes()
	p_http.EXPECT().Instance().Return(&http{}).AnyTimes()

	err := r.Add(p_http)
	a.NoError(err)

	p, err := r.Interface(Interface)
	a.NoError(err)
	a.NotNil(p)
}

func TestRegistry_ResolveFromOne(t *testing.T) {
	a := assert.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	l := mock_logger.NewMockLogger(ctrl)
	l.EXPECT().With(gomock.Any()).Return(l).AnyTimes()
	r := registry.NewRegistry(logger.New())

	dep := meta.Dependency{
		Constraint: "^1.0.0",
		Interface:  "nori/Http",
	}
	p_ID := meta.ID{ID: "nori/Http", Version: "1.0.0"}
	p_http := mock_plugin.NewMockPlugin(ctrl)
	p_http.EXPECT().Meta().Return(&meta.Data{
		ID:           p_ID,
		Author:       meta.Author{},
		Dependencies: []meta.Dependency{},
		Description:  meta.Description{},
		Core:         meta.Core{},
		Interface:    "nori/Http@1.0.0",
		License:      nil,
		Links:        nil,
		Repository:   meta.Repository{},
		Tags:         nil,
	}).AnyTimes()
	p_http.EXPECT().Instance().Return(&http{}).AnyTimes()

	err := r.Add(p_http)
	a.NoError(err)

	p, err := r.Resolve(dep)
	a.NoError(err)
	a.NotNil(p)
}

func TestRegistry_ResolveFromList(t *testing.T) {
	a := assert.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	l := mock_logger.NewMockLogger(ctrl)
	l.EXPECT().With(gomock.Any()).Return(l).AnyTimes()
	r := registry.NewRegistry(logger.New())

	dep := meta.Dependency{
		Constraint: "^1.1.2",
		Interface:  "nori/Http",
	}

	p_http_v1_0 := mock_plugin.NewMockPlugin(ctrl)
	p_http_v1_0.EXPECT().Meta().Return(&meta.Data{
		ID:           meta.ID{ID: "nori/Http", Version: "1.0.0"},
		Author:       meta.Author{},
		Dependencies: []meta.Dependency{},
		Description:  meta.Description{},
		Core:         meta.Core{},
		Interface:    "nori/Http@1.0.0",
		License:      nil,
		Links:        nil,
		Repository:   meta.Repository{},
		Tags:         nil,
	}).AnyTimes()
	p_http_v1_0.EXPECT().Instance().Return(&http{v: "1.0.0"}).AnyTimes()

	p_http_v1_1 := mock_plugin.NewMockPlugin(ctrl)
	p_http_v1_1.EXPECT().Meta().Return(&meta.Data{
		ID:           meta.ID{ID: "nori/Http", Version: "1.1.0"},
		Author:       meta.Author{},
		Dependencies: []meta.Dependency{},
		Description:  meta.Description{},
		Core:         meta.Core{},
		Interface:    "nori/Http@1.1.0",
		License:      nil,
		Links:        nil,
		Repository:   meta.Repository{},
		Tags:         nil,
	}).AnyTimes()
	p_http_v1_1.EXPECT().Instance().Return(&http{v: "1.1.0"}).AnyTimes()

	p_http_v1_3 := mock_plugin.NewMockPlugin(ctrl)
	p_http_v1_3.EXPECT().Meta().Return(&meta.Data{
		ID:           meta.ID{ID: "nori/Http", Version: "1.3.0"},
		Author:       meta.Author{},
		Dependencies: []meta.Dependency{},
		Description:  meta.Description{},
		Core:         meta.Core{},
		Interface:    "nori/Http@1.3.0",
		License:      nil,
		Links:        nil,
		Repository:   meta.Repository{},
		Tags:         nil,
	}).AnyTimes()
	p_http_v1_3.EXPECT().Instance().Return(&http{v: "1.3.0"}).AnyTimes()

	p_http_v2_0 := mock_plugin.NewMockPlugin(ctrl)
	p_http_v2_0.EXPECT().Meta().Return(&meta.Data{
		ID:           meta.ID{ID: "nori/Http", Version: "2.0.0"},
		Author:       meta.Author{},
		Dependencies: []meta.Dependency{},
		Description:  meta.Description{},
		Core:         meta.Core{},
		Interface:    "nori/Http@2.0.0",
		License:      nil,
		Links:        nil,
		Repository:   meta.Repository{},
		Tags:         nil,
	}).AnyTimes()
	p_http_v2_0.EXPECT().Instance().Return(&http{v: "2.0.0"}).AnyTimes()

	err := r.Add(p_http_v1_0)
	a.NoError(err)
	err = r.Add(p_http_v1_1)
	a.NoError(err)
	err = r.Add(p_http_v1_3)
	a.NoError(err)
	err = r.Add(p_http_v2_0)
	a.NoError(err)

	p, err := r.Resolve(dep)
	a.NoError(err)
	a.NotNil(p)

	i, ok := (p).(Http)
	a.True(ok)
	if ok {
		a.Equal("1.3.0", i.Version())
	}
}

type Http interface {
	Version() string
}

type http struct {
	v string
}

func (h *http) Version() string {
	return h.v
}
