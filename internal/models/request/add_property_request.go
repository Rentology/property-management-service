package request

import "property-managment-service/internal/models"

type AddPropertyRequest struct {
	Property        *models.Property        `json:"property"`
	PropertyDetails *models.PropertyDetails `json:"propertyDetails"`
	Images          []string                `json:"images"`
}
