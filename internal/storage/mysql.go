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

	"github.com/nori-io/nori-common/logger"

	_ "github.com/go-sql-driver/mysql"
	"github.com/nori-io/nori-common/meta"
)

type mysql struct {
	plugins Plugins
	db      *sql.DB
	log     logger.Logger
}

type mysqlPlugins struct {
	db  *sql.DB
	log logger.Logger
}

func newMySQLStorage(source string, log logger.Logger) (Storage, error) {
	db, err := sql.Open("mysql", source)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS nori_plugins (
	id varchar(255) NOT NULL ,
	version varchar(32) not null ,
	author text ,
	deps text ,
	description text ,
	core text,
	interface varchar(255) ,
	license text ,
	links text ,
	tags varchar(255) ,
	installed bigint ,
	hash varchar(255) ,
	PRIMARY KEY (id, version)
)  ENGINE=InnoDB;`)

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
	rows, err := m.db.Query("SELECT id, version, author, deps, description, core, interface, license, links, tags, installed, hash FROM nori_plugins")
	if err != nil {
		return metas, err
	}

	defer rows.Close()
	for rows.Next() {
		mt, err := mysqlParseRow(rows.Scan)
		if err != nil {
			m.log.Error(err)
			continue
		}
		metas = append(metas, mt)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return metas, nil
}

func (m *mysqlPlugins) Get(id meta.ID) (meta.Meta, error) {
	stmtGet, err := m.db.Prepare("SELECT * FROM nori_plugins WHERE id = ? AND version = ?")
	if err != nil {
		return nil, err
	}

	return mysqlParseRow(stmtGet.QueryRow(id.ID, id.Version).Scan)
}

func (m *mysqlPlugins) Save(meta meta.Meta) error {
	statement, err := m.db.Prepare("INSERT INTO nori_plugins VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		return err
	}
	var author, deps, description, core, license, links []byte

	author, _ = json.Marshal(meta.GetAuthor())
	deps, _ = json.Marshal(meta.GetDependencies())
	description, _ = json.Marshal(meta.GetDescription())
	core, _ = json.Marshal(meta.GetCore())
	license, _ = json.Marshal(meta.GetLicense())
	links, _ = json.Marshal(meta.GetLinks())

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
	stmtDelete, err := m.db.Prepare("DELETE FROM nori_plugins WHERE id = ? AND version = ?")
	if err != nil {
		return err
	}

	_, err = stmtDelete.Exec(id.ID, id.Version)
	return err
}

func (m *mysqlPlugins) IsInstalled(id meta.ID) (bool, error) {
	var exists bool
	err := m.db.QueryRow(
		"SELECT exists (SELECT id FROM nori_plugins WHERE id = ? AND version = ?)",
		id.ID,
		id.Version,
	).Scan(&exists)

	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	return exists, nil
}

func mysqlParseRow(scan func(dest ...interface{}) error) (meta.Meta, error) {
	var m meta.Data
	var author, description, dependencies, core, iface, license, links, tags string
	err := scan(
		&m.ID.ID, &m.ID.Version,
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
		return nil, err
	}

	var authData meta.Author
	err = json.Unmarshal([]byte(author), &authData)
	if err == nil {
		m.Author = authData
	}

	var depSet []meta.Dependency
	err = json.Unmarshal([]byte(dependencies), &depSet)
	if err == nil {
		m.Dependencies = depSet
	}

	var descData meta.Description
	err = json.Unmarshal([]byte(description), &descData)
	if err == nil {
		m.Description = descData
	}

	var coreData meta.Core
	err = json.Unmarshal([]byte(core), coreData)
	if err == nil {
		m.Core = coreData
	}

	m.Interface = meta.Interface(iface)

	var licenseData meta.License
	err = json.Unmarshal([]byte(license), &licenseData)
	if err == nil {
		m.License = licenseData
	}

	var linkSet []meta.Link
	err = json.Unmarshal([]byte(links), &linkSet)
	if err == nil {
		m.Links = linkSet
	}
	m.Tags = strings.Split(tags, ",")

	return m, nil
}
