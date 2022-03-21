package models

import "time"

type CalendarClient interface {
	GetCalendarData(string) (CalendarData, error)
	GetCalendarEvents(string, time.Time) ([]CalendarEventData, error)
}

type CalendarData struct {
	Name        string
	URL         string
	Description string
}

type CalendarEventData struct {
	Name     string
	Location string
	URL      string
	Start    time.Time
	End      time.Time
}
