package discord

import (
	"log"

	dg "github.com/bwmarrin/discordgo"
	m "github.com/hmccarty/parca/internal/models"
	c "github.com/hmccarty/parca/internal/services/config"
)

type DiscordHandler func(*dg.Session, *dg.InteractionCreate)

type DiscordSession struct {
	Session            *dg.Session
	registeredCommands []*dg.ApplicationCommand
}

func NewDiscordSession(config *c.Config, commands []m.Command) (*DiscordSession, error) {
	discordSession := new(DiscordSession)

	session, err := dg.New(config.DiscordToken)
	if err != nil {
		return nil, err
	}
	discordSession.Session = session

	discordCommands := make([]*dg.ApplicationCommand, len(commands))
	discordHandlers := map[string]DiscordHandler{}
	for i, v := range commands {
		discordCommands[i] = appFromCommand(v)
		discordHandlers[v.Name()] = createDiscordHandler(v)
	}

	discordSession.Session.AddHandler(func(s *dg.Session, i *dg.InteractionCreate) {
		if h, ok := discordHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	err = session.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	discordSession.registeredCommands = make([]*dg.ApplicationCommand, len(commands))
	for i, v := range discordCommands {
		cmd, err := discordSession.Session.ApplicationCommandCreate(
			config.DiscordAppID, config.DiscordGuildID, v)
		if err != nil {
			log.Panicf("cannot create '%v' command: %v", v.Name, err)
		}
		discordSession.registeredCommands[i] = cmd
	}

	return discordSession, nil
}

func (d *DiscordSession) Close() {
	for _, v := range d.registeredCommands {
		err := d.Session.ApplicationCommandDelete(d.Session.State.User.ID, "", v.ID)
		if err != nil {
			log.Panicf("cannot delete '%v' command: %v", v.Name, err)
		}
	}

	d.Session.Close()
}

func createDiscordHandler(command m.Command) DiscordHandler {
	return func(s *dg.Session, i *dg.InteractionCreate) {
		data, _ := dataFromInteraction(i.Interaction)
		appData := i.ApplicationCommandData()
		options := make([]m.CommandOption, len(appData.Options))
		for i, v := range appData.Options {
			option, err := optionFromInteraction(v)
			if err != nil {
				log.Println(err)
			}
			options[i] = option
		}

		content := command.Run(data, options)
		s.InteractionRespond(i.Interaction, &dg.InteractionResponse{
			Type: dg.InteractionResponseChannelMessageWithSource,
			Data: &dg.InteractionResponseData{
				Content: content,
			},
		})
	}
}
