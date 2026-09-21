package models

type CreateAllocationDto struct {
	RenterId                int    `json:"renterId" validate:"required,number"`
	VehicleId               int    `json:"vehicleId" validate:"required,number"`
	PickupDate              string `json:"pickupDate" validate:"required,datetime=2006-01-02"`
	EstimatedDevolutionDate string `json:"estimatedDevolutionDate" validate:"required,datetime=2006-01-02"`
}

type ListAllocationDto struct {
	ID                      int     `json:"id"`
	PickupDate              string  `json:"pickupDate"`
	DevolutionDate          string  `json:"devolutionDate"`
	EstimatedDevolutionDate string  `json:"estimatedDevolutionDate"`
	VehicleModel            string  `json:"vehicleModel"`
	VehicleBrand            string  `json:"vehicleBrand"`
	AllocationCost          float32 `json:"allocationCost"`
	LateReturnFeeCost       float32 `json:"lateReturnFeeCost"`
}