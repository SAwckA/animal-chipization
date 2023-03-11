package psql

import (
	"animal-chipization/internal/domain"
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

func (r *LocationRepository) Location(id int) (*domain.Location, error) {
	query := fmt.Sprintf(`
	select id, latitude, longitude from %s where id=$1
	`, locationTable)

	var location domain.Location
	if err := r.db.QueryRow(query, id).Scan(&location.ID, &location.Latitude, &location.Longitude); err != nil {
		return nil, &domain.ApplicationError{
			OriginalError: err,
			SimplifiedErr: domain.ErrNotFound,
			Description:   "location not found by id",
		}
	}

	return &location, nil
}

func (r *LocationRepository) Create(lat, lon float64) (int, error) {
	query := fmt.Sprintf(`
	insert into %s(latitude, longitude)
	values ($1, $2)
	returning id
	`, locationTable)

	var locationID int
	if err := r.db.QueryRow(query, lat, lon).Scan(&locationID); err != nil {
		return 0, &domain.ApplicationError{
			OriginalError: err,
			SimplifiedErr: domain.ErrAlreadyExist,
			Description:   "location already exist",
		}
	}

	return locationID, nil
}

func (r *LocationRepository) Update(location *domain.Location) error {
	query := fmt.Sprintf(`
	update %s 
	set latitude = $1,
		longitude = $2
	where id = $3
	`, locationTable)

	result, err := r.db.Exec(query, location.Latitude, location.Longitude, location.ID)
	if err != nil {
		return &domain.ApplicationError{
			OriginalError: err,
			SimplifiedErr: domain.ErrAlreadyExist,
			Description:   "location already exist",
		}
	}

	if affected, _ := result.RowsAffected(); affected != 1 {
		return &domain.ApplicationError{
			OriginalError: nil,
			SimplifiedErr: domain.ErrNotFound,
			Description:   "location not found by id",
		}
	}

	return nil
}

func (r *LocationRepository) Delete(id int) error {

	query := fmt.Sprintf(`
	delete from %s
	where id = $1
	`, locationTable)

	result, err := r.db.Exec(query, id)
	if err != nil {
		return &domain.ApplicationError{
			OriginalError: err,
			SimplifiedErr: domain.ErrInvalidInput,
			Description:   "location linked with animal visited location",
		}
	}

	affected, err := result.RowsAffected()
	if affected == 0 {
		return &domain.ApplicationError{
			OriginalError: err,
			SimplifiedErr: domain.ErrNotFound,
			Description:   "location not found by id",
		}
	}

	return err
}
