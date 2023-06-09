package psql

import (
	"animal-chipization/internal/domain"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
)

const animalVisitedLocationsTable = "public.animal_locations_list"

type VisitedLocationRepository struct {
	db *sqlx.DB
}

func NewVisitedLocationRepository(db *sqlx.DB) *VisitedLocationRepository {
	return &VisitedLocationRepository{
		db: db,
	}
}

func (r *VisitedLocationRepository) VisitedLocation(id int) (*domain.VisitedLocation, error) {
	query := fmt.Sprintf(`
		select id, animal_id, location_id, date_time_of_visited_location_point from %s
		where id = $1
	`, animalVisitedLocationsTable)

	var location domain.VisitedLocation
	err := r.db.QueryRow(query, id).Scan(&location.ID, &location.AnimalID, &location.LocationPointID, &location.DateTime)
	if err != nil {
		return nil, &domain.ApplicationError{
			OriginalError: err,
			SimplifiedErr: domain.ErrNotFound,
			Description:   "visited location not found by id",
		}
	}
	return &location, nil
}

func (r *VisitedLocationRepository) Search(animalID int, params *domain.SearchVisitedLocation) ([]domain.VisitedLocation, error) {
	args := []string{
		"animal_id = $3",
	}
	var data []interface{}
	placeholder := 4
	data = append(data, params.From, params.Size, animalID)

	if params.StartDateTime != nil {
		args = append(args, fmt.Sprintf("date_time_of_visited_location_point > $%d", placeholder))
		data = append(data, params.StartDateTime)
		placeholder++
	}

	if params.EndDateTime != nil {
		args = append(args, fmt.Sprintf("date_time_of_visited_location_point < $%d", placeholder))
		data = append(data, params.EndDateTime)
		placeholder++
	}

	query := fmt.Sprintf(`
		select 
			id, location_id, date_time_of_visited_location_point
		from %s
		where 
			%s
		order by date_time_of_visited_location_point
		offset $1
		limit $2
	`, animalVisitedLocationsTable, strings.Join(args, " and "))

	rows, err := r.db.Query(query, data...)
	if err != nil {
		return nil, &domain.ApplicationError{
			OriginalError: err,
			SimplifiedErr: domain.ErrUnknown,
			Description:   "unknown error during search visited location point",
		}
	}

	var res []domain.VisitedLocation
	for rows.Next() {
		var location domain.VisitedLocation
		_ = rows.Scan(&location.ID, &location.LocationPointID, &location.DateTime)
		res = append(res, location)
	}

	return res, nil
}

func (r *VisitedLocationRepository) Save(animalID int, location *domain.VisitedLocation) (int, error) {
	query := fmt.Sprintf(`
		insert into %s(animal_id, location_id, date_time_of_visited_location_point)
			values
		($1, $2, $3)
		returning id
	`, animalVisitedLocationsTable)

	var locationID int
	err := r.db.Get(&locationID, query, animalID, location.LocationPointID, location.DateTime)
	if err != nil {
		return 0, &domain.ApplicationError{
			OriginalError: err,
			SimplifiedErr: domain.ErrUnknown,
			Description:   "unknown error during save visited location point",
		}
	}

	return locationID, nil
}

func (r *VisitedLocationRepository) Update(visitedLocation *domain.VisitedLocation) error {
	query := fmt.Sprintf(`
		update %s
		set location_id = $1
		where id = $2
	`, animalVisitedLocationsTable)

	res, err := r.db.Exec(query, visitedLocation.LocationPointID, visitedLocation.ID)
	if err != nil {
		return &domain.ApplicationError{
			OriginalError: err,
			SimplifiedErr: domain.ErrUnknown,
			Description:   "unknown error during update visited location point",
		}
	}

	if aff, err := res.RowsAffected(); err != nil || aff != 1 {
		return &domain.ApplicationError{
			OriginalError: err,
			SimplifiedErr: domain.ErrNotFound,
			Description:   "nothing updated",
		}
	}
	return nil
}

func (r *VisitedLocationRepository) Delete(id int) error {
	query := fmt.Sprintf(`
		delete from %s
		where id = $1
	`, animalVisitedLocationsTable)

	res, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	if aff, err := res.RowsAffected(); err != nil || aff != 1 {
		return &domain.ApplicationError{
			OriginalError: err,
			SimplifiedErr: domain.ErrNotFound,
			Description:   "visited location not found by id",
		}
	}

	return nil
}
