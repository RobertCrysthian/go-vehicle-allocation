package models

type CreatePickupLocationDto struct {
	CityId         int    `json:"cityId" validate:"required,number"`
	StateId        int    `json:"stateId" validate:"required,number"`
	District       string `json:"district" validate:"required"`
	Street         string `json:"street" validate:"required"`
	BuildingNumber string `json:"buildingNumber" validate:"required"`
}
