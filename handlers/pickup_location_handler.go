package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

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