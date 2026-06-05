package main

import (
	"html/template"
	"net/http"
)

type WeeklyScheduleComponentData struct {
	Schedule *Schedule
	CanWrite bool
}

func RenderWeeklySchedule(w http.ResponseWriter, schedule *Schedule, canWrite bool) error {
	tmpl, err := template.ParseFiles("./templates/weekly_schedule.html")
	if err != nil {
		return err
	}

	data := WeeklyScheduleComponentData{
		Schedule: schedule,
		CanWrite: canWrite,
	}

	return tmpl.ExecuteTemplate(w, "weekly_schedule", data)
}
