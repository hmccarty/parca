package discord

import (
	"fmt"

	dg "github.com/bwmarrin/discordgo"
	m "github.com/hmccarty/parca/internal/models"
)

func setupEventHandlers(session *DiscordSession, events []m.Event) {
	for _, event := range events {
		switch event.GetType() {
		case m.OnMessageCreate:
			session.Session.AddHandler(
				func(s *dg.Session, i *dg.MessageCreate) {
					if i.Author.ID == s.State.User.ID {
						return
					}

					eventData := m.EventData{
						Message: messageFromData(i.Message),
					}
					resp, err := event.Handle(eventData)
					if err != nil {
						fmt.Println(err)
						return
					} else if resp == nil {
						return
					}

					switch resp.Type {
					case m.MessageResponse:
						s.ChannelMessageSendEmbed(i.Message.ChannelID,
							&dg.MessageEmbed{
								Title:       resp.Title,
								Description: resp.Description,
								URL:         resp.URL,
							})
					case m.AddRoleResponse:
						err := s.GuildMemberRoleAdd(resp.GuildID,
							resp.UserID, resp.RoleID)
						if err != nil {
							fmt.Println(err)
						}
					case m.RemoveRoleResponse:
						err := s.GuildMemberRoleRemove(resp.GuildID,
							resp.UserID, resp.RoleID)
						if err != nil {
							fmt.Println(err)
						}
					}
				},
			)
		}
	}
}
