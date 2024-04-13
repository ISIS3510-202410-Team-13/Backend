package handlers

import (
	"encoding/json"
	"net/http"
)

// Message struct para manejar el mensaje de respuesta
type HelloMessage struct {
	Message string `json:"message"`
}

// HelloHandler maneja las solicitudes a la ruta /hello
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	// Definir el mensaje de respuesta
	message := HelloMessage{
		Message: "¡Hola desde el servidor Planner!",
	}

	// Establecer el tipo de contenido de la respuesta como JSON
	w.Header().Set("Content-Type", "application/json")

	// Codificar el mensaje en JSON y escribirlo en el cuerpo de la respuesta
	json.NewEncoder(w).Encode(message)
}
