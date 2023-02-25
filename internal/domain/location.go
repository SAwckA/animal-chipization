package domain

// Точка локации не найдена по id
// var ErrLocationNotFoundByID = errors.New("location not found by id")

type Location struct {
	ID        int      `json:"id"`
	Latitude  *float64 `json:"latitude" db:"latitude" binding:"required,lte=90,gte=-90"`
	Longitude *float64 `json:"longitude" db:"longitude" binding:"required,lte=180,gte=-180"`
}
