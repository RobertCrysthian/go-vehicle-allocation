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

	isVehicleAvailable, err := verifyIfVehicleIsAvailable(*AllocationsHandler.DB, allocation.PickupDate, allocation.EstimatedDevolutionDate)
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
		 allocation.PickupDate, 
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

func (AllocationsHandler *AllocationsHandler) FindAllAllocations(writer http.ResponseWriter, request *http.Request) {
	const query = `
		SELECT
			a.id,
			a.pick_up_date,
			a.devolution_date,
			a.estimated_devolution_date,
			v.model,
			v.brand,
			v.cost_per_day * (COALESCE(a.devolution_date::date, a.estimated_devolution_date::date) - a.pick_up_date::date) as "allocation_cost",
			CASE 
				WHEN a.devolution_date IS NULL THEN 0
				ELSE GREATEST(a.devolution_date::date - a.estimated_devolution_date::date, 0) * v.late_return_fee
			END as late_return_fee_cost
		FROM allocations a
		JOIN users u ON u.id = a.renter_id 
		JOIN vehicles v ON v.id = a.vehicle_id
	`
	rows, err := AllocationsHandler.DB.Query(query)
	if err != nil {
		utils.InternalServerError(writer, "Erro na query de buscar alocações " +err.Error())
		return
	}
	var allocations = make([]models.ListAllocationDto, 0)

	for rows.Next() {
		var allocation models.ListAllocationDto
		err := rows.Scan(
			&allocation.ID, 
			&allocation.PickupDate, 
			&allocation.DevolutionDate, 
			&allocation.EstimatedDevolutionDate, 
			&allocation.VehicleModel,
			&allocation.VehicleBrand,
			&allocation.AllocationCost,
			&allocation.LateReturnFeeCost,
		)

		if err != nil {
			utils.InternalServerError(writer, "Ocorreu um erro ao escanear a alocação " + err.Error())
			return
		}
		allocations = append(allocations, allocation)
	}

		if err := rows.Err(); err != nil {
		utils.InternalServerError(writer, err.Error())
   		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(allocations)
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