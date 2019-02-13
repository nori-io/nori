// Copyright © 2018 Nori info@nori.io
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

package storage

import (
	"database/sql"
	"encoding/json"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/nori-io/nori-common/meta"
	"github.com/sirupsen/logrus"
)

type mysql struct {
	plugins Plugins
	db      *sql.DB
	log     *logrus.Logger
}

type mysqlPlugins struct {
	db  *sql.DB
	log *logrus.Logger
}

func getMySqlStorage(source string, log *logrus.Logger) (Storage, error) {
	db, err := sql.Open("mysql", source)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	_, err = db.Exec(sqlCreateTablePlugins)
	if err != nil {
		return nil, err
	}

	return &mysql{
		plugins: &mysqlPlugins{
			db:  db,
			log: log,
		},
		db:  db,
		log: log,
	}, nil
}

func (m *mysql) Plugins() Plugins {
	return m.plugins
}

func (m *mysqlPlugins) All() ([]meta.Meta, error) {
	var metas []meta.Meta
	rows, err := m.db.Query(sqlQueryPluginsAll)
	if err != nil {
		return metas, err
	}

	defer rows.Close()
	for rows.Next() {
		var mt meta.Data
		var iface int
		var author, description, dependencies, core, license, links, tags string
		err := rows.Scan(
			&mt.ID.ID, &mt.ID.Version,
			&author,
			&dependencies,
			&description,
			&core,
			&iface,
			&license,
			&links,
			&tags,
		)
		if err != nil {
			m.log.Error(err)
			continue
		}

		var authData meta.Author
		err = json.Unmarshal([]byte(author), authData)
		if err == nil {
			mt.Author = authData
		}

		var depSet []meta.Dependency
		err = json.Unmarshal([]byte(dependencies), depSet)
		if err == nil {
			mt.Dependencies = depSet
		}

		var descData meta.Description
		err = json.Unmarshal([]byte(description), descData)
		if err == nil {
			mt.Description = descData
		}

		var coreData meta.Core
		err = json.Unmarshal([]byte(core), coreData)
		if err == nil {
			mt.Core = coreData
		}

		mt.Interface = meta.Interface(iface)

		var licenseData meta.License
		err = json.Unmarshal([]byte(license), &licenseData)
		if err == nil {
			mt.License = licenseData
		}

		var linkSet []meta.Link
		err = json.Unmarshal([]byte(links), &linkSet)
		if err == nil {
			mt.Links = linkSet
		}
		mt.Tags = strings.Split(tags, ",")

		metas = append(metas, &mt)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return metas, nil
}

func (m *mysqlPlugins) Save(meta meta.Meta) error {
	statement, err := m.db.Prepare(sqlInsertPlugin)
	if err != nil {
		return err
	}
	var author, deps, description, core, license, links []byte

	author, err = json.Marshal(meta.GetAuthor())
	deps, err = json.Marshal(meta.GetDependencies())
	description, err = json.Marshal(meta.GetDescription())
	core, err = json.Marshal(meta.GetCore())
	license, err = json.Marshal(meta.GetLicense())
	links, err = json.Marshal(meta.GetLinks())

	_, err = statement.Exec(
		meta.Id().ID,
		meta.Id().Version,
		author,
		deps,
		description,
		core,
		meta.GetInterface(),
		license,
		links,
		strings.Join(meta.GetTags(), ","),
	)
	return err
}

func (m *mysqlPlugins) Delete(id meta.ID) error {
	_, err := m.db.Exec(sqlDeletePlugin, id.ID, id.Version)
	return err
}

const sqlQueryPluginsAll = `SELECT id, version, author, deps, description, core, interface, license, links, tags, installed, hash FROM nori_plugins`
const sqlInsertPlugin = `INSERT INTO nori_plugins VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
const sqlDeletePlugin = `DELETE FROM nori_plugins WHERE id = ? AND version = ? LIMIT 1`

const sqlCreateTablePlugins = `
CREATE TABLE IF NOT EXISTS nori_plugins (
	id varchar(255) NOT NULL ,
	version varchar(32) not null ,
	author text ,
	deps text ,
	description text ,
	core text,
	interface int ,
	license text ,
	links text ,
	tags varchar(255) ,
	installed bigint ,
	hash varchar(255) ,
	PRIMARY KEY (id, version)
)  ENGINE=MyISAM;
`
