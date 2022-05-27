package discord

import (
	"log"
	"time"

	dg "github.com/bwmarrin/discordgo"
	m "github.com/hmccarty/parca/internal/models"
	c "github.com/hmccarty/parca/internal/services/config"
)

func NewDiscordClient(config *c.Config, cmds []m.Command, events []m.Event) (*DiscordClient, error) {
	client := new(DiscordClient)

	// Setup Discord service to use API key
	session, err := dg.New(config.DiscordToken)
	if err != nil {
		return nil, err
	}
	client.Session = session

	// Setup handlers
	setCmdHandler(client, cmds)
	setEventHandlers(client, events)

	// Open connection to Discord service
	err = session.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	// Post applications to Discord
	client.registeredApplications = make([]*dg.ApplicationCommand, len(cmds))
	for i, cmd := range cmds {
		app, err := cmdToApp(cmd)
		if err != nil {
			return nil, err
		}

		registeredApp, err := client.Session.ApplicationCommandCreate(
			config.DiscordAppID, config.DiscordGuildID, app)
		if err != nil {
			return nil, err
		}

		client.registeredApplications[i] = registeredApp
	}

	return client, nil
}

type DiscordClient struct {
	Session                *dg.Session
	registeredApplications []*dg.ApplicationCommand
}

func (d *DiscordClient) ExecPeriodically(p time.Duration, quit chan struct{}) {
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
