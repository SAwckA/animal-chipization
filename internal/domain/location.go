package domain

import "errors"

// Точка локации не найдена по id
var ErrLocationNotFoundByID = errors.New("location not found by id")

type Location struct {
	ID        int     `json:"id" db:"id"`
	Latitude  float64 `json:"latitude" db:"latitude"`
	Longitude float64 `json:"longitude" db:"longitude"`
}
