package repository

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sirupsen/logrus"
	"time"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

type postgresConfig interface {
	DataSourceString() string
	ConnString() string
}

func NewPostgresDB(config postgresConfig) (db *sqlx.DB, err error) {
	var maxRetries = 5
	var timeoutRetry = 2 * time.Second

	for i := 0; i < maxRetries; i++ {
		db, err = sqlx.Open("pgx", config.DataSourceString())

		if err != nil {
			logrus.Warnln(fmt.Sprintf("[RETRY %d OF %d] cause: %s", i, maxRetries, err.Error()))
			time.Sleep(timeoutRetry)
			continue
		}

		err = db.Ping()
		if err != nil {
			logrus.Warnln(fmt.Sprintf("[RETRY %d OF %d] cause: %s", i, maxRetries, err.Error()))
			time.Sleep(timeoutRetry)
			continue
		}
		upMigrations(config.ConnString())
		return
	}
	return
}

func upMigrations(connString string) {
	logrus.Infof("Start migrations up")
	m, err := migrate.New(
		"file://migrations",
		connString,
	)

	if err != nil {
		logrus.Fatalf("Unable connect to database: %s", err.Error())
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
		logrus.Infof("Apply migrations: success")
	}
}
