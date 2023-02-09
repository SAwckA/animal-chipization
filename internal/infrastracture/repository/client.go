package repository

import (
	"context"
	"fmt"

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

func NewPostgresDB(host, port, user, dbname, password, sslmode string) (*sqlx.DB, error) {
	db, err := sqlx.Open("pgx", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		host, port, user, dbname, password, sslmode))

	if err != nil {
		return nil, err
	}

	errPing := db.Ping()
	if errPing != nil {
		return nil, errPing
	}

	return db, nil
}
