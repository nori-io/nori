package dep

import (
	"strings"

	"github.com/hashicorp/go-version"
)

type Dependency interface {
	Id() string
	Version() *version.Version
}

func NewDependencies(deps []string) ([]Dependency, error) {
	var list []Dependency

	for _, dep := range deps {
		meta := strings.Split(dep, ":")
		if len(meta) == 2 {
			ver = meta[1]
		}
		d, err := NewDependency(meta[0], ver)
		if err != nil {
			return nil, err
		}
		list = append(list, d)
	}

	return list, nil
}

func NewDependency(id, ver string) (Dependency, error) {
	v, err := version.NewVersion(ver)
	if err != nil {
		return nil, err
	}

	return dependency{
		id:  id,
		ver: v,
	}, nil
}

type dependency struct {
	id  string
	ver *version.Version
}

func (d dependency) Id() string {
	return d.id
}

func (d dependency) Version() *version.Version {
	return d.ver
}
