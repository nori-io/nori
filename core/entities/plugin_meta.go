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

package entities

type PluginMeta interface {
	GetId() string

	GetAuthor() string
	GetAuthorURI() string

	GetDependencies() []string
	GetDescription() string

	GetInterface() PluginInterface

	GetLicense() string
	GetLicenseURI() string

	GetPluginName() string
	GetPluginURI() string

	GetTags() []string

	GetVersion() string
}

type PluginMetaStruct struct {
	Id           string
	Author       string
	AuthorURI    string
	Dependencies []string
	Description  string
	Interface    PluginInterface
	License      string
	LicenseURI   string
	PluginName   string
	PluginURI    string
	Tags         []string
	Version      string
}

func (p *PluginMetaStruct) GetId() string {
	return p.Id
}

func (p *PluginMetaStruct) GetAuthor() string {
	return p.Author
}

func (p *PluginMetaStruct) GetAuthorURI() string {
	return p.AuthorURI
}

func (p *PluginMetaStruct) GetDependencies() []string {
	return p.Dependencies
}

func (p *PluginMetaStruct) GetDescription() string {
	return p.Description
}

func (p *PluginMetaStruct) GetInterface() PluginInterface {
	return p.Interface
}

func (p *PluginMetaStruct) GetLicense() string {
	return p.License
}

func (p *PluginMetaStruct) GetLicenseURI() string {
	return p.LicenseURI
}

func (p *PluginMetaStruct) GetPluginName() string {
	return p.PluginName
}

func (p *PluginMetaStruct) GetPluginURI() string {
	return p.PluginURI
}

func (p *PluginMetaStruct) GetTags() []string {
	return p.Tags
}

func (p *PluginMetaStruct) GetVersion() string {
	return p.Version
}
