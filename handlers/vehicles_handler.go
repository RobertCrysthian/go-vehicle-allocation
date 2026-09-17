package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/robertcrysthian/backend-go/models"
	"github.com/robertcrysthian/backend-go/utils"
)

type VehiclesHandler struct {
	DB *sql.DB
}

func NewVehiclesHandler(db *sql.DB) *VehiclesHandler {
	return &VehiclesHandler{DB: db}
}

func (VehiclesHandler *VehiclesHandler) CreateVehicle(writer http.ResponseWriter, request *http.Request) {
	var vehicle models.CreateVehicleDto
	err := utils.ValidateRequest(request, &vehicle)
	if err != nil {
		utils.BadRequestError(writer, err.Error())
		return
	}

	const checkChassiAvailabilityQuery = `
		SELECT EXISTS (
			SELECT 1 FROM vehicles
			WHERE chassi = $1
		)
	`
	var existingVehicleChassi bool

	if err := VehiclesHandler.DB.QueryRow(checkChassiAvailabilityQuery, vehicle.Chassi).Scan(&existingVehicleChassi); err != nil {
		utils.InternalServerError(writer, "Erro na query que valida chassi "+err.Error())
		return
	}

	if existingVehicleChassi {
		utils.ConflitctError(writer, "Esse veículo já está cadastrado no sistema!")
		return
	}

	const createVehicleQuery = `
		INSERT INTO vehicles (brand, model, chassi, year, color, doors_amount, seats_amount, has_air_conditioning, pickup_location_id, cost_per_day)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err = VehiclesHandler.DB.Exec(
		createVehicleQuery,
		vehicle.Brand,
		vehicle.Model,
		vehicle.Chassi,
		vehicle.Year,
		vehicle.Color,
		vehicle.DoorsAmount,
		vehicle.SeatsAmount,
		vehicle.HasAirConditioning,
		vehicle.PickupLocationId,
		vehicle.CostPerDay,
	)
	if err != nil {
		utils.InternalServerError(writer, "Erro na query de inserir veículo "+err.Error())
		return
	}

	response := map[string]string{"response": "Veículo cadastrado com sucesso!"}
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(response)
}

func (VehiclesHandler *VehiclesHandler) FindAllVehicles (writer http.ResponseWriter, request *http.Request) {
	const findAllVehiclesQuery = `
		SELECT 
			id,
			brand,
			model,
			chassi,
			year,
			color,
			doors_amount,
			seats_amount,
			has_air_conditioning,
			cost_per_day
		FROM vehicles
	`
	rows, err := VehiclesHandler.DB.Query(findAllVehiclesQuery)
	if err != nil {
		utils.InternalServerError(writer, "Erro na query de buscar veículos " +err.Error())
		return
	}
	
	var vehicles = make([]models.ListVehiclesDto, 0)
	for rows.Next() {
		var vehicle models.ListVehiclesDto
		err := rows.Scan(
			&vehicle.ID, 
			&vehicle.Brand, 
			&vehicle.Model, 
			&vehicle.Chassi, 
			&vehicle.Year,
			&vehicle.Color,
			&vehicle.DoorsAmount,
			&vehicle.SeatsAmount,
			&vehicle.HasAirConditioning,
			&vehicle.CostPerDay,
		)

		if err != nil {
			utils.InternalServerError(writer, "Ocorreu um erro ao escanear o veículo " + err.Error())
			return
		}
		vehicles = append(vehicles, vehicle)
	}
	if err := rows.Err(); err != nil {
		utils.InternalServerError(writer, err.Error())
   		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(vehicles)

}