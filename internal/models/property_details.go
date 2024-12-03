package models

type PropertyDetails struct {
	PropertyID        int64  `json:"propertyId" db:"property_id"`
	Floor             int    `json:"floor" db:"floor"`
	MaxFloor          int    `json:"max_floor" db:"max_floor"`
	Area              int    `json:"area" db:"area"`
	Rooms             int    `json:"rooms" db:"rooms"`
	HouseCreationYear int    `json:"houseCreationYear" db:"house_creation_year"`
	HouseType         string `json:"houseType" db:"house_type"`
	Description       string `json:"description" db:"description"`
}
