package handlers

import (
	"database/sql"
	"net/http"
)

type AllocationsHandler struct {
	DB *sql.DB
}

func NewAllocationsHandler(db *sql.DB) *AllocationsHandler {
	return &AllocationsHandler{DB: db}
}

func (AllocationsHandler *AllocationsHandler) CreateAllocation(writer http.ResponseWriter, request *http.Request) {
	
}