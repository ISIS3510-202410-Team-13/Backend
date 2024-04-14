package main

import (
	"log"
	"net/http"
	"os"

	"Planner/api/handlers"
	"Planner/api/middlewares"

	"github.com/gorilla/mux"
)

func main() {

	// Crear un enrutador
	router := mux.NewRouter()

	// Configurar los manejadores de las rutas
	router.HandleFunc("/hello", handlers.HelloHandler).Methods("GET")
	router.Handle("/planner", middlewares.ValidatePlannerBodyMiddleware(http.HandlerFunc(handlers.PlannerHandler))).Methods("POST")

	// Leer el puerto del entorno
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	// Iniciar el servidor en el puerto 8080
	log.Printf("Server listening on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
