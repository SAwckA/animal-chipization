package migrations

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sirupsen/logrus"
)

func Migrate(connString string) error {

	m, err := migrate.New(
		"file://migrations/versions",
		connString,
	)

	if err != nil {
		return err
	}

	err = m.Up()

	if err == migrate.ErrNoChange {
		logrus.Warnln(err)
		return nil
	}
	return err
}
