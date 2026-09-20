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

type PickupLocationHandler struct {
	DB *sql.DB
}

func NewPickupLocationHandler(db *sql.DB) *PickupLocationHandler {
	return &PickupLocationHandler{DB: db}
}

func (PickupLocationHandler *PickupLocationHandler) CreatePickupLocation(writer http.ResponseWriter, request *http.Request) {
	var pickupLocation models.CreatePickupLocationDto
	err := utils.ValidateRequest(request, &pickupLocation)
	if err != nil {
		utils.BadRequestError(writer, err.Error())
		return
	}

	var isValidCity bool
	const findCityByIdQuery = "SELECT EXISTS (SELECT 1 FROM cities WHERE id = $1)"
	row := PickupLocationHandler.DB.QueryRow(findCityByIdQuery, strconv.Itoa(pickupLocation.CityId))
	err = row.Scan(&isValidCity)

	if err != nil {
		utils.InternalServerError(writer, "Ocorreu um erro na query de validar cidades " + err.Error())
		return
	}

	if !isValidCity {
		utils.BadRequestError(writer, "Nenhuma cidade encontrada com esse id")
		return
	}


	const pickupLocationCreationQuery = "INSERT INTO pickup_locations (city_id, district, street, building_number) VALUES ($1, $2, $3, $4)"
	_, err = PickupLocationHandler.DB.Exec(
		pickupLocationCreationQuery, 
		pickupLocation.CityId,
		pickupLocation.District,
		pickupLocation.Street,
		pickupLocation.BuildingNumber,
	)
	if err != nil {
		utils.InternalServerError(writer, "Erro na query de criar local de coleta " +err.Error())
		return
	}

	response := map[string]string{"response": "Ponto de coleta criado com sucesso!"}
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(response)
}

func (PickupLocationHandler *PickupLocationHandler) FindAllPickupLocations (writer http.ResponseWriter, request *http.Request) {
	const findAllPickupLocationsQuery = `
		SELECT 
			pl.id,
			pl.district,
			pl.street,
			pl.building_number,
			c.name as city_name,
			s.name as state_name
		FROM pickup_locations pl
		JOIN cities c ON c.id = pl.city_id 
		JOIN states s ON s.id = c.state_id
	`
	rows, err := PickupLocationHandler.DB.Query(findAllPickupLocationsQuery)
	if err != nil {
		utils.InternalServerError(writer, "Erro na query de buscar locais de coleta " +err.Error())
		return
	}
	
	var pickupLocations = make([]models.ListPickupLocationDto, 0)
	for rows.Next() {
		var pickupLocation models.ListPickupLocationDto
		err := rows.Scan(
			&pickupLocation.ID, 
			&pickupLocation.District, 
			&pickupLocation.Street, 
			&pickupLocation.BuildingNumber, 
			&pickupLocation.CityName,
			&pickupLocation.StateName,
		)

		if err != nil {
			utils.InternalServerError(writer, "Ocorreu um erro ao escanear local " + err.Error())
			return
		}
		pickupLocations = append(pickupLocations, pickupLocation)
	}
	if err := rows.Err(); err != nil {
		utils.InternalServerError(writer, err.Error())
   		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(pickupLocations)

}

func (PickupLocationHandler *PickupLocationHandler) FindPickupLocationById (writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.BadRequestError(writer, "Identificador do local de coleta não foi encontrado ou não é um número válido")
		return
	}

	var pickupLocation models.ListPickupLocationDto
	const findPickupLocationQuery = `
		SELECT
			pl.id,
			pl.district,
			pl.street,
			pl.building_number,
			c.name as city_name,
			s.name as state_name
		FROM pickup_locations pl
		JOIN cities c ON c.id = pl.city_id 
		JOIN states s ON s.id = c.state_id
		WHERE pl.id = $1
	`

	err = PickupLocationHandler.DB.QueryRow(findPickupLocationQuery, id).Scan(
		&pickupLocation.ID, 
		&pickupLocation.District,
		&pickupLocation.Street,
		&pickupLocation.BuildingNumber,
		&pickupLocation.CityName,
		&pickupLocation.StateName,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			utils.BadRequestError(writer, "Nenhum ponto de coleta encontrado com esse id")
			return
		}
		utils.InternalServerError(writer, "Ocorreu um erro ao escanear o ponto de coleta: " + err.Error())
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(pickupLocation)

}

func (PickupLocationHandler *PickupLocationHandler) UpdatePickupLocation (writer http.ResponseWriter, request *http.Request) {
	var pickupLocation models.UpdatePickupLocationDto
	err := utils.ValidateRequest(request, &pickupLocation)
	if err != nil {
		utils.BadRequestError(writer, err.Error())
		return
	}

	const query = `
		UPDATE pickup_locations SET
			city_id = $1,
			district = $2,
			street = $3, 
			building_number = $4
		WHERE id = $5
	`

	result, err := PickupLocationHandler.DB.Exec(
		query, 
		pickupLocation.CityId, 
		pickupLocation.District, 
		pickupLocation.Street, 
		pickupLocation.BuildingNumber, 
		pickupLocation.ID,
	)

	if err != nil {
		utils.InternalServerError(writer, "Ocorreu um erro ao realizar a query " + err.Error())
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		utils.InternalServerError(writer, err.Error())
		return
	}
	if rowsAffected == 0 {
		http.Error(writer, "Nenhum ponto de coleta encontrado com esse id", http.StatusNotFound)
		return
	}

	response := map[string]string{"response": "Ponto de coleta editado com sucesso!"}
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(response)
}

func (PickupLocationHandler *PickupLocationHandler) DeletePickupLocation (writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.BadRequestError(writer, "Identificador do local de coleta não foi encontrado ou não é um número válido")
		return
	}

	query := `DELETE FROM pickup_locations WHERE id = $1`
	result , err := PickupLocationHandler.DB.Exec(query, id)

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	
	if rowsAffected == 0 {
		http.Error(writer, "Nenhum local de coleta encontrado com esse id", http.StatusNotFound)
		return
	}

	response := map[string]string{"response": "Local de coleta excluido com sucesso!"}
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(response)
}
