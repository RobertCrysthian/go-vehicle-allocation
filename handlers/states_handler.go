package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/robertcrysthian/backend-go/models"
	"github.com/robertcrysthian/backend-go/utils"
)

type StatesHandler struct {
	DB *sql.DB
}

func NewStatesHandler(db *sql.DB) *StatesHandler {
	return &StatesHandler{DB: db}
}

func (StatesHandler *StatesHandler) PopulateStates(writer http.ResponseWriter, request *http.Request) {
	resp, err := http.Get("https://servicodados.ibge.gov.br/api/v1/localidades/estados")
	if err != nil {
		utils.InternalServerError(writer, "Erro ao buscar os estados na api do ibge")
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.InternalServerError(writer, "Erro ao ler a resposta da api do ibge")
		return
	}

	var States []models.IbgeStateEntity
	err = json.Unmarshal(body, &States)
	if err != nil {
		utils.InternalServerError(writer, "Erro ao converter os estados para struct")
		return
	}

	var populateStatesQuery strings.Builder; populateStatesQuery.WriteString(`INSERT INTO states (name, state_abbreviation, ibge_id) VALUES`)
	var lastLoopValue = 0
	var args = make([]any, 0, len(States)*3)
	for i :=0; i < len(States); i++ {
		fmt.Fprintf(&populateStatesQuery, `($%d, $%d, $%d)`, lastLoopValue+1, lastLoopValue+2, lastLoopValue+3)
		lastLoopValue += 3 
		if i != len(States) -1 {
			populateStatesQuery.WriteString(",")
		}

		args = append(args, States[i].Name, States[i].StateAbbreviation, strconv.Itoa(States[i].ID))
	}

	_, err = StatesHandler.DB.Exec(populateStatesQuery.String(), args...)
	if err != nil {
		utils.InternalServerError(writer, "Erro na query de inserir estados " +err.Error())
		return
	}

	response := map[string]string{"response": "Estados populados com sucesso!"}
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(response)
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
		err := rows.Scan(&state.ID, &state.Name, &state.StateAbbreviation, &state.IbgeId)

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