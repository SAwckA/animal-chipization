package domain

import (
	"time"
)

const (
	VisitedLocationsDefaultSize       = 10
	VisitedLocationsSearchDefaultFrom = 0
)

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

func (v *VisitedLocation) Map() map[string]interface{} {
	return map[string]interface{}{
		"id":                           v.ID,
		"dateTimeOfVisitLocationPoint": v.DateTime.Format(time.RFC3339),
		"locationPointId":              v.LocationPointID,
	}
}

type UpdateVisitedLocationDTO struct {
	VisitedLocationPointID int `json:"visitedLocationPointId" binding:"gt=0,required"`
	LocationPointID        int `json:"locationPointId" binding:"gt=0,required"`
}

type SearchVisitedLocation struct {
	StartDateTime *time.Time `form:"startDateTime" time_format:"2006-01-02T15:04:05Z07:00"`
	EndDateTime   *time.Time `form:"endDateTime" time_format:"2006-01-02T15:04:05Z07:00"`
	From          *int       `form:"from"`
	Size          *int       `form:"size"`
}

func (s *SearchVisitedLocation) Validate() error {
	err := &ApplicationError{
		OriginalError: nil,
		SimplifiedErr: ErrInvalidInput,
		Description:   "validation error",
	}
	var defaultFrom, defaultSize = VisitedLocationsSearchDefaultFrom, VisitedLocationsDefaultSize

	if s.From == nil {
		s.From = &defaultFrom
	}
	if s.Size == nil {
		s.Size = &defaultSize
	}

	if *s.From < 0 || *s.Size <= 0 {
		return err
	}

	return nil
}
