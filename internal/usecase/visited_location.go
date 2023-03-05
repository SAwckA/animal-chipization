package usecase

import (
	"animal-chipization/internal/domain"
)

type visitedLocationRepository interface {
	Get(id int) (*domain.VisitedLocation, error)
	Save(animalID int, location domain.VisitedLocation) (int, error)
	Update(domain.VisitedLocation) error
	Delete(id int) error
	Search(animalID int, params domain.SearchVisitedLocationDTO) (*[]domain.VisitedLocation, error)
}

type VisitedLocationUsecase struct {
	repo         visitedLocationRepository
	animalRepo   animalRepository
	locationRepo locationRepository
}

func NewVisitedLocationUsecase(repo visitedLocationRepository, locatoinRepo locationRepository, animalRepo animalRepository) *VisitedLocationUsecase {
	return &VisitedLocationUsecase{
		repo:         repo,
		locationRepo: locatoinRepo,
		animalRepo:   animalRepo,
	}
}

// Добавление точки локации, посещенной животным
//
func (u *VisitedLocationUsecase) Create(animalID, pointID int) (*domain.VisitedLocation, error) {
	animal, err := u.animalRepo.GetAnimal(animalID)
	if err != nil {
		return nil, err
	}

	if animal.LifeStatus == "DEAD" {
		return nil, &domain.ApplicationError{
			OriginalError: nil,
			SimplifiedErr: domain.ErrInvalidInput,
			Description:   "dead animal",
		}
	}

	// Животное находится в точке чипирования и никуда не перемещалось, попытка добавить точку локации, равную точке чипирования.
	if animal.ChippingLocationId == pointID && len(animal.VisitedLocations) == 0 {
		return nil, &domain.ApplicationError{
			OriginalError: nil,
			SimplifiedErr: domain.ErrInvalidInput,
			Description:   "location point equal chipping location",
		}
	}

	_, err = u.locationRepo.GetLocation(pointID)
	if err != nil {
		return nil, err
	}

	// Попытка добавить точку локации, в которой уже находится животное

	if len(animal.VisitedLocations) > 0 {
		if animal.VisitedLocations[len(animal.VisitedLocations)-1].LocationPointID == pointID {
			return nil, &domain.ApplicationError{
				OriginalError: nil,
				SimplifiedErr: domain.ErrInvalidInput,
			}
		}
	}

	visitedLocaton := domain.NewVisitedLocation(pointID)

	locationID, err := u.repo.Save(animalID, *visitedLocaton)
	if err != nil {
		return nil, err
	}

	visitedLocaton.ID = locationID

	return visitedLocaton, nil
}

// Изменение точки локации, посещенной животным
//
func (u *VisitedLocationUsecase) Update(animalID int, newLocation domain.UpdateVisitedLocationDTO) (*domain.VisitedLocation, error) {
	if err := newLocation.Validate(); err != nil {
		return nil, err
	}

	animal, err := u.animalRepo.GetAnimal(animalID)
	if err != nil {
		return nil, err
	}

	_, err = u.locationRepo.GetLocation(*newLocation.LocationPointID)
	if err != nil {
		return nil, err
	}

	visitedLocation, err := u.repo.Get(*newLocation.VisitedLocationPointID)
	if err != nil {
		return nil, err
	}

	// Обновление точки на такую же точку
	if visitedLocation.LocationPointID == *newLocation.LocationPointID {
		return nil, &domain.ApplicationError{
			OriginalError: nil,
			SimplifiedErr: domain.ErrInvalidInput,
			Description:   "update to same location point",
		}
	}

	if len(animal.VisitedLocations) > 0 {
		pos, err := animal.FindVisitedLocaionPos(*newLocation.VisitedLocationPointID)
		if err != nil {
			return nil, err
		}

		// Обновление первой посещенной точки на точку чипирования
		if pos == 0 {
			if animal.ChippingLocationId == *newLocation.LocationPointID {
				return nil, &domain.ApplicationError{
					OriginalError: nil,
					SimplifiedErr: domain.ErrInvalidInput,
					Description:   "first visited location equal to chipping location",
				}
			}
		}

		// Обновление точки локации на точку, совпадающую со следующей и/или с предыдущей точками
		if pos > 0 && len(animal.VisitedLocations) > 1 {
			if animal.VisitedLocations[pos-1].LocationPointID == *newLocation.LocationPointID {
				return nil, &domain.ApplicationError{
					OriginalError: nil,
					SimplifiedErr: domain.ErrInvalidInput,
					Description:   "last visited location equal to new",
				}
			}
		}
		if pos < (len(animal.VisitedLocations) - 1) {
			if animal.VisitedLocations[pos+1].LocationPointID == *newLocation.LocationPointID {
				return nil, &domain.ApplicationError{
					OriginalError: nil,
					SimplifiedErr: domain.ErrInvalidInput,
					Description:   "next visited location equal to new",
				}
			}
		}
	}

	visitedLocation.LocationPointID = *newLocation.LocationPointID

	err = u.repo.Update(*visitedLocation)

	return visitedLocation, err
}

// Удаление посещённой точки локации животного
func (u *VisitedLocationUsecase) Delete(animalID int, locatoinID int) error {
	// Животное с animalId не найдено
	animal, err := u.animalRepo.GetAnimal(animalID)
	if err != nil {
		return err
	}

	// Объект с информацией о посещенной точке локации с visitedPointId не найден.
	_, err = u.repo.Get(locatoinID)
	if err != nil {
		return err
	}
	// У животного нет объекта с информацией о посещенной точке локации с visitedPointId
	_, err = animal.FindVisitedLocaionPos(locatoinID)

	if err != nil {
		return err
	}

	//TODO: (Если удаляется первая посещенная точка локации, а вторая точка совпадает с точкой чипирования, то она удаляется автоматически)

	if len(animal.VisitedLocations) >= 2 {
		pos, err := animal.FindVisitedLocaionPos(locatoinID)
		if err != nil {
			return err
		}

		if pos == 0 {
			if animal.VisitedLocations[pos+1].LocationPointID == animal.ChippingLocationId {
				u.repo.Delete(animal.VisitedLocations[pos+1].ID)
			}
		}
	}

	return u.repo.Delete(locatoinID)
}

func (u *VisitedLocationUsecase) Search(animalID int, params domain.SearchVisitedLocationDTO) (*[]domain.VisitedLocation, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}
	_, err := u.animalRepo.GetAnimal(animalID)
	if err != nil {
		return nil, err
	}

	return u.repo.Search(animalID, params)
}
