package usecase

import (
	"animal-chipization/internal/domain"
	"time"
)

type animalRepository interface {
	GetAnimal(animalID int) (*domain.Animal, error)
	SearchAnimal(params *domain.AnimalSearchParams) (*[]domain.Animal, error)
	CreateAnimal(params *domain.Animal) (int, error)
	UpdateAnimal(animal domain.Animal) error
	DeleteAnimal(animalID int) error

	AttachTypeAnimal(animalID, typeID int) error
	EditAnimalType(animalID, oldTypeID, newTypeID int) error
	DeleteAnimalType(animalID, typeID int) error
}

type AnimalUsecase struct {
	repo     animalRepository
	typeRepo animalTypeRepository
}

func NewAnimalUsecase(repo animalRepository, typeRepo animalTypeRepository) *AnimalUsecase {
	return &AnimalUsecase{repo: repo, typeRepo: typeRepo}
}

func (u *AnimalUsecase) GetAnimal(animalID int) (*domain.Animal, error) {
	return u.repo.GetAnimal(animalID)
}

func (u *AnimalUsecase) SearchAnimal(params *domain.AnimalSearchParams) (*[]domain.Animal, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}
	return u.repo.SearchAnimal(params)
}

func (u *AnimalUsecase) CreateAnimal(params domain.AnimalCreateParams) (*domain.Animal, error) {

	newAnimal, err := domain.NewAnimal(params)
	if err != nil {
		return nil, &domain.ApplicationError{
			OriginalError: err,
			SimplifiedErr: domain.ErrInvalidInput,
		}
	}

	id, err := u.repo.CreateAnimal(newAnimal)
	newAnimal.ID = id
	return newAnimal, err
}
func (u *AnimalUsecase) UpdateAnimal(animalID int, params domain.AnimalUpdateParams) (*domain.Animal, error) {

	animal, err := u.repo.GetAnimal(animalID)

	if err != nil {
		return nil, err
	}

	if err := params.Validate(); err != nil {
		return nil, err
	}

	animal.Lenght = *params.Lenght
	animal.Weight = *params.Weight
	animal.Height = *params.Height
	animal.Gender = *params.Gender

	if animal.LifeStatus == "DEAD" && *params.LifeStatus == "ALIVE" {
		return nil, &domain.ApplicationError{
			OriginalError: nil,
			SimplifiedErr: domain.ErrInvalidInput,
			Description:   "updating dead animal to alive",
		}
	}

	animal.ChipperID = *params.ChipperID

	if len(animal.VisitedLocations) > 0 {
		if animal.VisitedLocations[0].LocationPointID == *params.ChippingLocationID {
			return nil, &domain.ApplicationError{
				OriginalError: nil,
				SimplifiedErr: domain.ErrInvalidInput,
				Description:   "chipping location id equal first visited location",
			}
		}
	}
	animal.ChippingLocationId = *params.ChippingLocationID

	if *params.LifeStatus == "DEAD" {
		if animal.DeathDateTime == nil {
			deathTime := time.Now()
			animal.DeathDateTime = &deathTime
		}
		animal.LifeStatus = "DEAD"
	}

	err = u.repo.UpdateAnimal(*animal)

	return animal, err

}
func (u *AnimalUsecase) DeleteAnimal(animalID int) error {

	animal, err := u.repo.GetAnimal(animalID)

	if err != nil {
		return err
	}

	if len(animal.VisitedLocations) > 0 {
		return &domain.ApplicationError{
			OriginalError: nil,
			SimplifiedErr: domain.ErrInvalidInput,
			Description:   "animal leaved chipping location",
		}
	}

	return u.repo.DeleteAnimal(animal.ID)
}

func (u *AnimalUsecase) AddAnimalType(animalID, typeID int) (*domain.Animal, error) {

	animal, err := u.repo.GetAnimal(animalID)

	if err != nil {
		return nil, err
	}

	_, err = u.typeRepo.Get(typeID)
	if err != nil {
		return nil, err
	}

	if duplicate := animal.AnimalTypesContains(typeID); duplicate {
		return nil, &domain.ApplicationError{
			OriginalError: nil,
			SimplifiedErr: domain.ErrAlreadyExist,
			Description:   "animal already has this type",
		}
	}

	if err := u.repo.AttachTypeAnimal(animalID, typeID); err != nil {
		return nil, err
	}

	animal.AnimalTypes = append(animal.AnimalTypes, typeID)

	return animal, nil
}

func (u *AnimalUsecase) EditAnimalType(animalID int, params domain.AnimalEditTypeParams) (*domain.Animal, error) {

	if err := params.Validate(); err != nil {
		return nil, err
	}
	animal, err := u.repo.GetAnimal(animalID)
	if err != nil {
		return nil, err
	}

	if contains := animal.AnimalTypesContains(*params.NewTypeID); contains {
		return nil, &domain.ApplicationError{
			OriginalError: nil,
			SimplifiedErr: domain.ErrAlreadyExist,
			Description:   "animal already has this type",
		}
	}

	if contains := animal.AnimalTypesContains(*params.OldTypeID); !contains {
		return nil, &domain.ApplicationError{
			OriginalError: nil,
			SimplifiedErr: domain.ErrNotFound,
			Description:   "animal already has this type",
		}
	}

	_, err = u.typeRepo.Get(*params.NewTypeID)
	if err != nil {
		return nil, err
	}

	_, err = u.typeRepo.Get(*params.OldTypeID)
	if err != nil {
		return nil, err
	}

	if err := u.repo.EditAnimalType(animal.ID, *params.OldTypeID, *params.NewTypeID); err != nil {
		return nil, err
	}

	animal.ReplaceAnimalType(*params.OldTypeID, *params.NewTypeID)

	return animal, nil
}

func (u *AnimalUsecase) DeleteAnimalType(animalID, typeID int) (*domain.Animal, error) {

	animal, err := u.repo.GetAnimal(animalID)
	if err != nil {
		return nil, err
	}

	if len(animal.AnimalTypes) <= 1 {
		return nil, &domain.ApplicationError{
			OriginalError: nil,
			SimplifiedErr: domain.ErrInvalidInput,
			Description:   "animal have no types after deletion",
		}
	}

	if contains := animal.AnimalTypesContains(typeID); !contains {
		return nil, &domain.ApplicationError{
			OriginalError: nil,
			SimplifiedErr: domain.ErrNotFound,
			Description:   "animal doesnt have type with given type id",
		}
	}

	animal.RemoveAnimalType(typeID)

	err = u.repo.DeleteAnimalType(animalID, typeID)

	return animal, err
}
