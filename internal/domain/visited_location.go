package domain

import (
	"errors"
	"time"
)

var ErrInvalidParams = errors.New("invalid params")

type VisitedLocation struct {
	ID              int
	DateTime        time.Time
	LocationPointID int
}

type UpdateVisitedLocationDTO struct {
	VisitedLocationPointID *int `json:"visitedLocationPointId"`
	LocationPointID        *int `json:"locationPointId"`
}

func (u *UpdateVisitedLocationDTO) Validate() error {
	var err = ErrInvalidParams

	switch {
	case u.VisitedLocationPointID == nil || *u.VisitedLocationPointID <= 0:
		return err

	case u.LocationPointID == nil || *u.LocationPointID <= 0:
		return err

	default:
		return nil
	}
}
