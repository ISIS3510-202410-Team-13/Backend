package middlewares

import (
	"encoding/json"
	"net/http"
	"regexp"

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
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		for _, events := range request {
			for _, event := range events {

				// Validar los campos del Event
				if event.DayOfWeek == "" || event.StartTime == "" || event.EndTime == "" {
					http.Error(w, "Missing required fields in request body", http.StatusBadRequest)
					return
				}

				// Validar el día de la semana
				if !isValidDayOfWeek(event.DayOfWeek) {
					http.Error(w, "Invalid day of week", http.StatusBadRequest)
					return
				}

				// Validar el formato de tiempo
				if !isValidTimeFormat(event.StartTime, event.EndTime) {
					http.Error(w, "Invalid time format", http.StatusBadRequest)
					return
				}

				// Verificar si la hora de inicio es anterior a la hora de finalización
				if event.StartTime >= event.EndTime {
					http.Error(w, "Start time must be before end time", http.StatusBadRequest)
					return
				}

				// Verificar el rango de tiempo
				if !isValidTimeRange(event.StartTime, event.EndTime) {
					http.Error(w, "Invalid time range", http.StatusBadRequest)
					return
				}

				// Verificar los minutos
				if !isValidMinutes(event.StartTime, event.EndTime) {
					http.Error(w, "Invalid minutes", http.StatusBadRequest)
					return
				}

			}
		}

		// Si todas las validaciones pasan, continuar con el siguiente middleware o controlador
		next.ServeHTTP(w, r)
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
