package nori_test

//import (
//	"context"
//	"testing"
//
//	"github.com/golang/mock/gomock"
//	dommeta "github.com/nori-io/common/v5/pkg/domain/meta"
//	logger "github.com/nori-io/common/v5/pkg/domain/mocks/logger"
//	plugin "github.com/nori-io/common/v5/pkg/domain/mocks/plugin"
//	"github.com/nori-io/common/v5/pkg/meta"
//	"github.com/nori-io/nori/pkg/nori"
//	"github.com/nori-io/nori/pkg/nori/domain/enum"
//)
//
//func TestNori_Add(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	// todo: mock
//	fl := logger.NewMockLogger(ctrl)
//
//	p1 := plugin.NewMockPlugin(ctrl)
//	p1.EXPECT().Meta().Return(meta.Meta{
//		ID: meta.ID{
//			ID:      "nori/Basic",
//			Version: "1.0.0",
//		},
//		Author:       nil,
//		Dependencies: []dommeta.Dependency{},
//		Description:  meta.Description{},
//		Interface:    dommeta.NewInterface("nori/Basic", "1.0.0"),
//		License:      []dommeta.License{},
//		Links:        nil,
//		Repository:   nil,
//		Tags:         nil,
//	}).AnyTimes()
//
//	n, err := nori.New(fl)
//	if err != nil {
//		t.Error(err)
//	}
//	err = n.Add(p1)
//	if err != nil {
//		t.Error(err)
//	}
//
//	err = n.Add(p1)
//	if err == nil {
//		t.Error("must return error: plugin [] already exists")
//	}
//}
//
//func TestNori_Dependent(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	fl := logger.NewMockLogger(ctrl)
//
//	p1 := plugin.NewMockPlugin(ctrl)
//	p1.EXPECT().Meta().Return(meta.Meta{
//		ID: meta.ID{
//			ID:      "nori/Basic",
//			Version: "1.0.0",
//		},
//		Author:       nil,
//		Dependencies: []dommeta.Dependency{},
//		Description:  meta.Description{},
//		Interface:    dommeta.NewInterface("nori/Basic", "1.0.0"),
//		License:      []dommeta.License{},
//		Links:        nil,
//		Repository:   nil,
//		Tags:         nil,
//	}).AnyTimes()
//	p1.EXPECT().Init(gomock.Any(), gomock.Any(), gomock.Any())
//	p1.EXPECT().Start(gomock.Eq(context.Background()), gomock.Any())
//
//	p2 := plugin.NewMockPlugin(ctrl)
//	p2.EXPECT().Meta().Return(meta.Meta{
//		ID: meta.ID{
//			ID:      "nori/Dependent",
//			Version: "2.0.1",
//		},
//		Author: nil,
//		Dependencies: []dommeta.Dependency{
//			dommeta.NewInterface("nori/Basic", "1.0.0"),
//		},
//		Description: meta.Description{},
//		Interface:   dommeta.NewInterface("nori/Dependent", "1.0.0"),
//		License:     []dommeta.License{},
//		Links:       nil,
//		Repository:  nil,
//		Tags:        nil,
//	}).AnyTimes()
//	p2.EXPECT().Init(gomock.Any(), gomock.Any(), gomock.Any())
//	p2.EXPECT().Start(gomock.Eq(context.Background()), gomock.Any())
//
//	n, err := nori.New(fl)
//	if err != nil {
//		t.Error(err)
//	}
//
//	err = n.Add(p1)
//	if err != nil {
//		t.Error(err)
//	}
//	err = n.Add(p2)
//	if err != nil {
//		t.Error(err)
//	}
//
//	err = n.Init(context.Background(), p2.Meta().GetID())
//	if err != nil {
//		t.Error(err)
//	}
//
//	ids := n.GetByFilter(nori.Filter{State: enum.Inited})
//	if len(ids) != 2 {
//		t.Error("not all plugins inited")
//	}
//
//	err = n.Start(context.Background(), p2.Meta().GetID())
//	if err != nil {
//		t.Error(err)
//	}
//
//	ids = n.GetByFilter(nori.Filter{State: enum.Running})
//	if len(ids) != 2 {
//		t.Error("not all plugins running")
//	}
//}
//
//func TestNori_All(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	// todo: mock
//	fl := logger.NewMockLogger(ctrl)
//
//	p1 := plugin.NewMockPlugin(ctrl)
//	p1.EXPECT().Meta().Return(meta.Meta{
//		ID: meta.ID{
//			ID:      "nori/Basic",
//			Version: "1.0.0",
//		},
//		Author:       nil,
//		Dependencies: []dommeta.Dependency{},
//		Description:  meta.Description{},
//		Interface:    dommeta.NewInterface("nori/Basic", "1.0.0"),
//		License:      []dommeta.License{},
//		Links:        nil,
//		Repository:   nil,
//		Tags:         nil,
//	}).AnyTimes()
//	p1.EXPECT().Init(gomock.Eq(context.Background()), gomock.Any(), gomock.Any()).AnyTimes()
//	p1.EXPECT().Start(gomock.Eq(context.Background()), gomock.Any()).AnyTimes()
//
//	p2 := plugin.NewMockPlugin(ctrl)
//	p2.EXPECT().Meta().Return(meta.Meta{
//		ID: meta.ID{
//			ID:      "nori/Dependent",
//			Version: "2.0.1",
//		},
//		Author: nil,
//		Dependencies: []dommeta.Dependency{
//			dommeta.NewInterface("nori/Basic", "1.0.0"),
//		},
//		Description: meta.Description{},
//		Interface:   dommeta.NewInterface("nori/Dependent", "1.0.0"),
//		License:     []dommeta.License{},
//		Links:       nil,
//		Repository:  nil,
//		Tags:        nil,
//	}).AnyTimes()
//	p2.EXPECT().Init(gomock.Eq(context.Background()), gomock.Any(), gomock.Any()).AnyTimes()
//	p2.EXPECT().Start(gomock.Eq(context.Background()), gomock.Any()).AnyTimes()
//
//	n, err := nori.New(fl)
//	if err != nil {
//		t.Error(err)
//	}
//
//	err = n.Add(p1)
//	if err != nil {
//		t.Error(err)
//	}
//	err = n.Add(p2)
//	if err != nil {
//		t.Error(err)
//	}
//
//	err = n.InitAll(context.Background())
//	if err != nil {
//		t.Error(err)
//	}
//
//	ids := n.GetByFilter(nori.Filter{State: enum.Inited})
//	if len(ids) != 2 {
//		t.Error("not all plugins inited")
//	}
//
//	err = n.StartAll(context.Background())
//	if err != nil {
//		t.Error(err)
//	}
//
//	ids = n.GetByFilter(nori.Filter{State: enum.Running})
//	if len(ids) != 2 {
//		t.Error("not all plugins running")
//	}
//}
