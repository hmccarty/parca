package calendar

import (
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

func (*Week) Options() []m.CommandOptionMetadata {
	return nil
}

func (cmd *Week) Run(ctx m.CommandContext) error {
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

	endTime := time.Now().UTC().Add(7 * 24 * time.Hour)
	var events []m.CalendarEventData
	for _, calendarID := range calendarIDs {
		calEvents, _ := cmd.calendarClient.GetCalendarEvents(calendarID, endTime)
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
		Type:        m.MessageResponse,
		Title:       "This week's events",
		Description: desc,
		Color:       m.ColorGreen,
	})
}
