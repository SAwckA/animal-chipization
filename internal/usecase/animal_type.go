package usecase

import "animal-chipization/internal/domain"

type animalTypeRepository interface {
	GetAnimalType(typeID int) *domain.AnimalType
	CreateAnimalType(typeName string) (int, error)
	UpdateAnimalType(typeID int, typeName string) (*domain.AnimalType, error)
	DeleteAnimalType(typeID int) error
}

type AnimalTypeUsecase struct {
	repo animalTypeRepository
}

func NewAnimalTypeUsecase(repo animalTypeRepository) *AnimalTypeUsecase {
	return &AnimalTypeUsecase{repo: repo}
}

func (u *AnimalTypeUsecase) GetType(typeID int) *domain.AnimalType {
	return u.repo.GetAnimalType(typeID)
}

func (u *AnimalTypeUsecase) CreateType(typeName string) (*domain.AnimalType, error) {
	var animalType domain.AnimalType
	typeID, err := u.repo.CreateAnimalType(typeName)

	if err != nil {
		return nil, err
	}

	animalType.ID = typeID
	animalType.Type = typeName

	return &animalType, nil
}

func (u *AnimalTypeUsecase) UpdateType(typeID int, typeName string) (*domain.AnimalType, error) {
	return u.repo.UpdateAnimalType(typeID, typeName)
}

func (u *AnimalTypeUsecase) DeleteType(typeID int) error {
	return u.repo.DeleteAnimalType(typeID)
}
