package calendar

import (
	m "github.com/hmccarty/parca/internal/models"
)

type RemoveCalendar struct {
	createDbClient func() m.DbClient
	calendarClient m.CalendarClient
}

func NewRemoveCalendarCommand(createDbClient func() m.DbClient, calendarClient m.CalendarClient) m.Command {
	return &RemoveCalendar{
		createDbClient: createDbClient,
		calendarClient: calendarClient,
	}
}

func (*RemoveCalendar) Name() string {
	return "removecalendar"
}

func (*RemoveCalendar) Description() string {
	return "Removes calendar from channel"
}

func (*RemoveCalendar) Options() []m.CommandOptionMetadata {
	return []m.CommandOptionMetadata{
		{
			Name:        "calendar-id",
			Description: "Unique ID of the calendar",
			Type:        m.StringOption,
			Required:    true,
		},
	}
}

func (cmd *RemoveCalendar) Run(ctx m.CommandContext) error {
	if len(ctx.Options()) != 1 {
		return m.ErrMissingOptions
	}

	calendarID, err := ctx.Options()[0].ToString()
	if err != nil {
		return err
	}

	client := cmd.createDbClient()
	hasCalendar, err := client.HasCalendar(calendarID, ctx.ChannelID(), ctx.GuildID())
	if err != nil {
		return err
	} else if !hasCalendar {
		return ctx.Respond(m.Response{
			Type:        m.MessageResponse,
			Description: "Calendar doesn't exist in this channel",
			Color:       m.ColorGreen,
		})
	}

	err = client.RemoveCalendar(calendarID, ctx.ChannelID(), ctx.GuildID())
	if err != nil {
		return err
	}

	return ctx.Respond(m.Response{
		Type:        m.MessageResponse,
		Description: "Removed calendar from channel",
	})
}
