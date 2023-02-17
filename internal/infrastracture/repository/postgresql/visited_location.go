package psql

import "github.com/jmoiron/sqlx"

type VisitedLocationRepository struct {
	db *sqlx.DB
}

func NewVisitedLocationRepository(db *sqlx.DB) *VisitedLocationRepository {
	return &VisitedLocationRepository{
		db: db,
	}
}
