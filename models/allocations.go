package models

type CreateAllocationDto struct {
	RenterId                int    `json:"renterId" validate:"required,number"`
	VehicleId               int    `json:"vehicleId" validate:"required,number"`
	PickUpDate              string `json:"pickUpDate" validate:"required,datetime=2006-01-02"`
	EstimatedDevolutionDate string `json:"estimatedDevolutionDate" validate:"required,datetime=2006-01-02"`
}