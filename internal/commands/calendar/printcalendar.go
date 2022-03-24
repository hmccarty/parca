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

func (*PrintCalendar) Options() []m.CommandOptionMetadata {
	return nil
}

func (cmd *PrintCalendar) Run(ctx m.CommandContext) error {
	client := cmd.createDbClient()
	calendarIDs, err := client.GetCalendars(ctx.ChannelID(), ctx.GuildID())
	if err != nil {
		return ctx.Respond(m.Response{
			Description: "Could not add calendar at this time, try again later",
		})
	}

	desc := ""
	if len(calendarIDs) == 0 {
		desc = "No calendars found"
	} else {
		for _, calendarID := range calendarIDs {
			calendar, err := cmd.calendarClient.GetCalendarData(calendarID)
			if err != nil {
				desc += fmt.Sprintf("- %s (Couldn't retrieve title)\n", calendarID)
			} else {
				desc += fmt.Sprintf("- `%s` \n", calendar.Name)
			}
		}
	}

	channelName, err := ctx.GetChannelNameFromIDs(ctx.ChannelID(), ctx.GuildID())
	if err != nil {
		return err
	}

	return ctx.Respond(m.Response{
		Title:       fmt.Sprintf("Calendars in #%s", channelName),
		Description: desc,
	})
}
