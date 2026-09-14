package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/robertcrysthian/backend-go/models"
	"github.com/robertcrysthian/backend-go/utils"
)

type StatesHandler struct {
	DB *sql.DB
}

func NewStatesHandler(db *sql.DB) *StatesHandler {
	return &StatesHandler{DB: db}
}

func (StatesHandler *StatesHandler) FindStates(writer http.ResponseWriter, request *http.Request) {
	const query = `SELECT * FROM STATES`

	rows, err := StatesHandler.DB.Query(query)
	if err != nil {
		utils.InternalServerError(writer, "Erro na query de buscar estados " +err.Error())
		return
	}
	var states = make([]models.ListStateDto, 0)

	for rows.Next() {
		var state models.ListStateDto
		err := rows.Scan(&state.ID, &state.Name, &state.StateAbbreviation)

		if err != nil {
			utils.InternalServerError(writer, err.Error())
			return
		}
		states = append(states, state)
	}
	if err := rows.Err(); err != nil {
		utils.InternalServerError(writer, err.Error())
   		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(states)

}