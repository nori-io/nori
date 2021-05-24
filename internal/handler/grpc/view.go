package grpc

import (
	"github.com/nori-io/common/v5/pkg/domain/meta"
	"github.com/nori-io/nori-grpc/pkg/api/proto"
	"github.com/nori-io/nori/internal/domain/entity"
)

type View struct{}

func (v View) Meta(m meta.Meta) *proto.Meta {
	deps := make([]*proto.Dependency, len(m.GetDependencies()))
	//license := make([]*proto.License, len(m.GetLicense()))
	//links := make([]*proto.Link, len(m.GetLinks()))
	tags := make([]string, len(m.GetTags()))

	for _, d := range m.GetDependencies() {
		deps = append(deps, &proto.Dependency{
			Constraint: d.Constraint(),
			Interface:  d.Name(),
		})
	}

	//for _, l := range m.GetLicense() {
	//	license = append(license, &proto.License{
	//		Title: l.GetTitle(),
	//		Type:  l.GetType(),
	//		Uri:   l.GetURL(),
	//	})
	//}
	//
	//for _, l := range m.GetLinks() {
	//	links = append(links, &proto.Link{
	//		Title: l.Title,
	//		Url:   l.URL,
	//	})
	//}

	for _, t := range m.GetTags() {
		tags = append(tags, t)
	}

	return &proto.Meta{
		Id: &proto.ID{
			PluginId: string(m.GetID().GetID()),
			Version:  m.GetID().GetVersion(),
		},
		Author: &proto.Author{
			Name: m.GetAuthor().GetName(),
			Uri:  m.GetAuthor().GetURL(),
		},
		Dependencies: deps,
		//Description: &proto.Description{
		//	Description: m.GetDescription().GetDescription(),
		//	Name:        m.GetDescription().GetTitle(),
		//},
		//Interface: m.GetInterface().String(),
		//Licenses:  license,
		//Links:     links,
		//Repository: &proto.Repository{
		//	Uri:  m.GetRepository().GetURL(),
		//	Type: string(m.GetRepository().GetType()),
		//},
		Tags: tags,
	}
}

func (v View) Plugin(p *entity.Plugin) *proto.Plugin {
	//var s proto.Status
	//switch p.Status() {
	//case status.Initialised:
	//	s = proto.Status_Inited
	//case status.Started:
	//	s = proto.Status_Running
	//default:
	//	s = proto.Status_Stopped
	//}

	return &proto.Plugin{
		Meta:     v.Meta(p.Plugin.Meta()),
		IsActive: false,
		//IsInstallable: p.IsInstallable(),
		//Status:        s,
	}
}
