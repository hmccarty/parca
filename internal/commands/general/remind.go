package general

import (
	"time"

	m "github.com/hmccarty/parca/internal/models"
)

type Remind struct{}

func NewRemindCommand() m.Command {
	return &Remind{}
}

func (*Remind) Name() string {
	return "remind"
}

func (*Remind) Description() string {
	return "Creates a reminder"
}

func (*Remind) Options() []m.CommandOptionMetadata {
	return []m.CommandOptionMetadata{}
}

func (*Remind) Run(ctx m.ChatContext) error {
	err := ctx.Respond(m.Response{
		Type:        m.MessageResponse,
		Description: "Better than ever",
	})
	if err != nil {
		return err
	}

	time.Sleep(20 * time.Second)
	return ctx.Respond(m.Response{
		GuildID:     ctx.GuildID(),
		ChannelID:   ctx.ChannelID(),
		Type:        m.MessageResponse,
		Description: "Reminder",
	})
}
