package seeders

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/robertcrysthian/backend-go/constants"
	"github.com/robertcrysthian/backend-go/models"
)

func findStatesFromIBGE () ([]models.IbgeStateEntity, error) {
	resp, err := http.Get(constants.IBGE_URL + "estados")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var States []models.IbgeStateEntity
	err = json.Unmarshal(body, &States)
	if err != nil {
		return nil, err
	}

	return States, nil;
}

func PopulateStates(db *sql.DB) (error) {
	var States, err = findStatesFromIBGE()
	if err != nil { return err } 

	var populateStatesQuery strings.Builder; populateStatesQuery.WriteString(`INSERT INTO states (id, name, state_abbreviation) VALUES`)
	var lastLoopValue = 0
	const parametersLength = 3
	var args = make([]any, 0, len(States)*parametersLength)
	for i :=range States {
		fmt.Fprintf(&populateStatesQuery, `($%d, $%d, $%d)`, lastLoopValue+1, lastLoopValue+2, lastLoopValue+3)
		lastLoopValue += parametersLength
		if i != len(States) -1 {
			populateStatesQuery.WriteString(",")
		}

		args = append(args, States[i].ID, States[i].Name, States[i].StateAbbreviation)
	}
	fmt.Fprintf(&populateStatesQuery, " ON CONFLICT (id) DO NOTHING;")

	_, err = db.Exec(populateStatesQuery.String(), args...)
	if err != nil {
		return err
	}
	return nil
}