package services

import (
	"reflect"
	"testing"

	"Planner/models"
)

func TestFindAvailableTimeSlotsByDay(t *testing.T) {
	// Caso de prueba: día sin ningún bloque de tiempo
	day := "l"
	userTimeBlocks := []models.UserTimeBlock{}
	expectedSlots := []models.PlannerEvent{
		{
			DayOfWeek:      "l",
			StartTime:      "00:00",
			EndTime:        "23:59",
			Attendees:      []string{},
			UsersAvailable: 0,
			Duration:       1440,
		},
	}
	assertTimeSlots(t, day, []string{}, userTimeBlocks, expectedSlots)

	// Caso de prueba: día con un solo bloque de tiempo que abarca todo el día
	day = "m"
	userTimeBlocks = []models.UserTimeBlock{
		{UserID: "1", TimeBlock: models.TimeBlock{StartMinute: 0, EndMinute: 1440 - 1}},
	}
	expectedSlots = []models.PlannerEvent{
		{
			DayOfWeek:      "m",
			StartTime:      "00:00",
			EndTime:        "23:59",
			Attendees:      []string{},
			UsersAvailable: 0,
			Duration:       1440,
		},
	}
	assertTimeSlots(t, day, []string{"1"}, userTimeBlocks, expectedSlots)

	// Caso de prueba: usuario sin bloques de tiempo que abarcan todo el día
	day = "m"
	userTimeBlocks = []models.UserTimeBlock{}
	expectedSlots = []models.PlannerEvent{
		{
			DayOfWeek:      "m",
			StartTime:      "00:00",
			EndTime:        "23:59",
			Attendees:      []string{"1"},
			UsersAvailable: 1,
			Duration:       1440,
		},
	}
	assertTimeSlots(t, day, []string{"1"}, userTimeBlocks, expectedSlots)

	// Caso de prueba: día con múltiples bloques de tiempo que no se solapan
	day = "i"
	userTimeBlocks = []models.UserTimeBlock{
		{UserID: "1", TimeBlock: models.TimeBlock{StartMinute: 0, EndMinute: 360 - 1}},
		{UserID: "2", TimeBlock: models.TimeBlock{StartMinute: 360, EndMinute: 720 - 1}},
		{UserID: "3", TimeBlock: models.TimeBlock{StartMinute: 720, EndMinute: 1080 - 1}},
		{UserID: "4", TimeBlock: models.TimeBlock{StartMinute: 1080, EndMinute: 1440 - 1}},
	}
	expectedSlots = []models.PlannerEvent{
		{
			DayOfWeek:      "i",
			StartTime:      "00:00",
			EndTime:        "05:59",
			Attendees:      []string{"2", "3", "4"},
			UsersAvailable: 3,
			Duration:       360,
		},
		{
			DayOfWeek:      "i",
			StartTime:      "06:00",
			EndTime:        "11:59",
			Attendees:      []string{"1", "3", "4"},
			UsersAvailable: 3,
			Duration:       360,
		},
		{
			DayOfWeek:      "i",
			StartTime:      "12:00",
			EndTime:        "17:59",
			Attendees:      []string{"1", "2", "4"},
			UsersAvailable: 3,
			Duration:       360,
		},
		{
			DayOfWeek:      "i",
			StartTime:      "18:00",
			EndTime:        "23:59",
			Attendees:      []string{"1", "2", "3"},
			UsersAvailable: 3,
			Duration:       360,
		},
	}
	assertTimeSlots(t, day, []string{"1", "2", "3", "4"}, userTimeBlocks, expectedSlots)

	// Caso de prueba: día con un solo bloque de tiempo que no cubre las 24 horas
	day = "v"
	userTimeBlocks = []models.UserTimeBlock{
		{UserID: "1", TimeBlock: models.TimeBlock{StartMinute: 60, EndMinute: 600 - 1}},
	}
	expectedSlots = []models.PlannerEvent{
		{
			DayOfWeek:      "v",
			StartTime:      "00:00",
			EndTime:        "00:59",
			Attendees:      []string{"1"},
			UsersAvailable: 1,
			Duration:       60,
		},
		{
			DayOfWeek:      "v",
			StartTime:      "01:00",
			EndTime:        "09:59",
			Attendees:      []string{},
			UsersAvailable: 0,
			Duration:       540,
		},
		{
			DayOfWeek:      "v",
			StartTime:      "10:00",
			EndTime:        "23:59",
			Attendees:      []string{"1"},
			UsersAvailable: 1,
			Duration:       840,
		},
	}
	assertTimeSlots(t, day, []string{"1"}, userTimeBlocks, expectedSlots)

	// Caso de prueba: día con múltiples bloques de tiempo que se superponen
	day = "j"
	userTimeBlocks = []models.UserTimeBlock{
		{UserID: "1", TimeBlock: models.TimeBlock{StartMinute: 0, EndMinute: 360}},
		{UserID: "2", TimeBlock: models.TimeBlock{StartMinute: 300, EndMinute: 660}},
	}
	expectedSlots = []models.PlannerEvent{
		{
			DayOfWeek:      "j",
			StartTime:      "00:00",
			EndTime:        "04:59",
			Attendees:      []string{"2"},
			UsersAvailable: 1,
			Duration:       300,
		},
		{
			DayOfWeek:      "j",
			StartTime:      "05:00",
			EndTime:        "06:00",
			Attendees:      []string{},
			UsersAvailable: 0,
			Duration:       61,
		},
		{
			DayOfWeek:      "j",
			StartTime:      "06:01",
			EndTime:        "11:00",
			Attendees:      []string{"1"},
			UsersAvailable: 1,
			Duration:       300,
		},
		{
			DayOfWeek:      "j",
			StartTime:      "11:01",
			EndTime:        "23:59",
			Attendees:      []string{"1", "2"},
			UsersAvailable: 2,
			Duration:       779,
		},
	}
	assertTimeSlots(t, day, []string{"1", "2"}, userTimeBlocks, expectedSlots)

	// Caso de prueba: horario que abarca varios horarios más
	day = "j"
	userTimeBlocks = []models.UserTimeBlock{
		{UserID: "1", TimeBlock: models.TimeBlock{StartMinute: 600, EndMinute: 1200}},
		{UserID: "2", TimeBlock: models.TimeBlock{StartMinute: 800, EndMinute: 900}},
		{UserID: "3", TimeBlock: models.TimeBlock{StartMinute: 1000, EndMinute: 1100}},
	}
	expectedSlots = []models.PlannerEvent{
		{
			DayOfWeek:      "j",
			StartTime:      "00:00",
			EndTime:        "09:59",
			Attendees:      []string{"1", "2", "3"},
			UsersAvailable: 3,
			Duration:       600,
		},
		{
			DayOfWeek:      "j",
			StartTime:      "10:00",
			EndTime:        "13:19",
			Attendees:      []string{"2", "3"},
			UsersAvailable: 2,
			Duration:       200,
		},
		{
			DayOfWeek:      "j",
			StartTime:      "13:20",
			EndTime:        "15:00",
			Attendees:      []string{"3"},
			UsersAvailable: 1,
			Duration:       101,
		},
		{
			DayOfWeek:      "j",
			StartTime:      "15:01",
			EndTime:        "16:39",
			Attendees:      []string{"2", "3"},
			UsersAvailable: 2,
			Duration:       99,
		},
		{
			DayOfWeek:      "j",
			StartTime:      "16:40",
			EndTime:        "18:20",
			Attendees:      []string{"2"},
			UsersAvailable: 1,
			Duration:       101,
		},
		{
			DayOfWeek:      "j",
			StartTime:      "18:21",
			EndTime:        "20:00",
			Attendees:      []string{"2", "3"},
			UsersAvailable: 2,
			Duration:       100,
		},
		{
			DayOfWeek:      "j",
			StartTime:      "20:01",
			EndTime:        "23:59",
			Attendees:      []string{"1", "2", "3"},
			UsersAvailable: 3,
			Duration:       239,
		},
	}
	assertTimeSlots(t, day, []string{"1", "2", "3"}, userTimeBlocks, expectedSlots)

}

// Función auxiliar para aserciones de los resultados
func assertTimeSlots(t *testing.T, day string, usersId []string, userTimeBlocks []models.UserTimeBlock, expectedSlots []models.PlannerEvent) {
	actualSlots := findAvailableTimeSlotsByDay(day, userTimeBlocks, usersId)
	if !reflect.DeepEqual(actualSlots, expectedSlots) {
		t.Errorf("ERROR Expected: %v, Got: %v", expectedSlots, actualSlots)
	} else {
		t.Logf("SUCCESS Expected: %v, Got: %v", expectedSlots, actualSlots)
	}
}
