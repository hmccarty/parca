package calendar

import (
	"fmt"

	m "github.com/hmccarty/parca/internal/models"
)

type AddCalendar struct {
	createDbClient func() m.DbClient
	calendarClient m.CalendarClient
}

func NewAddCalendarCommand(createDbClient func() m.DbClient, calendarClient m.CalendarClient) m.Command {
	return &AddCalendar{
		createDbClient: createDbClient,
		calendarClient: calendarClient,
	}
}

func (*AddCalendar) Name() string {
	return "addcalendar"
}

func (*AddCalendar) Description() string {
	return "Adds Google calendar to a channel for event reminders"
}

func (*AddCalendar) Options() []m.CommandOption {
	return []m.CommandOption{
		{
			Name:     "calendar-id",
			Type:     m.StringOption,
			Required: true,
		},
	}
}

func (command *AddCalendar) Run(data m.CommandData, opts []m.CommandOption) m.Response {
	if len(opts) != 1 {
		return m.Response{
			Description: "Missing `calendar-id` argument",
		}
	}

	calendarID := opts[0].Value.(string)
	calendarData, err := command.calendarClient.GetCalendarData(calendarID)
	if err != nil {
		return m.Response{
			Description: "Failed to get calendar, missing permissions or bad id",
		}
	}

	client := command.createDbClient()
	err = client.AddCalendar(calendarID, data.ChannelID, data.GuildID)
	if err != nil {
		return m.Response{
			Description: "Could not add calendar at this time, try again later",
		}
	}

	return m.Response{
		Description: fmt.Sprintf("Added '%s' calendar to channel",
			calendarData.Name),
	}
}
