package storage

import (
	"database/sql"
	"encoding/json"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/secure2work/nori/core/plugins/meta"
	"github.com/sirupsen/logrus"
)

type mysql struct {
	db  *sql.DB
	log *logrus.Logger
}

func getMySqlStorage(cfg noriStorage) (NoriStorage, error) {
	db, err := sql.Open("mysql", cfg.Source)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	db.Exec(`CREATE TABLE IF NOT EXISTS nori_plugins (
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
		)  ENGINE=MyISAM;`)

	return &mysql{
		db: db,
	}, nil
}

func (m *mysql) GetPluginMetas() ([]meta.Meta, error) {
	var metas []meta.Meta
	rows, err := m.db.Query(
		`SELECT id,
				version,
				author,
				deps,
				description,
				core,
				interface,
				license,
				links,
				tags,
				installed, 
				hash FROM nori_plugins`)
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

func (m *mysql) SavePluginMeta(meta meta.Meta) error {
	statement, err := m.db.Prepare("INSERT INTO nori_plugins VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
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

	// @todo store file hash (?)
	_, err = statement.Exec(
		meta.Id().ID,
		meta.Id().Version,
		author,
		deps,
		description,
		core,
		int(meta.GetInterface()),
		license,
		links,
		strings.Join(meta.GetTags(), ","),
	)
	return err
}

func (m *mysql) DeletePluginMeta(id meta.ID) error {
	_, err := m.db.Exec(
		"DELETE FROM nori_plugins WHERE id = ? AND version = ? LIMIT 1",
		id.ID, id.Version)
	return err
}
