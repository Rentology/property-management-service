package models

type PropertyAvailability struct {
	Id          int64  `json:"id"`
	PropertyId  int64  `json:"propertyId"`
	Date        string `json:"date" validate:"datetime=2006-01-02"`
	IsAvailable bool   `json:"isAvailable"`
}
