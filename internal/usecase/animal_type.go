package usecase

import "animal-chipization/internal/domain"

type animalTypeRepository interface {
	Get(id int) (*domain.AnimalType, error)
	Create(typeName string) (int, error)
	Update(id int, typeName string) error
	Delete(id int) error
}

type AnimalTypeUsecase struct {
	repo animalTypeRepository
}

func NewAnimalTypeUsecase(repo animalTypeRepository) *AnimalTypeUsecase {
	return &AnimalTypeUsecase{repo: repo}
}

func (u *AnimalTypeUsecase) Get(id int) (*domain.AnimalType, error) {
	return u.repo.Get(id)
}

func (u *AnimalTypeUsecase) Create(typeName string) (*domain.AnimalType, error) {
	var animalType domain.AnimalType

	typeID, err := u.repo.Create(typeName)

	if err != nil {
		return nil, err
	}

	animalType.ID = typeID
	animalType.Type = typeName

	return &animalType, nil
}

func (u *AnimalTypeUsecase) Update(id int, typeName string) (*domain.AnimalType, error) {
	err := u.repo.Update(id, typeName)
	if err != nil {
		return nil, err
	}
	return &domain.AnimalType{ID: id, Type: typeName}, nil
}

func (u *AnimalTypeUsecase) Delete(id int) error {
	return u.repo.Delete(id)
}
