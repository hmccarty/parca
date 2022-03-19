package discord

import (
	"log"
	"strings"

	dg "github.com/bwmarrin/discordgo"
	m "github.com/hmccarty/parca/internal/models"
	c "github.com/hmccarty/parca/internal/services/config"
)

func NewDiscordSession(config *c.Config, commands []m.Command, events []m.Event) (*DiscordSession, error) {
	discordSession := new(DiscordSession)

	// CSetup Discord service to use API key
	session, err := dg.New(config.DiscordToken)
	if err != nil {
		return nil, err
	}
	discordSession.Session = session

	// Setup commands as applications and assign their handlers
	applications := make([]*dg.ApplicationCommand, len(commands))
	commandHandlers := map[string]CommandHandler{}
	reactionHandlers := map[string]ReactionHandler{}
	for i, v := range commands {
		applications[i] = appFromCommand(v)
		commandHandlers[v.Name()] = handlerFromCommand(v)
		reactionHandlers[v.Name()] = reactionHandlerFromCommand(v)
	}

	discordSession.Session.AddHandler(func(s *dg.Session, i *dg.InteractionCreate) {
		switch i.Type {
		case dg.InteractionApplicationCommand:
			if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}
		case dg.InteractionMessageComponent:
			key := strings.Split(i.MessageComponentData().CustomID, "-")[0]
			if h, ok := reactionHandlers[key]; ok {
				h(s, i)
			}
		}
	})

	setupEventHandlers(discordSession, events)

	// Open connection to Discord service
	err = session.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	// Post applications to Discord
	discordSession.registeredApplications = make([]*dg.ApplicationCommand, len(commands))
	for i, v := range applications {
		cmd, err := discordSession.Session.ApplicationCommandCreate(
			config.DiscordAppID, config.DiscordGuildID, v)
		if err != nil {
			log.Panicf("cannot create '%v' command: %v", v.Name, err)
		}
		discordSession.registeredApplications[i] = cmd
	}

	return discordSession, nil
}

type DiscordSession struct {
	Session                *dg.Session
	registeredApplications []*dg.ApplicationCommand
}

func (d *DiscordSession) Close() {
	for _, v := range d.registeredApplications {
		err := d.Session.ApplicationCommandDelete(d.Session.State.User.ID, "", v.ID)
		if err != nil {
			log.Panicf("cannot delete '%v' command: %v", v.Name, err)
		}
	}

	d.Session.Close()
}
