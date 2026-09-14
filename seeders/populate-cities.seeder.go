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

func findCitiesFromIBGE () ([]models.IbgeCityEntity, error) {
	resp, err := http.Get(constants.IBGE_URL + "municipios")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var Cities []models.IbgeCityEntity
	err = json.Unmarshal(body, &Cities)
	if err != nil {
		return nil, err
	}

	return Cities, nil;
}

func PopulateCities(db *sql.DB) (error) {
	var Cities, err = findCitiesFromIBGE()
	if err != nil { return err } 

	var populateCitiesQuery strings.Builder; populateCitiesQuery.WriteString(`INSERT INTO cities (id, name, state_id) VALUES`)
	var lastLoopValue = 0
	var parametersLength = 3
	var args = make([]any, 0, len(Cities)*parametersLength)

	for i :=range Cities {
		fmt.Fprintf(&populateCitiesQuery, `($%d, $%d, $%d)`, lastLoopValue+1, lastLoopValue+2, lastLoopValue+3)
		lastLoopValue += parametersLength
		if i != len(Cities) -1 {
			populateCitiesQuery.WriteString(",")
		}

		args = append(args, Cities[i].ID, Cities[i].Name, Cities[i].RegiaoImediata.RegiaoIntermediaria.UF.ID)
	}
	fmt.Fprintf(&populateCitiesQuery, " ON CONFLICT (id) DO NOTHING;")

	_, err = db.Exec(populateCitiesQuery.String(), args...)
	if err != nil {
		return err
	}
	return nil
}