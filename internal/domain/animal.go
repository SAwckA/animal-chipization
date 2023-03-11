package domain

import (
	"time"
)

const (
	AnimalSearchDefaultSize = 10
	AnimalSearchDefaultFrom = 0
)

type Animal struct {
	ID                 int
	AnimalTypes        []int
	Length             float32
	Weight             float32
	Height             float32
	Gender             string
	LifeStatus         string
	ChippingDateTime   time.Time
	ChipperID          int
	ChippingLocationId int
	VisitedLocations   []VisitedLocation
	DeathDateTime      *time.Time
}

func (a *Animal) FindVisitedLocationPos(id int) (int, error) {
	for i, v := range a.VisitedLocations {
		if v.ID == id {
			return i, nil
		}
	}
	return 0, &ApplicationError{
		OriginalError: nil,
		SimplifiedErr: ErrNotFound,
		Description:   "animal not visited location",
	}
}

func (a *Animal) ReplaceAnimalType(oldTypeID, newTypeID int) {

	for index, val := range a.AnimalTypes {
		if val == oldTypeID {
			a.AnimalTypes[index] = newTypeID
		}
	}

}

func (a *Animal) AnimalTypesContains(typeID int) bool {
	for _, i := range a.AnimalTypes {
		if i == typeID {
			return true
		}
	}
	return false
}

func (a *Animal) RemoveAnimalType(typeID int) {
	for i := range a.AnimalTypes {
		if a.AnimalTypes[i] == typeID {
			a.AnimalTypes = append(a.AnimalTypes[:i], a.AnimalTypes[i+1:]...)
			return
		}
	}
}

func (a *Animal) Map() map[string]interface{} {
	var visitedLocations = make([]int, 0)

	for _, v := range a.VisitedLocations {
		visitedLocations = append(visitedLocations, v.ID)
	}

	resp := map[string]interface{}{
		"id":                 a.ID,
		"animalTypes":        a.AnimalTypes,
		"length":             a.Length,
		"weight":             a.Weight,
		"height":             a.Height,
		"gender":             a.Gender,
		"lifeStatus":         a.LifeStatus,
		"chippingDateTime":   a.ChippingDateTime.Format(time.RFC3339),
		"chippingLocationId": a.ChippingLocationId,
		"chipperId":          a.ChipperID,
		"visitedLocations":   visitedLocations,
	}

	if a.DeathDateTime == nil {
		resp["deathDateTime"] = nil
		return resp
	}

	resp["deathDateTime"] = a.DeathDateTime.Format(time.RFC3339)
	return resp
}

func NewAnimal(params *AnimalCreateParams) (*Animal, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	return &Animal{
		AnimalTypes:        params.AnimalTypes,
		Length:             params.Length,
		Weight:             params.Weight,
		Height:             params.Height,
		Gender:             params.Gender,
		ChipperID:          params.ChipperID,
		ChippingLocationId: params.ChippingLocationID,
		LifeStatus:         "ALIVE",
		ChippingDateTime:   time.Now(),
		DeathDateTime:      nil,
	}, nil
}

type AnimalCreateParams struct {
	AnimalTypes        []int   `json:"animalTypes"`
	Length             float32 `json:"length" binding:"gt=0,required"`
	Weight             float32 `json:"weight" binding:"gt=0,required"`
	Height             float32 `json:"height" binding:"gt=0,required"`
	Gender             string  `json:"gender" binding:"allowed_strings=MALE;FEMALE;OTHER"`
	ChipperID          int     `json:"chipperId" binding:"gt=0,required"`
	ChippingLocationID int     `json:"chippingLocationId" binding:"gt=0,required"`
}

func (p *AnimalCreateParams) Validate() error {
	err := &ApplicationError{
		OriginalError: nil,
		SimplifiedErr: ErrInvalidInput,
		Description:   "Invalid create animal params",
	}

	if p.AnimalTypes == nil || len(p.AnimalTypes) <= 0 {
		return err
	}
	for _, v := range p.AnimalTypes {
		if v <= 0 {
			return err
		}
	}
	return nil
}

type AnimalUpdateParams struct {
	Length             float32 `json:"length" binding:"gt=0,required"`
	Weight             float32 `json:"weight" binding:"gt=0,required"`
	Height             float32 `json:"height" binding:"gt=0,required"`
	Gender             string  `json:"gender" binding:"allowed_strings=MALE;FEMALE;OTHER"`
	LifeStatus         string  `json:"lifeStatus" binding:"allowed_strings=ALIVE;DEAD"`
	ChipperID          int     `json:"chipperId" binding:"gt=0,required"`
	ChippingLocationID int     `json:"chippingLocationId" binding:"gt=0,required"`
}

type AnimalEditTypeParams struct {
	OldTypeID int `json:"oldTypeId" binding:"gt=0,required"`
	NewTypeID int `json:"newTypeId" binding:"gt=0,required"`
}

type AnimalSearchParams struct {
	StartDateTime     *time.Time `form:"startDateTime" time_format:"2006-01-02T15:04:05Z07:00"`
	EndDateTime       *time.Time `form:"endDateTime" time_format:"2006-01-02T15:04:05Z07:00"`
	ChipperID         *int       `form:"chipperId"`
	ChippedLocationID *int       `form:"chippingLocationId"`
	LifeStatus        *string    `form:"lifeStatus"`
	Gender            *string    `form:"gender"`

	From *int `form:"from"`
	Size *int `form:"size"`
}

func (s *AnimalSearchParams) Validate() error {
	err := &ApplicationError{
		OriginalError: nil,
		SimplifiedErr: ErrInvalidInput,
		Description:   "validation error",
	}
	var defaultFrom, defaultSize = AnimalSearchDefaultFrom, AnimalSearchDefaultSize

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
