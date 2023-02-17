package usecase

import (
	"animal-chipization/internal/domain"
	"time"
)

type animalRepository interface {
	GetAnimal(animalID int) *domain.Animal
	SearchAnimal(params *domain.AnimalSearchParams) *[]domain.Animal
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

func (u *AnimalUsecase) GetAnimal(animalID int) *domain.Animal {
	return u.repo.GetAnimal(animalID)
}

func (u *AnimalUsecase) SearchAnimal(params *domain.AnimalSearchParams) *[]domain.Animal {
	return u.repo.SearchAnimal(params)
}

func (u *AnimalUsecase) CreateAnimal(params domain.AnimalCreateParams) (*domain.Animal, error) {

	newAnimal, err := domain.NewAnimal(params)
	if err != nil {
		return nil, err
	}

	id, err := u.repo.CreateAnimal(newAnimal)
	newAnimal.ID = id
	return newAnimal, err
}
func (u *AnimalUsecase) UpdateAnimal(animalID int, params domain.AnimalUpdateParams) (*domain.Animal, error) {

	animal := u.repo.GetAnimal(animalID)

	if animal == nil {
		return nil, domain.ErrAnimalNotFoundByID
	}

	if err := params.Validate(); err != nil {
		return nil, err
	}

	animal.Lenght = *params.Lenght
	animal.Weight = *params.Weight
	animal.Height = *params.Height
	animal.Gender = *params.Gender

	if animal.LifeStatus == "DEAD" && *params.LifeStatus == "ALIVE" {
		return nil, domain.ErrAnimalUpdateParamsInvalid
	}

	animal.ChipperID = *params.ChipperID
	animal.ChippingLocationId = *params.ChippingLocationID

	if *params.LifeStatus == "DEAD" {
		if animal.DeathDateTime == nil {
			deathTime := time.Now()
			animal.DeathDateTime = &deathTime
		}
		animal.LifeStatus = "DEAD"
	}

	err := u.repo.UpdateAnimal(*animal)

	return animal, err

}
func (u *AnimalUsecase) DeleteAnimal(animalID int) error {

	animal := u.repo.GetAnimal(animalID)

	if animal == nil {
		return domain.ErrAnimalNotFoundByID
	}

	// TODO:
	// После пункта 6, нужно добавить следующую проверку:
	//
	// Животное покинуло локацию чипирования, при этом
	// есть другие посещенные точки

	return u.repo.DeleteAnimal(animal.ID)
}

func (u *AnimalUsecase) AddAnimalType(animalID, typeID int) (*domain.Animal, error) {

	animal := u.repo.GetAnimal(animalID)

	if animal == nil {
		return nil, domain.ErrAnimalNotFoundByID
	}

	animalType := u.typeRepo.GetAnimalType(typeID)

	if animalType == nil {
		return nil, domain.ErrAnimalTypeNotFound
	}

	if duplicate := animal.AnimalTypesContains(typeID); duplicate {
		return nil, domain.ErrAnimalTypeParamsDuplicate
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

	animal := u.repo.GetAnimal(animalID)
	if animal == nil {
		return nil, domain.ErrAnimalNotFoundByID
	}

	if contains := animal.AnimalTypesContains(*params.NewTypeID); contains {
		return nil, domain.ErrAnimalTypeParamsDuplicate
	}

	if contains := animal.AnimalTypesContains(*params.OldTypeID); !contains {
		return nil, domain.ErrMissingAnimalType
	}

	newType := u.typeRepo.GetAnimalType(*params.NewTypeID)
	if newType == nil {
		return nil, domain.ErrAnimalTypeNotFound
	}

	oldType := u.typeRepo.GetAnimalType(*params.OldTypeID)
	if oldType == nil {
		return nil, domain.ErrAnimalTypeNotFound
	}

	if err := u.repo.EditAnimalType(animal.ID, *params.OldTypeID, *params.NewTypeID); err != nil {
		return nil, err
	}

	animal.ReplaceAnimalType(*params.OldTypeID, *params.NewTypeID)

	return animal, nil
}

func (u *AnimalUsecase) DeleteAnimalType(animalID, typeID int) (*domain.Animal, error) {

	animal := u.repo.GetAnimal(animalID)
	if animal == nil {
		return nil, domain.ErrAnimalNotFoundByID
	}

	if len(animal.AnimalTypes) <= 1 {
		return nil, domain.ErrAnimalTypeListEmpty
	}

	if contains := animal.AnimalTypesContains(typeID); !contains {
		return nil, domain.ErrAnimalTypeNotFound
	}

	animal.RemoveAnimalType(typeID)

	err := u.repo.DeleteAnimalType(animalID, typeID)

	return animal, err
}
