package middlewares

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"Planner/api/constants"
	"Planner/models"
)

// ValidateBodyMiddleware es un middleware que valida el cuerpo de la solicitud para cada Event
func ValidatePlannerBodyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verificar si el cuerpo de la solicitud está presente
		if r.Body == nil {
			http.Error(w, "Request body is missing", http.StatusBadRequest)
			return
		}

		// Decodificar el cuerpo de la solicitud en un Event
		var request map[string][]models.Event
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, "Error while decoding JSON request body in middleware: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Almacenar el cuerpo de la solicitud decodificado en el contexto
		ctx := context.WithValue(r.Context(), constants.ContextKey{Key: "Planner"}, request)

		for _, events := range request {
			for _, event := range events {

				// Validar los campos del Event
				if event.DayOfWeek == "" || event.StartTime == "" || event.EndTime == "" {
					http.Error(w, "Missing required fields dayOfWeek, startTime or endTime in request body", http.StatusBadRequest)
					return
				}

				// Validar el día de la semana
				if !isValidDayOfWeek(event.DayOfWeek) {
					http.Error(w, fmt.Sprintf("Invalid day of week '%s'", event.DayOfWeek), http.StatusBadRequest)
					return
				}

				// Validar el formato de tiempo
				if !isValidTimeFormat(event.StartTime, event.EndTime) {
					http.Error(w, fmt.Sprintf("Invalid time format for '%s' or '%s'", event.StartTime, event.EndTime), http.StatusBadRequest)
					return
				}

				// Verificar si la hora de inicio es anterior a la hora de finalización
				if event.StartTime >= event.EndTime {
					http.Error(w, fmt.Sprintf("Start time '%s' must be before end time '%s'", event.StartTime, event.EndTime), http.StatusBadRequest)
					return
				}

				// Verificar el rango de tiempo
				if !isValidTimeRange(event.StartTime, event.EndTime) {
					http.Error(w, fmt.Sprintf("Invalid time range for '%s' or '%s'", event.StartTime, event.EndTime), http.StatusBadRequest)
					return
				}

				// Verificar los minutos
				if !isValidMinutes(event.StartTime, event.EndTime) {
					http.Error(w, fmt.Sprintf("Invalid minutes for '%s' or '%s'", event.StartTime, event.EndTime), http.StatusBadRequest)
					return
				}

			}
		}

		// Si todas las validaciones pasan, continuar con el siguiente middleware o controlador
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Funciones de ayuda para validaciones específicas

func isValidDayOfWeek(dayOfWeek string) bool {
	return regexp.MustCompile(`^[lmijvsd]$`).MatchString(dayOfWeek)
}

func isValidTimeFormat(startTime, endTime string) bool {
	return regexp.MustCompile(`^\d{4}$`).MatchString(startTime) && regexp.MustCompile(`^\d{4}$`).MatchString(endTime)
}

func isValidTimeRange(startTime, endTime string) bool {
	return (startTime >= "0000" && startTime < "2400" && endTime >= "0000" && endTime < "2400")
}

func isValidMinutes(startTime, endTime string) bool {
	return (startTime[2:] < "60" && endTime[2:] < "60")
}
