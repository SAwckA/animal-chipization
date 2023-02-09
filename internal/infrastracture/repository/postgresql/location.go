package psql

import (
	"animal-chipization/internal/domain"
	"animal-chipization/internal/errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

const locationTable = "public.location"

type LocationRepository struct {
	db *sqlx.DB
}

func NewLocationRepository(db *sqlx.DB) *LocationRepository {
	return &LocationRepository{db: db}
}

func (r *LocationRepository) CreateLocation(lat, lon float64) (int, error) {
	query := fmt.Sprintf(`
	insert into %s(latitude, longitude)
	values ($1, $2)
	returning id
	`, locationTable)

	row := r.db.QueryRow(query, lat, lon)

	var locationID int
	if err := row.Scan(&locationID); err != nil {
		return 0, errors.ErrAlreadyExist
	}

	return locationID, nil
}

func (r *LocationRepository) GetLocation(locationID int) *domain.Location {
	query := fmt.Sprintf(`
	select id, latitude, longitude from %s where id=$1
	`, locationTable)

	var location domain.Location
	row := r.db.QueryRow(query, locationID)
	err := row.Scan(&location.ID, &location.Latitude, &location.Longitude)

	if err != nil {
		return nil
	}

	return &location
}

func (r *LocationRepository) UpdateLocation(location *domain.Location) error {
	query := fmt.Sprintf(`
	update %s 
	set latitude = $1,
		longitude = $2

	where id = $3
	`, locationTable)

	result, err := r.db.Exec(query, location.Latitude, location.Longitude, location.ID)

	if err != nil {
		return errors.ErrAlreadyExist
	}

	if affected, _ := result.RowsAffected(); affected != 1 {
		return errors.ErrNotFound
	}

	return nil
}

func (r *LocationRepository) DeleteLocation(locationID int) error {

	query := fmt.Sprintf(`
	delete from %s
	where id = $1
	`, locationTable)

	result, err := r.db.Exec(query, locationID)

	if err != nil {
		return errors.ErrLinked
	}

	affected, err := result.RowsAffected()

	if affected == 0 {
		return errors.ErrNotFound
	}

	return err
}
