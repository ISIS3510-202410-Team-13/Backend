package main

import (
	"log"
	"net/http"

	"Planner/api/handlers"
)

func main() {
	// Configurar los manejadores de las rutas
	http.HandleFunc("/hello", handlers.HelloHandler)
	//http.HandleFunc("/planner", handlers.PlannerHandler)

	// Iniciar el servidor en el puerto 8080
	log.Println("Server listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
