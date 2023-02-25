package domain

import (
	"errors"
	"time"
)

var ErrInvalidParams = errors.New("invalid params")
var ErrDeadAnimal = errors.New("animal is dead")
var ErrLocationPointEqualChippingLocation = errors.New("attempt to add a location point equal to the chipping point")
var ErrAlreadyLocated = errors.New("attempt to add a location point where the animal is already located")
var ErrEqualNewVisitLocation = errors.New("cant update visit location point to same location point")

type VisitedLocation struct {
	ID              int       `json:"id"`
	DateTime        time.Time `json:"date_time_of_visited_location_point"`
	LocationPointID int       `json:"location_id"`
	AnimalID        int       `json:"animal_id"`
}

func NewVisitedLocation(pointID int) *VisitedLocation {
	return &VisitedLocation{
		DateTime:        time.Now(),
		LocationPointID: pointID,
	}
}

func (v *VisitedLocation) Response() map[string]interface{} {
	return map[string]interface{}{
		"id":                           v.ID,
		"dateTimeOfVisitLocationPoint": v.DateTime.Format(time.RFC3339),
		"locationPointId":              v.LocationPointID,
	}
}

type UpdateVisitedLocationDTO struct {
	VisitedLocationPointID *int `json:"visitedLocationPointId"`
	LocationPointID        *int `json:"locationPointId"`
}

func (u *UpdateVisitedLocationDTO) Validate() error {
	err := &ApplicationError{
		OriginalError: nil,
		SimplifiedErr: ErrInvalidInput,
		Description:   "invalid visited location params",
	}

	switch {
	case u.VisitedLocationPointID == nil || *u.VisitedLocationPointID <= 0:
		return err

	case u.LocationPointID == nil || *u.LocationPointID <= 0:
		return err

	default:
		return nil
	}
}

type SearchVisitedLocationDTO struct {
	StartDateTime *time.Time `form:"startDateTime" time_format:"2006-01-02T15:04:05Z07:00"`
	EndDateTime   *time.Time `form:"endDateTime" time_format:"2006-01-02T15:04:05Z07:00"`
	From          *int       `form:"from"`
	Size          *int       `form:"size"`
}

func (s *SearchVisitedLocationDTO) Validate() error {
	err := &ApplicationError{
		OriginalError: nil,
		SimplifiedErr: ErrInvalidInput,
		Description:   "validation error",
	}
	var defaultFrom, defaultSize = 0, 10

	if s.From == nil {
		s.From = &defaultFrom
	}
	if s.Size == nil {
		s.Size = &defaultSize
	}

	switch {
	case *s.From < 0:
		return err
	case *s.Size <= 0:
		return err

	// TODO: не в формате ISO-8601

	default:
		return nil
	}

}
