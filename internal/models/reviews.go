package models

type Review struct {
	Id         int64  `json:"id"`
	PropertyId int64  `json:"propertyId"`
	UserId     int64  `json:"userId"`
	Rating     int    `json:"rating"`
	Comment    string `json:"comment"`
	CreatedAt  string `json:"createdAt" validate:"datetime=2006-01-02"`
}
