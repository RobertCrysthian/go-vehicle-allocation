package models

type CreatePickupLocationDto struct {
	CityId         int    `json:"cityId" validate:"required,number"`
	District       string `json:"district" validate:"required"`
	Street         string `json:"street" validate:"required"`
	BuildingNumber int    `json:"buildingNumber" validate:"required,number"`
}

type UpdatePickupLocationDto struct {
	ID int `json:"id" validate:"required,number"`
	CreatePickupLocationDto
}

type ListPickupLocationDto struct {
	ID             int    `json:"id"`
	District       string `json:"district"`
	Street         string `json:"street"`
	BuildingNumber int    `json:"buildingNumber"`
	CityName       string `json:"cityName"`
	StateName      string `json:"stateName"`
}
