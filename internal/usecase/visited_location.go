package usecase

import (
	"animal-chipization/internal/domain"
)

type visitedLocationRepository interface {
	Get(id int) *domain.VisitedLocation
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

func (u *VisitedLocationUsecase) Create(animalID, pointID int) (*domain.VisitedLocation, error) {
	animal := u.animalRepo.GetAnimal(animalID)
	if animal == nil {
		return nil, domain.ErrAnimalNotFoundByID
	}

	if animal.LifeStatus == "DEAD" {
		return nil, domain.ErrDeadAnimal
	}

	// Животное находится в точке чипирования и никуда не перемещалось, попытка добавить точку локации, равную точке чипирования.
	if animal.ChippingLocationId == pointID && len(animal.VisitedLocations) == 0 {
		return nil, domain.ErrLocationPointEqualChippingLocation
	}

	point := u.locationRepo.GetLocation(pointID)
	if point == nil {
		return nil, domain.ErrLocationNotFoundByID
	}

	// Попытка добавить точку локации, в которой уже находится животное

	if len(animal.VisitedLocations) > 0 {
		if animal.VisitedLocations[len(animal.VisitedLocations)-1].LocationPointID == pointID {
			return nil, domain.ErrAlreadyLocated
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

func (u *VisitedLocationUsecase) Update(animalID int, newLocation domain.UpdateVisitedLocationDTO) (*domain.VisitedLocation, error) {
	if err := newLocation.Validate(); err != nil {
		return nil, err
	}

	animal := u.animalRepo.GetAnimal(animalID)
	if animal == nil {
		return nil, domain.ErrAnimalNotFoundByID
	}

	point := u.locationRepo.GetLocation(*newLocation.LocationPointID)
	if point == nil {
		return nil, domain.ErrLocationNotFoundByID
	}

	visitedLocation := u.repo.Get(*newLocation.VisitedLocationPointID)
	if visitedLocation == nil {
		return nil, domain.ErrLocationNotFoundByID
	}

	// Обновление точки на такую же точку
	if visitedLocation.LocationPointID == *newLocation.LocationPointID {
		return nil, domain.ErrEqualNewVisitLocation
	}

	if len(animal.VisitedLocations) > 0 {
		pos, err := animal.FindVisitedLocaionPos(*newLocation.VisitedLocationPointID)
		if err != nil {
			return nil, err
		}

		// Обновление первой посещенной точки на точку чипирования
		if pos == 0 {
			if animal.ChippingLocationId == *newLocation.LocationPointID {
				return nil, domain.ErrLocationPointEqualChippingLocation
			}
		}

		// Обновление точки локации на точку, совпадающую со следующей и/или с предыдущей точками
		if animal.VisitedLocations[pos-1].LocationPointID == *newLocation.LocationPointID {
			return nil, domain.ErrAlreadyLocated
		}
		if pos < (len(animal.VisitedLocations)) {
			if animal.VisitedLocations[pos+1].LocationPointID == *newLocation.LocationPointID {
				return nil, domain.ErrAlreadyLocated
			}
		}
	}

	visitedLocation.LocationPointID = *newLocation.LocationPointID

	err := u.repo.Update(*visitedLocation)

	return visitedLocation, err
}

func (u *VisitedLocationUsecase) Delete(animalID int, locatoinID int) error {
	// Животное с animalId не найдено
	animal := u.animalRepo.GetAnimal(animalID)
	if animal == nil {
		return domain.ErrAnimalNotFoundByID
	}

	// Объект с информацией о посещенной точке локации с visitedPointId не найден.
	visitedLocation := u.repo.Get(locatoinID)
	if visitedLocation == nil {
		return domain.ErrLocationNotFoundByID
	}
	// У животного нет объекта с информацией о посещенной точке локации с visitedPointId
	_, err := animal.FindVisitedLocaionPos(locatoinID)

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
		return nil, domain.ErrInvalidParams
	}
	animal := u.animalRepo.GetAnimal(animalID)
	if animal == nil {
		return nil, domain.ErrAnimalNotFoundByID
	}

	return u.repo.Search(animalID, params)
}
