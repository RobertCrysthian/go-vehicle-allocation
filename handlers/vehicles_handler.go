package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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
		INSERT INTO vehicles (brand, model, chassi, year, color, doors_amount, seats_amount, has_air_conditioning, pickup_location_id, cost_per_day, late_return_fee)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
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
		vehicle.LateReturnFee,
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
			cost_per_day,
			late_return_fee
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
			&vehicle.LateReturnFee,
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

func (VehiclesHandler *VehiclesHandler) FindVehicleById (writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.BadRequestError(writer, "Identificador do veículo não foi encontrado ou não é um número válido")
		return
	}

	const findVehicleQuery = `
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
			cost_per_day,
			late_return_fee
		FROM vehicles
		WHERE id = $1
	`
	var vehicle models.ListVehiclesDto

	err = VehiclesHandler.DB.QueryRow(findVehicleQuery, id).Scan(
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
			&vehicle.LateReturnFee,
		)
	
	if err != nil {
		if err == sql.ErrNoRows {
			utils.BadRequestError(writer, "Nenhum veículo encontrado com esse id")
			return
		}
		utils.InternalServerError(writer, "Ocorreu um erro ao escanear o veículo: " + err.Error())
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(vehicle)
}

func (VehiclesHandler *VehiclesHandler) UpdateVehicle (writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.BadRequestError(writer, "Identificador do veículo não foi encontrado ou não é um número válido")
		return
	}

	var vehicle models.UpdateVehicleDto
	err = json.NewDecoder(request.Body).Decode(&vehicle)
	if err != nil {
		http.Error(writer, "A resposta deve ser um JSON válido " + err.Error(), http.StatusBadRequest)
		return
	}

	query := `
		UPDATE vehicles set 
			brand = $1, 
			model = $2, 
			chassi = $3,
			year = $4,
			color = $5,
			doors_amount = $6,
			seats_amount = $7,
			has_air_conditioning = $8,
			pickup_location_id = $9,
			cost_per_day = $10,
			late_return_fee = $11
		WHERE id = $12`

	result, err := VehiclesHandler.DB.Exec(query, 
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
		vehicle.LateReturnFee,
		id,
	)

	if err != nil {
		utils.InternalServerError(writer, err.Error())
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		utils.InternalServerError(writer, err.Error())
		return
	}
	if rowsAffected == 0 {
		http.Error(writer, "Nenhum veículo encontrado com esse id", http.StatusNotFound)
		return
	}

	response := map[string]string{"response": "Veículo editado com sucesso!"}
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(response)
}

func (VehiclesHandler *VehiclesHandler) DeleteVehicle (writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.BadRequestError(writer, "Identificador do veículo não foi encontrado ou não é um número válido")
		return
	}

	query := `DELETE FROM vehicles WHERE id = $1`
	result , err := VehiclesHandler.DB.Exec(query, id)

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	
	if rowsAffected == 0 {
		http.Error(writer, "Nenhum veículo encontrado com esse id", http.StatusNotFound)
		return
	}

	response := map[string]string{"response": "Veículo excluido com sucesso!"}
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(response)
}