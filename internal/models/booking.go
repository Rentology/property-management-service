package models

type Booking struct {
	Id           int64  `json:"id"`
	PropertyId   int64  `json:"propertyId"`
	UserId       int64  `json:"userId"`
	CheckInDate  string `json:"checkInDate" validate:"datetime=2006-01-02"`
	CheckOutDate string `json:"checkOutDate" validate:"datetime=2006-01-02"`
	TotalPrice   int    `json:"totalPrice"`
	Status       string `json:"status" validate:"oneof= confirmed pending cancelled"`
	CreatedAt    string `json:"createdAt" validate:"datetime=2006-01-02"`
}
