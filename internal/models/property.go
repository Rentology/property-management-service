package models

type Property struct {
	ID           int64  `json:"id"`
	OwnerId      int64  `json:"owner_id"`
	Title        string `json:"title"`
	Location     string `json:"location"`
	Price        int    `json:"price"`
	PropertyType string `json:"propertyType" validate:"oneof= house apartment"`
	RentalType   string `json:"rentalType" validate:"oneof = shortTerm longTerm"`
	MaxGuests    int    `json:"maxGuests"`
	CreatedAt    string `json:"createdAt" validate:"datetime=2006-01-02"`
}
