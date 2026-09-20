package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/robertcrysthian/backend-go/models"
	"github.com/robertcrysthian/backend-go/utils"
)

type AllocationsHandler struct {
	DB *sql.DB
}

func NewAllocationsHandler(db *sql.DB) *AllocationsHandler {
	return &AllocationsHandler{DB: db}
}

func (AllocationsHandler *AllocationsHandler) CreateAllocation(writer http.ResponseWriter, request *http.Request) {
	var allocation models.CreateAllocationDto
	err := utils.ValidateRequest(request, &allocation)
	if err != nil {
		utils.BadRequestError(writer, err.Error())
		return
	}

	const findUserQuery = `
		SELECT EXISTS (
			SELECT 1 FROM users WHERE id = $1
		)
	`
	var isUserValid bool
	err = AllocationsHandler.DB.QueryRow(findUserQuery, allocation.RenterId).Scan(&isUserValid)
	if err != nil {
		utils.InternalServerError(writer, "Ocorreu um erro ao validar usuário " + err.Error())
		return
	}

	if !isUserValid {
		utils.NotFoundError(writer, "Usuário não encontrado")
		return
	}


	const findVehicleQuery = `
		SELECT EXISTS (
			SELECT 1 FROM vehicles WHERE id = $1
		)
	`
	var isValidVehicle bool
	err = AllocationsHandler.DB.QueryRow(findVehicleQuery, allocation.VehicleId).Scan(&isValidVehicle)
	if err != nil {
		utils.InternalServerError(writer, "Ocorreu um erro ao validar o veículo " + err.Error())
		return
	}

	if !isValidVehicle {
		utils.NotFoundError(writer, "Veículo não encontrado")
		return
	}

	isVehicleAvailable, err := verifyIfVehicleIsAvailable(*AllocationsHandler.DB, allocation.PickUpDate, allocation.EstimatedDevolutionDate)
	if err != nil {
		utils.InternalServerError(writer, "Ocorreu um erro ao validar a data de disponibilidade do veículo " + err.Error())
	}

	if !isVehicleAvailable {
		utils.BadRequestError(writer, "O veículo selecionado não se encontra disponível nessa data")
	}

	const creationQuery = `
		INSERT INTO allocations (renter_id, vehicle_id, pick_up_date, estimated_devolution_date)
		VALUES ($1, $2, $3, $4)
	`

	_, err = AllocationsHandler.DB.Exec(
		creationQuery, 
		allocation.RenterId,
		 allocation.VehicleId, 
		 allocation.PickUpDate, 
		 allocation.EstimatedDevolutionDate,
	)
	if err != nil {
		utils.InternalServerError(writer, "Erro na query de criar alocação "+err.Error())
		return
	}

	response := map[string]string{"response": "Alocação criada com sucesso!"}
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(response)

}

func verifyIfVehicleIsAvailable(db sql.DB, initialDate string, finalDate string) (bool, error) {

	const query = `
		SELECT EXISTS(
			SELECT 1 FROM vehicles v
			LEFT JOIN allocations a ON a.vehicle_id = v.id 
			WHERE a.vehicle_id = 2 AND a.pick_up_date < $2 AND COALESCE(a.devolution_date, a.estimated_devolution_date) > $1
		)
	`
	var isVehicleAvailable bool
	err := db.QueryRow(query, initialDate, finalDate).Scan(&isVehicleAvailable)
	if err != nil {
		return false, err
	}
	return !isVehicleAvailable, nil

}