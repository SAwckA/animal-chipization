package usecase

type visitedLocationRepository interface {
}

type VisitedLocationUsecase struct {
	repo         visitedLocationRepository
	locationRepo locationRepository
}

func NewVisitedLocationUsecase(repo visitedLocationRepository, locatoinRepo locationRepository) *VisitedLocationUsecase {
	return &VisitedLocationUsecase{
		repo:         repo,
		locationRepo: locatoinRepo,
	}
}
