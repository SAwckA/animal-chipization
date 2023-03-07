package repository

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"time"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Инициализация нового соединения к mongo
func NewMongoClient(connString string) (*mongo.Client, error) {

	// Новый клиент
	client, err := mongo.NewClient(options.Client().ApplyURI(connString))

	if err != nil {
		return nil, err
	}

	// Инициализация соединения
	err = client.Connect(context.TODO())
	if err != nil {
		return nil, err
	}

	// Проверка соединения
	err = client.Ping(context.TODO(), nil)

	return client, err
}

type postgresConfig interface {
	DataSourceString() string
}

func NewPostgresDB(config postgresConfig) (db *sqlx.DB, err error) {
	var maxRetries = 5
	var timeoutRetry = 2 * time.Second

	for i := 0; i < maxRetries; i++ {
		db, err = sqlx.Open("pgx", config.DataSourceString())

		if err != nil {
			logrus.Warnln(fmt.Sprintf("[RETRY %d] cause: %s", i, err.Error()))
			time.Sleep(timeoutRetry)
			continue
		}

		err = db.Ping()
		if err != nil {
			logrus.Warnln(fmt.Sprintf("[RETRY %d] cause: %s", i, err.Error()))
			time.Sleep(timeoutRetry)
			continue
		}
		return
	}
	return
}
