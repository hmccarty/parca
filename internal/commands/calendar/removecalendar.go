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

func (*RemoveCalendar) Options() []m.CommandOption {
	return []m.CommandOption{
		{
			Name:     "calendar-id",
			Type:     m.StringOption,
			Required: true,
		},
	}
}

func (command *RemoveCalendar) Run(data m.CommandData, opts []m.CommandOption) m.Response {
	if len(opts) != 1 {
		return m.Response{
			Description: "Missing `calendar-id` argument",
		}
	}

	calendarID := opts[0].Value.(string)
	client := command.createDbClient()
	hasCalendar, err := client.HasCalendar(calendarID, data.ChannelID, data.GuildID)
	if err != nil {
		return m.Response{
			Description: "Couldn't check calendar at this time, try again later",
		}
	} else if !hasCalendar {
		return m.Response{
			Description: "Calendar doesn't exist in this channel",
		}
	}

	err = client.RemoveCalendar(calendarID, data.ChannelID, data.GuildID)
	if err != nil {
		return m.Response{
			Description: "Couldn't remove calendar at this time, try again later",
		}
	}

	return m.Response{
		Description: "Removed calendar from channel",
	}
}
