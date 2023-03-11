package domain

type Location struct {
	ID        int      `json:"id"`
	Latitude  *float64 `json:"latitude" db:"latitude" binding:"required,lte=90,gte=-90"`
	Longitude *float64 `json:"longitude" db:"longitude" binding:"required,lte=180,gte=-180"`
}

func (l *Location) Map() map[string]interface{} {
	return map[string]interface{}{
		"id":        l.ID,
		"latitude":  l.Latitude,
		"longitude": l.Longitude,
	}
}
