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
	GetAuthor() string
	GetAuthorURI() string

	GetDependencies() []string
	GetDescription() string

	GetLicense() string
	GetLicenseURI() string

	GetId() string

	GetPluginName() string
	GetPluginURI() string

	GetTags() []string

	GetKind() PluginKind

	GetVersion() string
}

// @todo rewrite this comments, make it clear
// Structure Description:
// Author - author's name or company's name
// AuithorURI - link to author's web-site
// Dependencies []string - this is an array of PluginName
// Id - this is plugin's id in format {company}/{kind} or {team}/{technology}/{plugin-name}", ex: nori/http, nori/redis/cluster
// PluginName - this is the name of the plugin, ex: "HTML Pages Manager v1.0"
// PluginURI - this is a link to plugin documentation / web-page / support
type PluginMetaStruct struct {
	Id           string
	Author       string
	AuthorURI    string
	Dependencies []string
	Description  string
	License      string
	LicenseURI   string
	PluginName   string
	PluginURI    string
	Tags         []string
	Kind         PluginKind
	Version      string
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

func (p *PluginMetaStruct) GetLicense() string {
	return p.License
}

func (p *PluginMetaStruct) GetLicenseURI() string {
	return p.LicenseURI
}

func (p *PluginMetaStruct) GetId() string {
	return p.Id
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

func (p *PluginMetaStruct) GetKind() PluginKind {
	return p.Kind
}

func (p *PluginMetaStruct) GetVersion() string {
	return p.Version
}
