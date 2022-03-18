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
					}
					s.ChannelMessageSendEmbed(i.Message.ChannelID,
						&dg.MessageEmbed{
							Title:       resp.Title,
							Description: resp.Description,
							URL:         resp.URL,
						})
				},
			)
		}
	}
}
