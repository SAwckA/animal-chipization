package domain

import "errors"

// Тип животного не найден
var ErrAnimalTypeNotFound = errors.New("animal type not found")

type AnimalType struct {
	ID   int    `json:"id" db:"id"`
	Type string `json:"type" db:"type"`
}
