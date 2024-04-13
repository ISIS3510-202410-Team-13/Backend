package models

// Event representa un evento en el horario de una persona
type Event struct {
	DayOfWeek string `json:"dayOfWeek"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

// PlannerEvent representa un evento en el horario con información adicional
type PlannerEvent struct {
	DayOfWeek       string   `json:"dayOfWeek"`
	StartTime       string   `json:"startTime"`
	EndTime         string   `json:"endTime"`
	UsersAvailable  []string `json:"usersAvailable"`
	AmountAvailable int      `json:"amountAvailable"`
}

// TimeBlock representa un bloque de tiempo con un inicio y fin en minutos
type TimeBlock struct {
	StartMinute int `json:"startMinute"`
	EndMinute   int `json:"endMinute"`
}

// UserTimeBlock representa un bloque de tiempo con un inicio y fin en minutos y un ID de usuario
type UserTimeBlock struct {
	UserID string `json:"userID"`
	TimeBlock
}

// ScheduleModel representa el horario de una persona
type ScheduleModel struct {
	ID        string      `json:"id"`
	Monday    []TimeBlock `json:"l"`
	Tuesday   []TimeBlock `json:"m"`
	Wednesday []TimeBlock `json:"i"`
	Thursday  []TimeBlock `json:"j"`
	Friday    []TimeBlock `json:"v"`
	Saturday  []TimeBlock `json:"s"`
	Sunday    []TimeBlock `json:"d"`
}
