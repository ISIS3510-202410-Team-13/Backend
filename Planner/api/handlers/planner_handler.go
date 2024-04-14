package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"Planner/models"
	"Planner/services"
)

// PlannerRequest es la estructura que representa la solicitud para la ruta /planner
type PlannerRequest map[string][]models.Event

// PlannerResponse es la estructura que representa la respuesta para la ruta /planner
type PlannerResponse []models.PlannerEvent

// PlannerHandler maneja las solicitudes a la ruta /planner
func PlannerHandler(w http.ResponseWriter, r *http.Request) {
	// Decodificar el cuerpo de la solicitud en un PlannerRequest
	var request PlannerRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Error while decoding request:"+err.Error(), http.StatusBadRequest)
		return
	}

	// Convertir PlannerRequest a un mapa con el ID del usuario y valor ScheduleModel
	var userSchedules []models.ScheduleModel
	for userID, events := range request {
		// Crear ScheduleModel para el usuario actual
		userSchedule := models.ScheduleModel{
			ID: userID,
		}

		// Convertir los eventos a TimeBlock y agregarlos al ScheduleModel
		for _, event := range events {
			// Convertir la hora de inicio y fin a minutos (asumiendo que están en formato "HHMM")
			startHour := event.StartTime[:2]
			startMinute := event.StartTime[2:]
			endHour := event.EndTime[:2]
			endMinute := event.EndTime[2:]

			startHourInt, err := strconv.Atoi(startHour)
			if err != nil {
				http.Error(w, "Error al convertir la hora de inicio a entero", http.StatusBadRequest)
			}
			startMinuteInt, err := strconv.Atoi(startMinute)
			if err != nil {
				http.Error(w, "Error al convertir el minuto de inicio a entero", http.StatusBadRequest)
			}
			endHourInt, err := strconv.Atoi(endHour)
			if err != nil {
				http.Error(w, "Error al convertir la hora de fin a entero", http.StatusBadRequest)
			}
			endMinuteInt, err := strconv.Atoi(endMinute)
			if err != nil {
				http.Error(w, "Error al convertir el minuto de fin a entero", http.StatusBadRequest)
			}

			// Calcular los minutos totales desde la medianoche para la hora de inicio y fin
			startTimeMinutes := (60 * startHourInt) + startMinuteInt
			endTimeMinutes := (60 * endHourInt) + endMinuteInt

			// Crear un TimeBlock con la hora de inicio y fin convertidas a minutos
			timeBlock := models.TimeBlock{
				StartMinute: startTimeMinutes,
				EndMinute:   endTimeMinutes,
			}

			// Determinar el día de la semana y agregar el TimeBlock al ScheduleModel correspondiente
			switch event.DayOfWeek {
			case "l":
				userSchedule.Monday = append(userSchedule.Monday, timeBlock)
			case "m":
				userSchedule.Tuesday = append(userSchedule.Tuesday, timeBlock)
			case "i":
				userSchedule.Wednesday = append(userSchedule.Wednesday, timeBlock)
			case "j":
				userSchedule.Thursday = append(userSchedule.Thursday, timeBlock)
			case "v":
				userSchedule.Friday = append(userSchedule.Friday, timeBlock)
			case "s":
				userSchedule.Saturday = append(userSchedule.Saturday, timeBlock)
			case "d":
				userSchedule.Sunday = append(userSchedule.Sunday, timeBlock)
			default:
				http.Error(w, "Día de la semana no válido", http.StatusBadRequest)
			}
		}
		// Agregar el ScheduleModel al array
		userSchedules = append(userSchedules, userSchedule)
	}

	// Llamar al servicio GetAvailableTimeSlots
	availableSlots := services.GetAvailableTimeSlots(userSchedules)

	// Establecer el tipo de contenido de la respuesta como JSON
	w.Header().Set("Content-Type", "application/json")

	// Codificar la respuesta en JSON y escribirla en el cuerpo de la respuesta
	json.NewEncoder(w).Encode(availableSlots)
}
