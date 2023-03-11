package psql

import (
	"animal-chipization/internal/domain"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

const animalTypeTable = "public.animal_type"
const uniqueTypeConstraint = "unique_type"
const animalTypeFkey = "animal_types_list_type_id_fkey"

type AnimalTypeRepository struct {
	db *sqlx.DB
}

func NewAnimalTypeRepository(db *sqlx.DB) *AnimalTypeRepository {
	return &AnimalTypeRepository{db: db}
}

func (r *AnimalTypeRepository) AnimalType(id int) (*domain.AnimalType, error) {
	query := fmt.Sprintf(`select id, type from %s where id = $1`, animalTypeTable)

	var animalType domain.AnimalType
	if err := r.db.QueryRow(query, id).Scan(&animalType.ID, &animalType.Type); err != nil {
		return nil, &domain.ApplicationError{
			OriginalError: err,
			SimplifiedErr: domain.ErrNotFound,
			Description:   "animal type not found",
		}
	}

	return &animalType, nil
}

func (r *AnimalTypeRepository) Create(typeName string) (int, error) {
	query := fmt.Sprintf(`insert into %s(type) values ($1) returning id`, animalTypeTable)

	var typeID int
	if err := r.db.QueryRow(query, typeName).Scan(&typeID); err != nil {
		if strings.Contains(err.Error(), uniqueTypeConstraint) {
			return 0, &domain.ApplicationError{
				OriginalError: err,
				SimplifiedErr: domain.ErrAlreadyExist,
				Description:   "animal type with this type already exist",
			}
		}
		return 0, &domain.ApplicationError{
			OriginalError: err,
			SimplifiedErr: domain.ErrUnknown,
			Description:   "unknown error in (r *AnimalTypeRepository) Create",
		}
	}
	return typeID, nil
}

func (r *AnimalTypeRepository) Update(id int, typeName string) error {
	query := fmt.Sprintf(`
	update %s 
		set type = $1
	where
		id = $2
	`, animalTypeTable)

	res, err := r.db.Exec(query, typeName, id)
	if err != nil {
		if strings.Contains(err.Error(), uniqueTypeConstraint) {
			return &domain.ApplicationError{
				OriginalError: err,
				SimplifiedErr: domain.ErrAlreadyExist,
				Description:   "animal type with this type already exist",
			}
		}
		return &domain.ApplicationError{
			OriginalError: err,
			SimplifiedErr: domain.ErrUnknown,
			Description:   "unknown error in (r *AnimalTypeRepository) Create",
		}
	}

	if aff, err := res.RowsAffected(); aff == 0 || err != nil {
		return &domain.ApplicationError{
			OriginalError: err,
			SimplifiedErr: domain.ErrNotFound,
			Description:   "animal type not found by id during update animal type",
		}
	}

	return nil
}

func (r *AnimalTypeRepository) Delete(id int) error {
	query := fmt.Sprintf(`
	delete from %s
	where id = $1
	`, animalTypeTable)

	res, err := r.db.Exec(query, id)
	if err != nil {
		if strings.Contains(err.Error(), animalTypeFkey) {
			return &domain.ApplicationError{
				OriginalError: err,
				SimplifiedErr: domain.ErrLinked,
				Description:   "animal type linked with animal",
			}
		}

		return &domain.ApplicationError{
			OriginalError: err,
			SimplifiedErr: domain.ErrUnknown,
			Description:   "unknown error in (r *AnimalTypeRepository) Create",
		}
	}

	if aff, err := res.RowsAffected(); aff == 0 || err != nil {
		return &domain.ApplicationError{
			OriginalError: err,
			SimplifiedErr: domain.ErrNotFound,
			Description:   "animal type not found by id during delete animal type",
		}
	}

	return nil
}
