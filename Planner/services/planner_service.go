package services

import (
	"Planner/models"
	"fmt"
	"sort"
)

// Función para normalizar el horario de un usuario
func normalizeSchedule(schedule models.ScheduleModel) models.ScheduleModel {
	// Copiar los TimeBlocks de cada día de la semana a un mapa para realizar la fusión
	timeBlocksMap := make(map[string][]models.TimeBlock)
	timeBlocksMap["l"] = schedule.Monday
	timeBlocksMap["m"] = schedule.Tuesday
	timeBlocksMap["i"] = schedule.Wednesday
	timeBlocksMap["j"] = schedule.Thursday
	timeBlocksMap["v"] = schedule.Friday
	timeBlocksMap["s"] = schedule.Saturday
	timeBlocksMap["d"] = schedule.Sunday

	// Fusionar los TimeBlocks de cada día de la semana
	for _, day := range []string{"l", "m", "i", "j", "v", "s", "d"} {
		timeBlocksMap[day] = mergeBlocks(timeBlocksMap[day])
	}

	// Actualizar el horario del usuario con los TimeBlocks fusionados
	normalizedSchedule := models.ScheduleModel{
		ID:        schedule.ID,
		Monday:    timeBlocksMap["l"],
		Tuesday:   timeBlocksMap["m"],
		Wednesday: timeBlocksMap["i"],
		Thursday:  timeBlocksMap["j"],
		Friday:    timeBlocksMap["v"],
		Saturday:  timeBlocksMap["s"],
		Sunday:    timeBlocksMap["d"],
	}

	return normalizedSchedule
}

// Función para fusionar TimeBlocks sin overlapping
func mergeBlocks(timeBlocks []models.TimeBlock) []models.TimeBlock {
	// Caso que no tiene elementos
	if len(timeBlocks) < 1 {
		return timeBlocks
	}

	// Ordenar el array ascendente según el StartTime
	sort.Slice(timeBlocks, func(i, j int) bool {
		return timeBlocks[i].StartMinute < timeBlocks[j].StartMinute
	})

	// Fusionar TimeBlocks sin overlapping
	mergedBlocks := []models.TimeBlock{timeBlocks[0]}
	for _, block := range timeBlocks[1:] {
		lastAdded := mergedBlocks[len(mergedBlocks)-1]
		if block.StartMinute < lastAdded.EndMinute {
			lastAdded.EndMinute = max(lastAdded.EndMinute, block.EndMinute)
		} else {
			mergedBlocks = append(mergedBlocks, block)
		}
	}

	return mergedBlocks
}

// Función para unir los time blocks de todos los usuarios en un único array por día de la semana
func mergeTimeBlocksByDay(schedules []models.ScheduleModel) map[string][]models.UserTimeBlock {
	mergedTimeBlocks := make(map[string][]models.UserTimeBlock)

	for _, schedule := range schedules {

		// Unir los time blocks por día de la semana
		for _, day := range []string{"l", "m", "i", "j", "v", "s", "d"} {
			switch day {
			case "l":
				mergedTimeBlocks["l"] = append(mergedTimeBlocks["l"], convertToUserTimeBlocks(schedule.Monday, schedule.ID)...)
			case "m":
				mergedTimeBlocks["m"] = append(mergedTimeBlocks["m"], convertToUserTimeBlocks(schedule.Tuesday, schedule.ID)...)
			case "i":
				mergedTimeBlocks["i"] = append(mergedTimeBlocks["i"], convertToUserTimeBlocks(schedule.Wednesday, schedule.ID)...)
			case "j":
				mergedTimeBlocks["j"] = append(mergedTimeBlocks["j"], convertToUserTimeBlocks(schedule.Thursday, schedule.ID)...)
			case "v":
				mergedTimeBlocks["v"] = append(mergedTimeBlocks["v"], convertToUserTimeBlocks(schedule.Friday, schedule.ID)...)
			case "s":
				mergedTimeBlocks["s"] = append(mergedTimeBlocks["s"], convertToUserTimeBlocks(schedule.Saturday, schedule.ID)...)
			case "d":
				mergedTimeBlocks["d"] = append(mergedTimeBlocks["d"], convertToUserTimeBlocks(schedule.Sunday, schedule.ID)...)
			}
		}
	}

	return mergedTimeBlocks
}

// Función auxiliar para convertir TimeBlock a UserTimeBlock con el ID de usuario especificado
func convertToUserTimeBlocks(timeBlocks []models.TimeBlock, userID string) []models.UserTimeBlock {
	userTimeBlocks := make([]models.UserTimeBlock, len(timeBlocks))
	for i, block := range timeBlocks {
		userTimeBlocks[i] = models.UserTimeBlock{
			UserID:    userID,
			TimeBlock: block,
		}
	}
	return userTimeBlocks
}

// Struct para representar un punto de tiempo con información de ocupación
type timePoint struct {
	time   int    // Tiempo en minutos
	userID string // ID del usuario
	isBusy bool   // Indica si el usuario está ocupado (true) o libre (false)
}

// Función para encontrar los time slots disponibles por día
func findAvailableTimeSlotsByDay(day string, timeBlocks []models.UserTimeBlock, usersId []string) []models.PlannerEvent {
	// Array para almacenar los puntos de tiempo con información de ocupación
	var timePoints []timePoint

	// Agregar los timePoints correspondientes a los bloques de tiempo de los usuarios
	for _, tb := range timeBlocks {
		timePoints = append(timePoints, timePoint{time: tb.StartMinute, userID: tb.UserID, isBusy: true})
		timePoints = append(timePoints, timePoint{time: tb.EndMinute, userID: tb.UserID, isBusy: false})
	}

	// Ordenar los timePoints en orden ascendente
	sort.Slice(timePoints, func(i, j int) bool {
		return timePoints[i].time < timePoints[j].time
	})

	// Agregar los tiempos de inicio y fin únicos al array de timePoints
	timePoints = append([]timePoint{{time: 0, userID: "--BEGIN--"}}, timePoints...) // Tiempo inicial del día
	timePoints = append(timePoints, timePoint{time: 1440 - 1, userID: "--END--"})   // Tiempo final del día

	// Lista para almacenar los eventos planificados
	var plannerEvents []models.PlannerEvent

	// Lista para rastrear los usuarios disponibles en cada momento
	availableUsers := make([]string, len(usersId))

	// Inicializar los usuarios disponibles en el tiempo inicial
	copy(availableUsers, usersId)
	sort.Strings(availableUsers)

	// Iterar sobre cada timePoint
	for i := 1; i < len(timePoints); i++ {

		// Actualizar los usuarios disponibles en este intervalo de tiempo
		if timePoints[i].isBusy {
			decrementAvailableUsers(&availableUsers, timePoints[i].userID)
		}

		// Calcular la duración entre el tiempo previo y el tiempo actual
		duration := timePoints[i].time - timePoints[i-1].time

		// Si la duración es mayor que cero, significa que hay un intervalo disponible
		if duration > 0 {

			// Crear un evento de PlannerEvent para este intervalo de tiempo
			event := models.PlannerEvent{
				DayOfWeek:       day,
				StartTime:       convertToTimeString(timePoints[i-1].time),
				EndTime:         convertToTimeString(timePoints[i].time),
				UsersAvailable:  make([]string, len(availableUsers)),
				AmountAvailable: len(availableUsers),
			}

			// Copiar los usuarios disponibles al evento
			copy(event.UsersAvailable, availableUsers)

			// Agregar el evento al slice de PlannerEvent
			plannerEvents = append(plannerEvents, event)
		}

		// Actualizar los usuarios disponibles en este intervalo de tiempo
		if !timePoints[i].isBusy {
			incrementAvailableUsers(&availableUsers, timePoints[i].userID)
		}
	}

	return plannerEvents
}

// Función auxiliar para actualizar los usuarios disponibles en cada intervalo de tiempo
func incrementAvailableUsers(users *[]string, userID string) {
	if userID == "--END--" || userID == "--BEGIN--" {
		return
	}
	*users = append(*users, userID)
	sort.Strings(*users)
}

// Función auxiliar para actualizar los usuarios disponibles en cada intervalo de tiempo
func decrementAvailableUsers(users *[]string, userID string) {
	if userID == "--BEGIN--" || userID == "--END--" {
		return
	}
	for i := len(*users) - 1; i >= 0; i-- {
		if (*users)[i] == userID {
			*users = append((*users)[:i], (*users)[i+1:]...)
			break
		}
	}
	sort.Strings(*users)
}

// Función auxiliar para convertir minutos en formato de cadena de tiempo (HH:MM)
func convertToTimeString(minutes int) string {
	hours := minutes / 60
	mins := minutes % 60
	return fmt.Sprintf("%02d:%02d", hours, mins)
}

// GetAvailableTimeSlots simula la obtención de los slots de tiempo disponibles
func GetAvailableTimeSlots(schedules []models.ScheduleModel) []models.PlannerEvent {

	// Normalizar los horarios de todos los usuarios
	for i := range schedules {
		schedules[i] = normalizeSchedule(schedules[i])
	}

	// Obtener los IDs de los usuarios
	usersId := make([]string, len(schedules))
	for i, schedule := range schedules {
		usersId[i] = schedule.ID
	}

	// Unir los time blocks de todos los usuarios por día de la semana
	mergedTimeBlocksByDay := mergeTimeBlocksByDay(schedules)

	// Canal para recibir los resultados de las coroutines
	results := make(chan []models.PlannerEvent, len(mergedTimeBlocksByDay))

	// Lanzar una coroutine para encontrar los time slots disponibles por día
	for day, userTimeBlocks := range mergedTimeBlocksByDay {
		go func(day string, userTimeBlocks []models.UserTimeBlock) {
			availableSlots := findAvailableTimeSlotsByDay(day, userTimeBlocks, usersId)
			results <- availableSlots
		}(day, userTimeBlocks)
	}

	// Recopilar los resultados de todas las coroutines
	availableSlotsByDay := make([]models.PlannerEvent, 0)
	for range mergedTimeBlocksByDay {
		availableSlots := <-results
		availableSlotsByDay = append(availableSlotsByDay, availableSlots...)
	}

	return availableSlotsByDay
}
