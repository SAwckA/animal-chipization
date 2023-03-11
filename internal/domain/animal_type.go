package domain

type AnimalType struct {
	ID   int    `json:"id" db:"id"`
	Type string `json:"type" db:"type"`
}

func (d *AnimalType) Map() map[string]interface{} {
	return map[string]interface{}{
		"id":   d.ID,
		"type": d.Type,
	}
}

type AnimalTypeCreate struct {
	Type string `json:"type" binding:"required,exclude_whitespace"`
}
