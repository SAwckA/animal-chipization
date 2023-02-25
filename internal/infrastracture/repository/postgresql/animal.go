package psql

import (
	"animal-chipization/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

const (
	animalTable                 = "public.animal"
	animalTypesListTable        = "public.animal_types_list"
	animalVisitedLocationsTable = "public.animal_locations_list"
)

type AnimalRepository struct {
	db *sqlx.DB
}

func NewAnimalRepository(db *sqlx.DB) *AnimalRepository {
	return &AnimalRepository{db: db}
}

func (r *AnimalRepository) GetAnimal(animalID int) (*domain.Animal, error) {

	query := fmt.Sprintf(`
	with locations as (
		select 
			tmp.animal_id as animal_id,
			json_agg(tmp.obj) as locations_list 
		from (
			select
				animal_id,
				json_build_object(
					'id', all2.id,
					'animal_id', all2.animal_id,
					'location_id', all2.location_id,
					'date_time_of_visited_location_point', all2.date_time_of_visited_location_point 
				) as obj
			from %s all2
			order by all2.date_time_of_visited_location_point
		) as tmp
		group by tmp.animal_id
	), 
	types1 as (
		select
			atl.animal_id as animal_id,
			jsonb_agg(atl.type_id) as types_list
		from %s atl
		group by atl.animal_id 
	)
	select 
		an.*,
		types1.types_list,
		locations.locations_list
	from %s an
	left join locations on locations.animal_id = an.id
	left join types1 on types1.animal_id = an.id
	where an.id = $1`,
		animalVisitedLocationsTable,
		animalTypesListTable,
		animalTable,
	)

	row := r.db.QueryRow(query, animalID)

	var typesString *string
	var visitedLocationString *string
	var animal domain.Animal

	if err := row.Scan(
		&animal.ID,
		&animal.Weight,
		&animal.Lenght,
		&animal.Height,
		&animal.Gender,
		&animal.LifeStatus,
		&animal.ChippingDateTime,
		&animal.ChipperID,
		&animal.ChippingLocationId,
		&animal.DeathDateTime,
		&typesString,
		&visitedLocationString,
	); err != nil {
		return nil, &domain.ApplicationError{
			OriginalError: err,
			SimplifiedErr: domain.ErrNotFound,
			Description:   "animal not found by id",
		}
	}

	if typesString != nil {
		json.Unmarshal([]byte(*typesString), &animal.AnimalTypes)
	} else {
		animal.AnimalTypes = make([]int, 0)
	}

	if visitedLocationString != nil {
		json.Unmarshal([]byte(*visitedLocationString), &animal.VisitedLocations)
	} else {
		animal.VisitedLocations = make([]domain.VisitedLocation, 0)
	}

	return &animal, nil
}

func (r *AnimalRepository) SearchAnimal(params *domain.AnimalSearchParams) (*[]domain.Animal, error) {

	var searchParams []string
	var searchData []interface{}
	placeholder := 1

	if params.StartDateTime != nil {
		searchParams = append(searchParams, "an.chippingdatetime > $1")
		searchData = append(searchData, params.StartDateTime)
		placeholder++
	}

	if params.EndDateTime != nil {
		searchParams = append(searchParams, "an.chippingdatetime < $2")
		searchData = append(searchData, params.EndDateTime)
		placeholder++
	}

	if params.ChipperID != nil {
		searchParams = append(searchParams, fmt.Sprintf(`an.chipperid = $%d`, placeholder))
		searchData = append(searchData, params.ChipperID)
		placeholder++
	}
	if params.ChippedLocationID != nil {
		searchParams = append(searchParams, fmt.Sprintf(`an.chippinglocationid = $%d`, placeholder))
		searchData = append(searchData, params.ChippedLocationID)
		placeholder++
	}
	if params.LifeStatus != nil {
		searchParams = append(searchParams, fmt.Sprintf(`an.lifestatus = $%d`, placeholder))
		searchData = append(searchData, params.LifeStatus)
		placeholder++
	}
	if params.Gender != nil {
		searchParams = append(searchParams, fmt.Sprintf(`an.gender = $%d`, placeholder))
		searchData = append(searchData, params.Gender)
		placeholder++
	}

	isSearch := ""
	if len(searchData) > 0 {
		isSearch = "where"
	}

	query := fmt.Sprintf(`
	with locations as (
		select 
			tmp.animal_id as animal_id,
			json_agg(tmp.obj) as locations_list 
		from (
			select
				animal_id,
				json_build_object(
					'id', all2.id,
					'animal_id', all2.animal_id,
					'location_id', all2.location_id,
					'date_time_of_visited_location_point', all2.date_time_of_visited_location_point 
				) as obj
			from %s all2
			order by all2.date_time_of_visited_location_point
		) as tmp
		group by tmp.animal_id
	), 
	types1 as (
		select
			atl.animal_id as animal_id,
			jsonb_agg(atl.type_id) as types_list
		from %s atl
		group by atl.animal_id 
	)
	select 
		an.*,
		types1.types_list,
		locations.locations_list
	from %s an
	left join locations on locations.animal_id = an.id
	left join types1 on types1.animal_id = an.id
	%s
		%s
	offset $%d
	limit $%d`,
		animalVisitedLocationsTable,
		animalTypesListTable,
		animalTable,
		isSearch,
		strings.Join(searchParams, " and "),
		placeholder,
		placeholder+1,
	)
	searchData = append(searchData, params.From)
	searchData = append(searchData, params.Size)

	rows, err := r.db.Query(query, searchData...)

	if err != nil {
		return nil, &domain.ApplicationError{
			OriginalError: err,
			SimplifiedErr: domain.ErrUnknown,
			Description:   "invalid query",
		}
	}

	var res []domain.Animal

	for rows.Next() {
		var animal domain.Animal
		var typesString *string
		var visitedLocationsString *string
		if err := rows.Scan(
			&animal.ID,
			&animal.Weight,
			&animal.Lenght,
			&animal.Height,
			&animal.Gender,
			&animal.LifeStatus,
			&animal.ChippingDateTime,
			&animal.ChipperID,
			&animal.ChippingLocationId,
			&animal.DeathDateTime,
			&typesString,
			&visitedLocationsString,
		); err != nil {
			return nil, &domain.ApplicationError{
				OriginalError: err,
				SimplifiedErr: domain.ErrNotFound,
				Description:   "0 animals found in search",
			}
		}

		if typesString != nil {
			json.Unmarshal([]byte(*typesString), &animal.AnimalTypes)
		} else {
			animal.AnimalTypes = make([]int, 0)
		}

		if visitedLocationsString != nil {
			json.Unmarshal([]byte(*visitedLocationsString), &animal.VisitedLocations)
		} else {
			animal.VisitedLocations = make([]domain.VisitedLocation, 0)
		}

		res = append(res, animal)
	}
	return &res, nil
}

func (r *AnimalRepository) CreateAnimal(animal *domain.Animal) (int, error) {

	tx, err := r.db.BeginTx(context.TODO(), nil)
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}

		tx.Commit()
	}()

	if err != nil {
		return 0, err
	}

	query := fmt.Sprintf(`
	insert into %s(
		weight, 
		length, 
		height, 
		gender, 
		lifestatus, 
		chippingdatetime, 
		chipperid, 
		chippinglocationid, 
		deathdatetime
	) values ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	
	returning id
	`, animalTable)

	row := tx.QueryRow(query,
		animal.Weight,
		animal.Lenght,
		animal.Height,
		animal.Gender,
		animal.LifeStatus,
		animal.ChippingDateTime,
		animal.ChipperID,
		animal.ChippingLocationId,
		animal.DeathDateTime,
	)

	var id int
	if err := row.Scan(&id); err != nil {
		if strings.Contains(err.Error(), "animal_chipperid_fkey") {
			return 0, &domain.ApplicationError{
				OriginalError: err,
				SimplifiedErr: domain.ErrNotFound,
				Description:   "Account not found by id",
			}
		}

		if strings.Contains(err.Error(), "animal_chippinglocationid_fkey") {
			return 0, &domain.ApplicationError{
				OriginalError: err,
				SimplifiedErr: domain.ErrNotFound,
				Description:   "Location not found by id",
			}
		}

		return 0, &domain.ApplicationError{
			OriginalError: err,
			SimplifiedErr: domain.ErrUnknown,
			Description:   "unknown error",
		}
	}

	baseQuery := fmt.Sprintf(`
		insert into %s(animal_id, type_id) 
			values 
	`, animalTypesListTable)

	var argsQuery []string
	argValues := make([]interface{}, len(animal.AnimalTypes))

	for index, value := range animal.AnimalTypes {
		argValues[index] = value
		argsQuery = append(argsQuery, fmt.Sprintf(`(%d, $%d)`, id, index+1))

	}
	_, err = tx.Exec(fmt.Sprintf("%s %s", baseQuery, strings.Join(argsQuery, ",")), argValues...)

	if err != nil {
		if strings.Contains(err.Error(), "animal_types_list_type_id_fkey") {
			return 0, &domain.ApplicationError{
				OriginalError: err,
				SimplifiedErr: domain.ErrNotFound,
				Description:   "Animal type not found by id",
			}
		}
		return 0, &domain.ApplicationError{
			OriginalError: err,
			SimplifiedErr: domain.ErrUnknown,
			Description:   "unknown error",
		}
	}
	return id, nil
}

// UpdateAnimal не обновляет поля AnimalType и VisitedLocations,
// для этого используются другие методы предназначенные для
// изменений только этих полей
func (r *AnimalRepository) UpdateAnimal(animal domain.Animal) error {

	query := fmt.Sprintf(`
	update %s
	set
		
		length = $1,
		weight = $2,
		height = $3,
		gender = $4,
		lifestatus = $5,
		chipperid = $6,
		chippinglocationid = $7,
		deathDateTime = $8

	where
		id = $9
	`, animalTable)

	res, err := r.db.Exec(query,
		animal.Lenght,
		animal.Weight,
		animal.Height,
		animal.Gender,
		animal.LifeStatus,
		animal.ChipperID,
		animal.ChippingLocationId,
		animal.DeathDateTime,
		animal.ID,
	)

	if err != nil {
		if strings.Contains(err.Error(), "animal_chipperid_fkey") {
			return &domain.ApplicationError{
				OriginalError: err,
				SimplifiedErr: domain.ErrNotFound,
				Description:   "Account not found by id",
			}
		}

		if strings.Contains(err.Error(), "animal_chippinglocationid_fkey") {
			return &domain.ApplicationError{
				OriginalError: err,
				SimplifiedErr: domain.ErrNotFound,
				Description:   "Location not found by id",
			}
		}

		return &domain.ApplicationError{
			OriginalError: err,
			SimplifiedErr: domain.ErrUnknown,
			Description:   "unknown error",
		}
	}

	affected, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if affected != 1 {
		return &domain.ApplicationError{
			OriginalError: err,
			SimplifiedErr: domain.ErrNotFound,
			Description:   "Animal not found by id",
		}
	}

	return nil
}

func (r *AnimalRepository) DeleteAnimal(animalID int) error {

	query := fmt.Sprintf(`delete from %s where id = $1`, animalTable)

	res, err := r.db.Exec(query, animalID)

	if err != nil {
		return &domain.ApplicationError{
			OriginalError: err,
			SimplifiedErr: domain.ErrUnknown,
			Description:   "Account not found by id",
		}
	}

	if affected, err := res.RowsAffected(); err != nil || affected != 1 {
		return &domain.ApplicationError{
			OriginalError: err,
			SimplifiedErr: domain.ErrNotFound,
			Description:   "Animal not found by id",
		}
	}

	return nil
}

func (r *AnimalRepository) AttachTypeAnimal(animalID, typeID int) error {
	query := fmt.Sprintf(`
	insert into %s(animal_id, type_id) 
		values 
	($1, $2)
	`, animalTypesListTable)

	_, err := r.db.Exec(query, animalID, typeID)

	if err != nil {
		return &domain.ApplicationError{
			OriginalError: err,
			SimplifiedErr: domain.ErrAlreadyExist,
			Description:   "animal already have this type",
		}
	}

	return nil
}

func (r *AnimalRepository) EditAnimalType(animalID, oldTypeID, newTypeID int) error {

	query := fmt.Sprintf(`
	update %s
	set
		type_id = $1
	where
		animal_id = $2 and type_id = $3
	`, animalTypesListTable)

	_, err := r.db.Exec(query, newTypeID, animalID, oldTypeID)

	if err != nil {
		return &domain.ApplicationError{
			OriginalError: err,
			SimplifiedErr: domain.ErrAlreadyExist,
			Description:   "animal already have this type",
		}
	}

	return nil
}

func (r *AnimalRepository) DeleteAnimalType(animalID, typeID int) error {

	query := fmt.Sprintf(`
	
	delete from %s
	where
		animal_id = $1 and type_id = $2

	`, animalTypesListTable)

	_, err := r.db.Exec(query, animalID, typeID)

	if err != nil {
		return &domain.ApplicationError{
			OriginalError: err,
			SimplifiedErr: domain.ErrAlreadyExist,
			Description:   "animal already have this type",
		}
	}
	return nil
}
