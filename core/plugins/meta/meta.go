package meta

import (
	"fmt"

	"github.com/hashicorp/go-version"
)

type ID struct {
	ID      string
	Version string
}

func (id ID) GetVersion() (*version.Version, error) {
	v, err := version.NewVersion(id.Version)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (id ID) String() string {
	return fmt.Sprintf("%s:%s", id.ID, id.Version)
}

type Author struct {
	Name string
	URI  string
}

type Core struct {
	VersionConstraint string
}

type Dependency struct {
	ID         string
	Constraint string
}

func (d Dependency) GetConstraint() (version.Constraints, error) {
	constraints, err := version.NewConstraint(d.Constraint)
	if err != nil {
		return nil, err
	}
	return constraints, nil
}

type Description struct {
	Name        string
	Description string
}

type License struct {
	Title string
	Type  string
	URI   string
}

type Link struct {
	Title string
	URL   string
}

type Meta interface {
	Id() ID
	GetAuthor() Author
	GetDependencies() []Dependency
	GetDescription() Description
	GetCore() Core
	GetInterface() Interface
	GetLicense() License
	GetLinks() []Link
	GetTags() []string
}

type Data struct {
	ID           ID
	Author       Author
	Dependencies []Dependency
	Description  Description
	Core         Core
	Interface    Interface
	License      License
	Links        []Link
	Tags         []string
}

func (m Data) Id() ID {
	return m.ID
}

func (m Data) GetAuthor() Author {
	return m.Author
}

func (m Data) GetDependencies() []Dependency {
	return m.Dependencies
}

func (m Data) GetDescription() Description {
	return m.Description
}

func (m Data) GetCore() Core {
	return m.Core
}

func (m Data) GetInterface() Interface {
	return m.Interface
}

func (m Data) GetLicense() License {
	return m.License
}

func (m Data) GetLinks() []Link {
	return m.Links
}

func (m Data) GetTags() []string {
	return m.Tags
}
