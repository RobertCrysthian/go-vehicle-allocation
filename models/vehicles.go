package models

type CreateVehicleDto struct {
	Brand              string `json:"brand" validate:"required"`
	Model              string `json:"model" validate:"required"`
	Chassi             string `json:"chassi" validate:"required"`
	Year               int    `json:"year" validate:"required,number"`
	Color              string `json:"color" validate:"required"`
	DoorsAmount        int    `json:"doorsAmount" validate:"required,number"`
	SeatsAmount        int    `json:"seatsAmount" validate:"required,number"`
	HasAirConditioning bool   `json:"hasAirConditioning"`
	PickupLocationId   int    `json:"pickupLocationId" validate:"required,number"`
	CostPerDay 				 string `json:"costPerDay" validate:"required,number"`
}

type ListVehiclesDto struct {
	ID                 int    `json:"id"`
	Brand              string `json:"brand"`
	Model              string `json:"model"`
	Chassi             string `json:"chassi"`
	Year               int    `json:"year"`
	Color              string `json:"color"`
	DoorsAmount        int    `json:"doorsAmount"`
	SeatsAmount        int    `json:"seatsAmount"`
	HasAirConditioning bool   `json:"hasAirConditioning"`
	PickupLocationId   int    `json:"pickupLocationId"`
	CostPerDay         string `json:"costPerDay"`
}
