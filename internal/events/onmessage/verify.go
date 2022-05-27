package events

import (
	"fmt"

	m "github.com/hmccarty/parca/internal/models"
)

func NewVerifyOnMessageEvent(createDbClient func() m.DbClient) m.Event {
	return &VerifyOnMessageEvent{
		createDbClient: createDbClient,
	}
}

type VerifyOnMessageEvent struct {
	createDbClient func() m.DbClient
}

func (*VerifyOnMessageEvent) GetType() m.EventType {
	return m.OnMessageCreate
}

func (event *VerifyOnMessageEvent) Handle(ctx m.EventContext) error {
	isDM, err := ctx.IsChannelDM(ctx.ChannelID(), ctx.GuildID())
	if err != nil {
		return err
	} else if !isDM {
		return nil
	} else if ctx.Message() == nil {
		return m.ErrMissingData
	}

	client := event.createDbClient()
	code, guildID, err := client.GetVerifyCode(ctx.UserID())
	if err != nil {
		return err
	}

	if code == "" {
		return nil
	} else if ctx.Message().Content != code {
		return ctx.Respond(m.Response{
			Type:        m.MessageResponse,
			Description: "Invalid code",
		})
	}

	_, roleID, err := client.GetVerifyConfig(guildID)
	if err != nil {
		return err
	}

	err = ctx.Respond(m.Response{
		Type:    m.AddRoleResponse,
		GuildID: guildID,
		UserID:  ctx.UserID(),
		RoleID:  roleID,
	})
	if err != nil {
		return err
	}

	guildName, err := ctx.GetGuildNameFromID(guildID)

	return ctx.Respond(m.Response{
		Type:        m.MessageResponse,
		GuildID:     guildID,
		ChannelID:   ctx.ChannelID(),
		Description: fmt.Sprintf("You have been verified on %s", guildName),
	})
}
