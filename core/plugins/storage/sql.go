package storage

import (
	"database/sql"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/secure2work/nori/core/entities"
	"github.com/sirupsen/logrus"
)

type mysql struct {
	db  *sql.DB
	log *logrus.Logger
}

func getMySqlStorage(cfg noriCoreStorage) (NoriStorage, error) {
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
		    id varchar(255) NOT NULL,
		    author varchar(255),
		    author_uri varchar(255),
		    description varchar(255),
		    license varchar(255),
		    license_uri varchar(255),
		    plugin_name varchar(255),
		    plugin_uri varchar(1000),
		    tags varchar(255),
		    interface varchar(255),
		    version varchar(255),
		    dependencies text,
		    installed bigint,
		    hash varchar(255),
		    PRIMARY KEY (id)
		)  ENGINE=MyISAM;`)

	return &mysql{
		db: db,
	}, nil
}

func (m *mysql) GetInstallations() ([]entities.PluginMeta, error) {
	var meta []entities.PluginMeta
	rows, err := m.db.Query("SELECT id, author, author_uri, description, license, license_uri, plugin_name, plugin_uri, tags, interface, version, dependencies FROM nori_plugins")
	if err != nil {
		return meta, err
	}

	defer rows.Close()
	for rows.Next() {
		var ms entities.PluginMetaStruct
		var tags, dependencies string
		err := rows.Scan(&ms.Id, &ms.Author, &ms.AuthorURI, &ms.Description, &ms.License, &ms.LicenseURI,
			&ms.PluginName, &ms.PluginURI, &tags, &ms.Kind, &ms.Version, &dependencies)
		if err != nil {
			m.log.Error(err)
		}
		ms.Tags = strings.Split(tags, ",")
		ms.Dependencies = strings.Split(dependencies, ",")
		meta = append(meta, &ms)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return meta, nil
}

func (m *mysql) SaveInstallation(meta entities.PluginMeta) error {
	statement, err := m.db.Prepare("INSERT INTO nori_plugins VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		return err
	}
	// @todo store file hash (?)
	_, err = statement.Exec(meta.GetId(), meta.GetAuthor(), meta.GetAuthorURI(), meta.GetDescription(),
		meta.GetLicense(), meta.GetLicenseURI(), meta.GetPluginName(), meta.GetPluginURI(),
		strings.Join(meta.GetTags(), ","), meta.GetInterface(), meta.GetVersion(),
		strings.Join(meta.GetDependencies(), ","), time.Now().Unix(), "")

	return err
}

func (m *mysql) RemoveInstallation(id string) error {
	_, err := m.db.Exec("DELETE FROM nori_plugins WHERE id = ? LIMIT 1", id)
	return err
}
