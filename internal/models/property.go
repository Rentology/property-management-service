package models

type Property struct {
	ID           int64  `json:"id"`
	OwnerId      int64  `json:"ownerId" db:"owner_id"`
	Title        string `json:"title"`
	Location     string `json:"location"`
	Price        int    `json:"price"`
	PropertyType string `json:"propertyType" db:"property_type" validate:"oneof=house apartment"`
	RentalType   string `json:"rentalType" db:"rental_type" validate:"oneof=shortTerm longTerm"`
	MaxGuests    int    `json:"maxGuests" db:"max_guests"`
	CreatedAt    string `json:"createdAt" db:"created_at" validate:"datetime=2006-01-02"`
}
