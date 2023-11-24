package postgres

import (
	"time"

	"github.com/golang-migrate/migrate/v4"

	// load the postgres migration driver
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	// load the file source
	_ "github.com/golang-migrate/migrate/v4/source/file"
	log "github.com/sirupsen/logrus"
)

// NewMigration ...
func NewMigration(datasource *Datasource) *Migration {
	return &Migration{
		datasource: datasource,
	}
}

// Run ...
func (m *Migration) Run() error {
	// if location is not db, do not execute
	source := m.datasource.AsFileSource()
	datasource := m.datasource.AsDatasourceString()
	log.
		WithField("source", source).
		WithField("datasource", datasource).
		Info("Running Database Migrations")
	migrations, err := migrate.New(source, datasource)
	if err != nil {
		log.WithError(err).Error("Unable to create migrations")
		return err
	}
	if err := migrations.Up(); err != nil {
		switch err {
		case migrate.ErrNilVersion:
			fallthrough
		case migrate.ErrInvalidVersion:
			fallthrough
		case migrate.ErrLockTimeout:
			log.WithError(err).Error("an error occurred during db migrations")
			return err
		case migrate.ErrLocked:
			log.Warn("Migrations is currently locked: waiting 5 seconds before retrying")
			time.Sleep(5 * time.Second)
			return m.Run()
		case migrate.ErrNoChange:
			log.Info("Database already at latest version. No changes were made")
		default:
			log.WithError(err).Error("error running migrations")
		}
	}
	return nil
}
