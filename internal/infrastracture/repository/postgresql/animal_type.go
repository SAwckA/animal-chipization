package psql

import (
	"animal-chipization/internal/domain"
	"animal-chipization/internal/errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

const animalTypeTable = "public.animal_type"

type AnimalTypeRepository struct {
	db *sqlx.DB
}

func NewAnimalTypeRepository(db *sqlx.DB) *AnimalTypeRepository {
	return &AnimalTypeRepository{db: db}
}

func (r *AnimalTypeRepository) GetAnimalType(typeID int) *domain.AnimalType {
	query := fmt.Sprintf(`select id, type from %s where id = $1`, animalTypeTable)

	var animalType domain.AnimalType
	row := r.db.QueryRow(query, typeID)

	if err := row.Scan(&animalType.ID, &animalType.Type); err != nil {
		return nil
	}

	return &animalType
}

func (r *AnimalTypeRepository) CreateAnimalType(typeName string) (int, error) {
	query := fmt.Sprintf(`insert into %s(type) values ($1) returning id`, animalTypeTable)

	var typeID int
	row := r.db.QueryRow(query, typeName)

	if err := row.Scan(&typeID); err != nil {
		return 0, errors.ErrAlreadyExist
	}
	return typeID, nil
}

func (r *AnimalTypeRepository) UpdateAnimalType(typeID int, typeName string) (*domain.AnimalType, error) {
	query := fmt.Sprintf(`
	update %s 
		set type = $1
	where
		id = $2
	`, animalTypeTable)

	res, err := r.db.Exec(query, typeName, typeID)

	if err != nil {
		return nil, errors.ErrAlreadyExist
	}

	if affected, err := res.RowsAffected(); affected == 0 || err != nil {
		return nil, errors.ErrNotFound
	}

	return &domain.AnimalType{ID: typeID, Type: typeName}, nil
}

func (r *AnimalTypeRepository) DeleteAnimalType(typeID int) error {
	query := fmt.Sprintf(`
	delete from %s
	where id = $1
	`, animalTypeTable)

	res, err := r.db.Exec(query, typeID)

	if err != nil {
		return err
	}

	if affected, err := res.RowsAffected(); affected == 0 || err != nil {
		return errors.ErrNotFound
	}

	return nil
}
