package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/robertcrysthian/backend-go/config"
	"github.com/robertcrysthian/backend-go/handlers"
	"github.com/robertcrysthian/backend-go/models"
	"github.com/robertcrysthian/backend-go/seeders"
)

func main() {
	dbConnection := config.SetupDB()

	_, err := dbConnection.Exec(models.CreateTablesSQL);

	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()

	err = seeders.PopulateStates(dbConnection)
	if err != nil {
		log.Fatal("Ocorreu um erro no seeder de estados " + err.Error())
	}
	err = seeders.PopulateCities(dbConnection)
		if err != nil {
		log.Fatal("Ocorreu um erro no seeder de cidades " + err.Error())
	}

	statesHandler := handlers.NewStatesHandler(dbConnection)
	userHandler := handlers.NewUsersHandler(dbConnection)
	router.HandleFunc("/users/create", userHandler.CreateUser).Methods("POST")

	vehiclesHandler := handlers.NewVehiclesHandler(dbConnection)
	router.HandleFunc("/vehicles/create", vehiclesHandler.CreateVehicle).Methods("POST")
	router.HandleFunc("/states", statesHandler.FindStates).Methods("GET")

	defer dbConnection.Close()

	log.Fatal(http.ListenAndServe(":8080", router))
}