package domain

// var ErrAnimalTypeNotFound = errors.New("animal type not found")                        // Тип животного не найден
// var ErrAnimalTypeLinked = errors.New("animal type linked to animal")                   // Тип животного свзязан с животным
// var ErrAnimalTypeAlreadyExist = errors.New("animal type with this name already exist") // Тип животного с таким type уже существует

type AnimalType struct {
	ID   int    `json:"id" db:"id"`
	Type string `json:"type" db:"type"`
}

func (d *AnimalType) Response() map[string]interface{} {
	return map[string]interface{}{
		"id":   d.ID,
		"type": d.Type,
	}
}

type TypeId struct {
	ID int `uri:"typeId" binding:"required,gt=0"`
}

type AnimalTypeDTO struct {
	Type string `json:"type" binding:"required,exclude_whitespace"`
}
