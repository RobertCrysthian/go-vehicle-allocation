package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/robertcrysthian/backend-go/models"
	"github.com/robertcrysthian/backend-go/utils"
)

type CitiesHandler struct {
	DB *sql.DB
}

func NewCitiesHandler(db *sql.DB) *CitiesHandler {
	return &CitiesHandler{DB: db}
}

func  (CitiesHandler *CitiesHandler) FindCities (writer http.ResponseWriter, request *http.Request) {
	page, err := strconv.Atoi(request.URL.Query().Get("page"))
	if err != nil {
		utils.BadRequestError(writer, "É necessário informar a página como um parâmetro válido")
		return 
	}

	results, err := strconv.Atoi(request.URL.Query().Get("results"))
	if err != nil {
		utils.BadRequestError(writer, "É necessário informar os resultados como um parâmetro válido")
		return 
	}

	var total int 
	totalQuery := "SELECT COUNT(id) FROM cities"
	err = CitiesHandler.DB.QueryRow(totalQuery).Scan(&total)
	if err != nil {
		utils.InternalServerError(writer, "Erro ao buscar total de cidades " + err.Error())
	}

	findCitiesQuery := `
	SELECT 
		c.id,
		c.name,
		s.state_abbreviation
	FROM cities c
	JOIN states s ON s.id = c.state_id 
	LIMIT $1
	OFFSET $2
	`
	offset := (page - 1) * results
	rows, err := CitiesHandler.DB.Query(findCitiesQuery, results, offset)
	if err != nil {
		utils.InternalServerError(writer, "Erro na query de buscar cidades " +err.Error())
		return
	}

	var cities = make([]models.ListCityDto, 0)

	for rows.Next() {
		var city models.ListCityDto
		err := rows.Scan(&city.ID, &city.Name, &city.StateAbbreviation)

		if err != nil {
			utils.InternalServerError(writer, err.Error())
			return
		}
		cities = append(cities, city)
	}

	if err := rows.Err(); err != nil {
	utils.InternalServerError(writer, err.Error())
	return
	}

	var response models.ListCityWithPaginationDto

	response.Cities = cities
	response.Total = total
	response.Page = page
	response.Results = results

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(response)
}