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

func (*AddCalendar) Options() []m.CommandOptionMetadata {
	return []m.CommandOptionMetadata{
		{
			Name:        "calendar-id",
			Description: "Unique ID of the calendar",
			Type:        m.StringOption,
			Required:    true,
		},
	}
}

func (cmd *AddCalendar) Run(ctx m.ChatContext) error {
	if len(ctx.Options()) != 1 {
		return m.ErrMissingOptions
	}

	calendarID, err := ctx.Options()[0].ToString()
	if err != nil {
		return err
	}

	calendarData, err := cmd.calendarClient.GetCalendarData(calendarID)
	if err != nil {
		return ctx.Respond(m.Response{
			Type:        m.AckResponse,
			Description: "Failed to get calendar, missing permissions or bad id",
		})
	}

	client := cmd.createDbClient()
	err = client.AddCalendar(calendarID, ctx.ChannelID(), ctx.GuildID())
	if err != nil {
		return ctx.Respond(m.Response{
			Type:        m.AckResponse,
			Description: "Could not add calendar at this time, try again later",
		})
	}

	return ctx.Respond(m.Response{
		Type: m.AckResponse,
		Description: fmt.Sprintf("Added '%s' calendar to channel",
			calendarData.Name),
	})
}
