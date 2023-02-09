package domain

type AnimalType struct {
	ID   int    `json:"id" db:"id"`
	Type string `json:"type" db:"type"`
}
