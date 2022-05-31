package calendar

import (
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

func (cmd *Today) Run(ctx m.ChatContext) error {
	client := cmd.createDbClient()
	calendarIDs, err := client.GetCalendars(ctx.ChannelID(), ctx.GuildID())
	if err != nil {
		return err
	}

	if len(calendarIDs) == 0 {
		return ctx.Respond(m.Response{
			Type:        m.AckResponse,
			Description: "No calendars added to this channel",
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
			desc += m.ConstructCalendarEventMsg(event)

			// Prevent message overflow
			if len(desc) > 4000 {
				break
			}
		}
	}

	return ctx.Respond(m.Response{
		Type:        m.AckResponse,
		Title:       "Today's events",
		Description: desc,
	})
}
