package handlers

import (
	"database/sql"
	"net/http"

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

	const pickupLocationCreationQuery = `INSERT INTO pickup_locations`
}