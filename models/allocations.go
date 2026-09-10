package models

type CreateAllocationDto struct {
	RenterId                int    `json:"renterId" validate:"required,number"`
	VehicleId               int    `json:"vehicleId" validate:"required,number"`
	PickUpDate              string `json:"pickUpDate" validate:"required,datetime"`
	EstimatedDevolutionDate string `json:"estimatedDevolutionDate" validate:"required,datetime"`
}