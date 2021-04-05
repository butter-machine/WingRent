package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"wingrent/api"
	"wingrent/database"
	"wingrent/database/postgres"
)

var BaseURL = "/api/v1"


func route() {
	r := mux.NewRouter()

	r.HandleFunc(BaseURL + "/planes" , api.ListPlanes).Methods("GET")
	r.HandleFunc(BaseURL + "/planes/{id:[0-9]+}", api.RetrievePlane).Methods("GET")
	r.HandleFunc(BaseURL + "/planes", api.CreatePlane).Methods("POST")
	r.HandleFunc(BaseURL + "/planes/{id:[0-9]+}", api.DeletePlane).Methods("DELETE")
	r.HandleFunc(BaseURL + "/planes/{id:[0-9]+}", api.UpdatePlane).Methods("PATCH")

	http.ListenAndServe(":8080", r)
}

func setupPostgresDBFromENV() {
	dbUser, dbPassword, dbName, dbHost :=
		"wingrent",
		"CoolPass123",
		"wingrent_db",
		"database"
		//os.Getenv("POSTGRES_USER"),
		//os.Getenv("POSTGRES_PASSWORD"),
		//os.Getenv("POSTGRES_DB"),
		//os.Getenv("DB_HOST")

	dbPort, err := strconv.ParseInt("5432", 10, 64)
	p, err := postgres.Initialize(dbUser, dbPassword, dbName, dbHost, int(dbPort))
	if err != nil {
		log.Fatalf("Could not set up database: %v", err)
	}
	database.DBSingleton = p
}

func main() {
	log.Println("Starting api server...")

	setupPostgresDBFromENV()
	log.Println("Database connection established")
	defer database.DBSingleton.Close()

	route()
	log.Println("Ready to receive requests")
}
