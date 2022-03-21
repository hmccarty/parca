package discord

import (
	"fmt"
	"strings"

	dg "github.com/bwmarrin/discordgo"
	m "github.com/hmccarty/parca/internal/models"
)

func setCmdHandler(client *DiscordClient, cmds []m.Command) error {
	cmdHandlers := make(map[string]m.Command, len(cmds))
	for _, cmd := range cmds {
		cmdHandlers[cmd.Name()] = cmd
	}

	client.Session.AddHandler(
		func(s *dg.Session, i *dg.InteractionCreate) {
			var name string
			switch i.Type {
			case dg.InteractionApplicationCommand:
				name = i.ApplicationCommandData().Name
			case dg.InteractionMessageComponent:
				// TODO: Error check
				name = strings.Split(i.MessageComponentData().CustomID, "-")[0]
			}

			if cmd, ok := cmdHandlers[name]; ok {
				ctx, err := createCmdCtx(s, i.Interaction, cmd)
				if err != nil {
					fmt.Println("could not create command context")
				} else {
					err = cmd.Run(ctx)
					if err != nil {
						fmt.Println(err)
					}
				}
			} else {
				fmt.Println("command could not be found")
			}
		},
	)

	return nil
}

func setEventHandlers(client *DiscordClient, events []m.Event) error {
	for _, event := range events {
		switch event.GetType() {
		case m.OnMessageCreate:
			client.Session.AddHandler(
				func(s *dg.Session, e *dg.MessageCreate) {
					if e.Author.ID == s.State.User.ID {
						return
					}

					eventCtx := &DiscordEventContext{
						session:   client.Session,
						guildID:   e.GuildID,
						userID:    e.Author.ID,
						channelID: e.ChannelID,
						messageID: e.Message.ID,
						message: &m.ChatMessage{
							ID:      e.Message.ID,
							Content: e.Message.Content,
						},
					}
					event.Handle(eventCtx)
				},
			)
		}
	}

	return nil
}
