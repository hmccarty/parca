package models

import (
	"fmt"
	"time"
)

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

func ConstructCalendarEventMsg(e CalendarEventData) string {
	timeMsg := e.Start.Format("Monday (Jan 02) ")
	if e.Start.Day() == e.End.Day() {
		timeMsg += e.Start.Format("from 03:04 PM ")
		timeMsg += e.End.Format("to 03:04 PM")
	} else {
		timeMsg += e.End.Format("to Monday (Jan 02) ")
	}

	msg := fmt.Sprintf("[%s](%s)\n %s \n", e.Name, e.URL, timeMsg)
	if e.Location != "" {
		msg += fmt.Sprintf("%s\n", e.Location)
	}
	return msg + "\n"
}
