package usecase

import "animal-chipization/internal/domain"

type locationRepository interface {
	Location(id int) (*domain.Location, error)
	Create(lat, lon float64) (int, error)
	Update(location *domain.Location) error
	Delete(id int) error
}

type LocationUsecase struct {
	repo locationRepository
}

func NewLocationUsecase(repo locationRepository) *LocationUsecase {
	return &LocationUsecase{repo: repo}
}

func (u *LocationUsecase) Location(id int) (*domain.Location, error) {
	return u.repo.Location(id)
}

func (u *LocationUsecase) Create(lat, lon float64) (*domain.Location, error) {
	locationID, err := u.repo.Create(lat, lon)
	if err != nil {
		return nil, err
	}

	return &domain.Location{
		ID:        locationID,
		Latitude:  &lat,
		Longitude: &lon,
	}, nil
}

func (u *LocationUsecase) Update(id int, location *domain.Location) (*domain.Location, error) {
	location.ID = id
	return location, u.repo.Update(location)
}

func (u *LocationUsecase) Delete(id int) error {
	return u.repo.Delete(id)
}
