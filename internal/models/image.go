package models

type Image struct {
	Id         int64  `json:"id"`
	PropertyId int64  `json:"propertyId" db:"property_id"`
	ImageUrl   string `json:"imageUrl" db:"image_url"`
}
