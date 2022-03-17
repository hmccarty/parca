package calendar

import (
	"fmt"

	m "github.com/hmccarty/parca/internal/models"
)

type PrintCalendar struct {
	createDbClient func() m.DbClient
	calendarClient m.CalendarClient
}

func NewPrintCalendarCommand(createDbClient func() m.DbClient, calendarClient m.CalendarClient) m.Command {
	return &PrintCalendar{
		createDbClient: createDbClient,
		calendarClient: calendarClient,
	}
}

func (*PrintCalendar) Name() string {
	return "printcalendars"
}

func (*PrintCalendar) Description() string {
	return "Lists all calendars active in channel"
}

func (*PrintCalendar) Options() []m.CommandOption {
	return []m.CommandOption{}
}

func (command *PrintCalendar) Run(data m.CommandData, _ []m.CommandOption) m.Response {
	client := command.createDbClient()
	calendarIDs, err := client.GetCalendars(data.ChannelID, data.GuildID)
	if err != nil {
		return m.Response{
			Description: "Could not add calendar at this time, try again later",
		}
	}

	desc := ""
	if len(calendarIDs) == 0 {
		desc = "No calendars found"
	} else {
		for _, calendarID := range calendarIDs {
			calendar, err := command.calendarClient.GetCalendarData(calendarID)
			if err != nil {
				desc += fmt.Sprintf("- %s (Couldn't retrieve title)\n", calendarID)
			} else {
				desc += fmt.Sprintf("- `%s` \n", calendar.Name)
			}
		}
	}

	return m.Response{
		Title:       fmt.Sprintf("Calendars in <#%s>", data.ChannelID),
		Description: desc,
	}
}
