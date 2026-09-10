package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/robertcrysthian/backend-go/config"
	"github.com/robertcrysthian/backend-go/handlers"
	"github.com/robertcrysthian/backend-go/models"
)

func main() {
	dbConnection := config.SetupDB()

	_, err := dbConnection.Exec(models.CreateTablesSQL);

	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()

	// taskHandler := handlers.NewTaskHandler(dbConnection)
	userHandler := handlers.NewUsersHandler(dbConnection)
	router.HandleFunc("/users/create", userHandler.CreateUser).Methods("POST")

	vehiclesHandler := handlers.NewVehiclesHandler(dbConnection)
	router.HandleFunc("/vehicles/create", vehiclesHandler.CreateVehicle).Methods("POST")


	defer dbConnection.Close()

	log.Fatal(http.ListenAndServe(":8080", router))
}