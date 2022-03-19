package discord

import (
	"log"
	"time"

	dg "github.com/bwmarrin/discordgo"
	m "github.com/hmccarty/parca/internal/models"
	c "github.com/hmccarty/parca/internal/services/config"
)

func NewDiscordClient(config *c.Config, commands []m.Command, events []m.Event) (*DiscordClient, error) {
	client := new(DiscordClient)

	// Setup Discord service to use API key
	session, err := dg.New(config.DiscordToken)
	if err != nil {
		return nil, err
	}
	client.Session = session

	// Setup command applications and event handlers
	applications := make([]*dg.ApplicationCommand, len(commands))

	client.Session.AddHandler()

	setupEventHandlers(client, events)

	// Open connection to Discord service
	err = session.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	// Post applications to Discord
	client.registeredApplications = make([]*dg.ApplicationCommand, len(commands))
	for i, v := range applications {
		cmd, err := client.Session.ApplicationCommandCreate(
			config.DiscordAppID, config.DiscordGuildID, v)
		if err != nil {
			log.Panicf("cannot create '%v' command: %v", v.Name, err)
		}
		client.registeredApplications[i] = cmd
	}

	return client, nil
}

type DiscordClient struct {
	Session                *dg.Session
	registeredApplications []*dg.ApplicationCommand
}

func (d *DiscordClient) ExecHourly(quit chan struct{}) {
	ticker := time.NewTicker(time.Hour)
	for {
		select {
		case <-ticker.C:

		case <-quit:
			ticker.Stop()
			return
		}
	}
}

func (d *DiscordClient) Close() {
	for _, v := range d.registeredApplications {
		err := d.Session.ApplicationCommandDelete(d.Session.State.User.ID, "", v.ID)
		if err != nil {
			log.Panicf("cannot delete '%v' command: %v", v.Name, err)
		}
	}

	d.Session.Close()
}
