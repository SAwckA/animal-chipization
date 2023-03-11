package usecase

import (
	"animal-chipization/internal/domain"
)

type visitedLocationRepository interface {
	VisitedLocation(id int) (*domain.VisitedLocation, error)
	Search(animalID int, params *domain.SearchVisitedLocation) ([]domain.VisitedLocation, error)
	Save(animalID int, location *domain.VisitedLocation) (int, error)
	Update(visitedLocation *domain.VisitedLocation) error
	Delete(id int) error
}

type VisitedLocationUsecase struct {
	repo         visitedLocationRepository
	animalRepo   animalRepository
	locationRepo locationRepository
}

func NewVisitedLocationUsecase(repo visitedLocationRepository, locationRepo locationRepository, animalRepo animalRepository) *VisitedLocationUsecase {
	return &VisitedLocationUsecase{
		repo:         repo,
		locationRepo: locationRepo,
		animalRepo:   animalRepo,
	}
}

func (u *VisitedLocationUsecase) Create(animalID, pointID int) (*domain.VisitedLocation, error) {
	animal, err := u.animalRepo.Animal(animalID)
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

	_, err = u.locationRepo.Location(pointID)
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

	visitedLocation := domain.NewVisitedLocation(pointID)

	locationID, err := u.repo.Save(animalID, visitedLocation)
	if err != nil {
		return nil, err
	}

	visitedLocation.ID = locationID

	return visitedLocation, nil
}

func (u *VisitedLocationUsecase) Search(animalID int, params *domain.SearchVisitedLocation) ([]domain.VisitedLocation, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	_, err := u.animalRepo.Animal(animalID)
	if err != nil {
		return nil, err
	}

	return u.repo.Search(animalID, params)
}

func (u *VisitedLocationUsecase) Update(animalID int, location *domain.UpdateVisitedLocationDTO) (*domain.VisitedLocation, error) {
	animal, err := u.animalRepo.Animal(animalID)
	if err != nil {
		return nil, err
	}

	_, err = u.locationRepo.Location(location.LocationPointID)
	if err != nil {
		return nil, err
	}

	visitedLocation, err := u.repo.VisitedLocation(location.VisitedLocationPointID)
	if err != nil {
		return nil, err
	}

	// Обновление точки на такую же точку
	if visitedLocation.LocationPointID == location.LocationPointID {
		return nil, &domain.ApplicationError{
			OriginalError: nil,
			SimplifiedErr: domain.ErrInvalidInput,
			Description:   "update to same location point",
		}
	}

	if len(animal.VisitedLocations) > 0 {
		pos, err := animal.FindVisitedLocationPos(location.VisitedLocationPointID)
		if err != nil {
			return nil, err
		}

		// Обновление первой посещенной точки на точку чипирования
		if pos == 0 {
			if animal.ChippingLocationId == location.LocationPointID {
				return nil, &domain.ApplicationError{
					OriginalError: nil,
					SimplifiedErr: domain.ErrInvalidInput,
					Description:   "first visited location equal to chipping location",
				}
			}
		}

		// Обновление точки локации на точку, совпадающую со следующей и/или с предыдущей точками
		if pos > 0 && len(animal.VisitedLocations) > 1 {
			if animal.VisitedLocations[pos-1].LocationPointID == location.LocationPointID {
				return nil, &domain.ApplicationError{
					OriginalError: nil,
					SimplifiedErr: domain.ErrInvalidInput,
					Description:   "last visited location equal to new",
				}
			}
		}
		if pos < (len(animal.VisitedLocations) - 1) {
			if animal.VisitedLocations[pos+1].LocationPointID == location.LocationPointID {
				return nil, &domain.ApplicationError{
					OriginalError: nil,
					SimplifiedErr: domain.ErrInvalidInput,
					Description:   "next visited location equal to new",
				}
			}
		}
	}

	visitedLocation.LocationPointID = location.LocationPointID

	err = u.repo.Update(visitedLocation)

	return visitedLocation, err
}

func (u *VisitedLocationUsecase) Delete(animalID int, locationID int) error {
	// Животное с animalId не найдено
	animal, err := u.animalRepo.Animal(animalID)
	if err != nil {
		return err
	}

	// Объект с информацией о посещенной точке локации с visitedPointId не найден.
	_, err = u.repo.VisitedLocation(locationID)
	if err != nil {
		return err
	}

	// У животного нет объекта с информацией о посещенной точке локации с visitedPointId
	_, err = animal.FindVisitedLocationPos(locationID)
	if err != nil {
		return err
	}

	if len(animal.VisitedLocations) >= 2 {
		pos, err := animal.FindVisitedLocationPos(locationID)
		if err != nil {
			return err
		}

		if pos == 0 {
			if animal.VisitedLocations[pos+1].LocationPointID == animal.ChippingLocationId {
				_ = u.repo.Delete(animal.VisitedLocations[pos+1].ID)
			}
		}
	}

	return u.repo.Delete(locationID)
}
