package migrations

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sirupsen/logrus"
)

func Migrate(connString string) {
	logrus.Infof("Migrations up")
	m, err := migrate.New(
		"file://migrations/versions",
		connString,
	)

	if err != nil {
		logrus.Fatalf("Unable to find migration files cause: %s", err.Error())
	}

	err = m.Up()

	if err == migrate.ErrNoChange {
		logrus.Infof("Apply migrations: %s", err.Error())
		return
	}

	if err != nil {
		logrus.Fatalf("Failed to up migrations, cause: %s", err.Error())
	}

	if err == nil {
		logrus.Infof("Migrations up to date")
	}
}
