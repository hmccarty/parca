package calendar

import (
	"fmt"
	"time"

	m "github.com/hmccarty/parca/internal/models"
)

type Today struct {
	createDbClient func() m.DbClient
	calendarClient m.CalendarClient
}

func NewTodayCommand(createDbClient func() m.DbClient, calendarClient m.CalendarClient) m.Command {
	return &Today{
		createDbClient: createDbClient,
		calendarClient: calendarClient,
	}
}

func (*Today) Name() string {
	return "today"
}

func (*Today) Description() string {
	return "List today's events for each calendar"
}

func (*Today) Options() []m.CommandOptionMetadata {
	return nil
}

func (cmd *Today) Run(ctx m.CommandContext) error {
	client := cmd.createDbClient()
	calendarIDs, err := client.GetCalendars(ctx.ChannelID(), ctx.GuildID())
	if err != nil {
		return err
	}

	if len(calendarIDs) == 0 {
		return ctx.Respond(m.Response{
			Type:        m.MessageResponse,
			Description: "No calendars added to this channel",
			Color:       m.ColorGreen,
		})
	}

	endTime := time.Now().UTC().Add(24 * time.Hour)
	var events []m.CalendarEventData
	for _, calendarID := range calendarIDs {
		calEvents, err := cmd.calendarClient.GetCalendarEvents(calendarID, endTime)
		if err != nil {
			return err
		}
		events = append(events, calEvents...)
	}

	desc := ""
	if len(events) == 0 {
		desc = "No events found"
	} else {
		for _, event := range events {
			desc += fmt.Sprintf("[%s](%s) \n%d/%d/%d at %d:%d\n%s\n",
				event.Name, event.URL, event.Start.Day(), event.Start.Month(),
				event.Start.Year(), event.Start.Hour(), event.Start.Minute(),
				event.Location,
			)
		}
	}

	return ctx.Respond(m.Response{
		Type:        m.MessageResponse,
		Title:       "Today's events",
		Description: desc,
		Color:       m.ColorGreen,
	})
}
