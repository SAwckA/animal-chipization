package usecase

import "animal-chipization/internal/domain"

type locationRepository interface {
	CreateLocation(lat, lon float64) (int, error)
	GetLocation(locationID int) (*domain.Location, error)
	UpdateLocation(location *domain.Location) error
	DeleteLocation(locationID int) error
}

type LocationUsecase struct {
	repo locationRepository
}

func NewLocationUsecase(repo locationRepository) *LocationUsecase {
	return &LocationUsecase{repo: repo}
}

func (u *LocationUsecase) CreateLocation(lat, lon float64) (*domain.Location, error) {

	locationID, err := u.repo.CreateLocation(lat, lon)

	if err != nil {
		return nil, err
	}

	return &domain.Location{
		ID:        locationID,
		Latitude:  &lat,
		Longitude: &lon,
	}, nil
}

func (u *LocationUsecase) GetLocation(locationID int) (*domain.Location, error) {
	return u.repo.GetLocation(locationID)
}

func (u *LocationUsecase) UpdateLocation(locationID int, location *domain.Location) (*domain.Location, error) {
	location.ID = locationID
	return location, u.repo.UpdateLocation(location)
}

func (u *LocationUsecase) DeleteLocation(locationID int) error {
	return u.repo.DeleteLocation(locationID)
}
