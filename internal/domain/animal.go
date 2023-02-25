package domain

import (
	"errors"
	"time"
)

// Животное не найдено по id
var ErrAnimalNotFoundByID = errors.New("animal not found by id")

// Встречается только в AnimalCreateParams
// Ошибка дублирования типа животного в массиве
var ErrAnimalTypeParamsDuplicate = errors.New("animal type list has duplicates")

// Отсутствие типа в списке типов животного
var ErrMissingAnimalType = errors.New("missing animal type in list of animal types")

// Универсальная ошибка валидации параметров создания животного
var ErrAnimalCreateParamsInvalid = errors.New("invalid create animal params")

// Универсальная ошибка валидации параметров обновления животного
var ErrAnimalUpdateParamsInvalid = errors.New("invalid update animal params")

// Универсальная ошибка валидации параметров изменения типа животного
var ErrAnimalEditTypeParamsInvalid = errors.New("invalid edit animal type params")

// Отсутсвие типов в животного
var ErrAnimalTypeListEmpty = errors.New("empty animal type list")

var ErrAnimalVisitLocationNotFound = errors.New("animal dont visited this location")

type Animal struct {
	ID                 int               `json:"id" db:"id"`
	AnimalTypes        []int             `json:"animalTypes"`
	Lenght             float32           `json:"lenght" db:"lenght"`
	Weight             float32           `json:"weight" db:"weight"`
	Height             float32           `json:"height" db:"height"`
	Gender             string            `json:"gender" db:"gender"`
	LifeStatus         string            `json:"lifeStatus" db:"lifestatus"`
	ChippingDateTime   time.Time         `json:"chippinDateTime" db:"chippinglocationid"`
	ChipperID          int               `json:"chippedId" db:"chipperid"`
	ChippingLocationId int               `json:"chippingLocationId" db:"chippinglocationid"`
	VisitedLocations   []VisitedLocation `json:"visitedLocations" db:"animaltypes"`
	DeathDateTime      *time.Time        `json:"deathDateTime" db:"deathdatetime"`
}

func (a *Animal) FindVisitedLocaionPos(id int) (int, error) {
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

func (a *Animal) Response() map[string]interface{} {
	var visitedLocations = make([]int, 0)

	for _, v := range a.VisitedLocations {
		visitedLocations = append(visitedLocations, v.ID)
	}

	resp := map[string]interface{}{
		"id":                 a.ID,
		"animalTypes":        a.AnimalTypes,
		"length":             a.Lenght,
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

type AnimalCreateParams struct {
	AnimalTypes        *[]int   `json:"animalTypes"`
	Lenght             *float32 `json:"length"`
	Weight             *float32 `json:"weight"`
	Height             *float32 `json:"height"`
	Gender             *string  `json:"gender"`
	ChipperID          *int     `json:"chipperId"`
	ChippingLocationID *int     `json:"chippingLocationId"`
}

func validateList(l *[]int) bool {
	for _, v := range *l {
		if v <= 0 {
			return false
		}
	}
	return true
}

func (p *AnimalCreateParams) Validate() error {
	var err = ErrAnimalCreateParamsInvalid

	switch {

	case p.AnimalTypes == nil || len(*p.AnimalTypes) <= 0 || !validateList(p.AnimalTypes):
		return err
	case p.Weight == nil || *p.Weight <= 0:
		return err
	case p.Lenght == nil || *p.Lenght <= 0:
		return err
	case p.Height == nil || *p.Height <= 0:
		return err
	case p.Gender == nil || (*p.Gender != "MALE" && *p.Gender != "FEMALE" && *p.Gender != "OTHER"):
		return err
	case p.ChipperID == nil || *p.ChipperID <= 0:
		return err
	case p.ChippingLocationID == nil || *p.ChippingLocationID <= 0:
		return err

	default:
		return nil
	}
}

type AnimalUpdateParams struct {
	Lenght             *float32 `json:"length"`
	Weight             *float32 `json:"weight"`
	Height             *float32 `json:"height"`
	Gender             *string  `json:"gender"`
	LifeStatus         *string  `json:"lifeStatus"`
	ChipperID          *int     `json:"chipperId"`
	ChippingLocationID *int     `json:"chippingLocationId"`
}

func (p *AnimalUpdateParams) Validate() error {
	err := &ApplicationError{
		OriginalError: ErrInvalidInput,
		SimplifiedErr: ErrInvalidInput,
	}

	switch {

	case p.Weight == nil || *p.Weight <= 0:
		err.Description = "invalid weight"
		return err
	case p.Lenght == nil || *p.Lenght <= 0:
		err.Description = "invalid lenght"
		return err
	case p.Height == nil || *p.Height <= 0:
		err.Description = "invalid height"
		return err
	case p.Gender == nil || (*p.Gender != "MALE" && *p.Gender != "FEMALE" && *p.Gender != "OTHER"):
		err.Description = "invalid gender"
		return err
	case p.LifeStatus == nil || (*p.LifeStatus != "ALIVE" && *p.LifeStatus != "DEAD"):
		err.Description = "invalid lifestatus"
		return err
	case p.ChipperID == nil || *p.ChipperID <= 0:
		err.Description = "invalid chipperid"
		return err
	case p.ChippingLocationID == nil || *p.ChippingLocationID <= 0:
		err.Description = "invalid chippinglocationid"
		return err

	default:
		return nil
	}
}

type AnimalEditTypeParams struct {
	OldTypeID *int `json:"oldTypeId"`
	NewTypeID *int `json:"newTypeId"`
}

func (p *AnimalEditTypeParams) Validate() error {
	err := &ApplicationError{
		OriginalError: ErrInvalidParams,
		SimplifiedErr: ErrInvalidInput,
	}

	switch {
	case p.OldTypeID == nil || *p.OldTypeID <= 0:
		err.Description = "invalid OldtypeID"
		return err

	case p.NewTypeID == nil || *p.NewTypeID <= 0:
		err.Description = "invalid NewTypeID"
		return err

	default:
		return nil
	}
}

func NewAnimal(params AnimalCreateParams) (*Animal, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	return &Animal{
		AnimalTypes:        *params.AnimalTypes,
		Lenght:             *params.Lenght,
		Weight:             *params.Weight,
		Height:             *params.Height,
		Gender:             *params.Gender,
		ChipperID:          *params.ChipperID,
		ChippingLocationId: *params.ChippingLocationID,
		LifeStatus:         "ALIVE",
		ChippingDateTime:   time.Now(),
		DeathDateTime:      nil,
	}, nil
}

type AnimalSearchParams struct {
	StartDateTime *time.Time
	EndDateTime   *time.Time

	ChipperID         *int
	ChippedLocationID *int

	LifeStatus *string
	Gender     *string

	From int
	Size int
}
