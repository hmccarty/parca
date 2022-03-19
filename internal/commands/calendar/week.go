package calendar

import (
	"fmt"
	"time"

	m "github.com/hmccarty/parca/internal/models"
)

type Week struct {
	createDbClient func() m.DbClient
	calendarClient m.CalendarClient
}

func NewWeekCommand(createDbClient func() m.DbClient, calendarClient m.CalendarClient) m.Command {
	return &Week{
		createDbClient: createDbClient,
		calendarClient: calendarClient,
	}
}

func (*Week) Name() string {
	return "week"
}

func (*Week) Description() string {
	return "List this week's events for each calendar"
}

func (*Week) Options() []m.CommandOption {
	return []m.CommandOption{}
}

func (command *Week) Run(data m.CommandData, _ []m.CommandOption) m.Response {
	client := command.createDbClient()
	calendarIDs, err := client.GetCalendars(data.ChannelID, data.GuildID)
	if err != nil {
		return m.Response{
			Description: "Could not add calendar at this time, try again later",
		}
	}

	if len(calendarIDs) == 0 {
		return m.Response{
			Description: "No calendars added to this channel",
		}
	}

	endTime := time.Now().UTC().Add(7 * 24 * time.Hour)
	var events []m.CalendarEventData
	for _, calendarID := range calendarIDs {
		calEvents, _ := command.calendarClient.GetCalendarEvents(calendarID, endTime)
		events = append(events, calEvents...)
	}

	desc := ""
	if len(events) == 0 {
		desc = "No events found"
	} else {
		for _, event := range events {
			desc += fmt.Sprintf("%s\n%s\n", event.Name, event.Location)
		}
	}

	return m.Response{
		Title:       "Week's events",
		Description: desc,
	}
}

func (*Week) HandleReaction(data m.CommandData, reaction string) m.Response {
	return m.Response{
		Description: "Not expecting a reaction",
	}
}
