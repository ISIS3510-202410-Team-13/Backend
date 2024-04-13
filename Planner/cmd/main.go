package main

import (
	"log"
	"net/http"

	"Planner/api/handlers"
	"Planner/api/middlewares"

	"github.com/gorilla/mux"
)

func main() {

	// Crear un enrutador
	router := mux.NewRouter()

	// Configurar los manejadores de las rutas
	router.HandleFunc("/hello", handlers.HelloHandler).Methods("GET")
	router.Handle("/planner", middlewares.ValidatePlannerBodyMiddleware(http.HandlerFunc(handlers.PlannerHandler))).Methods("GET")

	// Iniciar el servidor en el puerto 8080
	log.Println("Server listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
