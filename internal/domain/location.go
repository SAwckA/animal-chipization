package domain

type Location struct {
	ID        int     `json:"id" db:"id"`
	Latitude  float64 `json:"latitude" db:"latitude"`
	Longitude float64 `json:"longitude" db:"longitude"`
}
