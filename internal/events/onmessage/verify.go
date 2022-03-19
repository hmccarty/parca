package events

import (
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

func (event *VerifyOnMessageEvent) Handle(data m.EventData) (*m.Response, error) {
	client := event.createDbClient()
	code, guildID, _ := client.GetVerifyCode(data.Message.Author.ID)

	if code == "" {
		return nil, nil
	} else if data.Message.Content != code {
		return &m.Response{
			Description: "Invalid code",
		}, nil
	}

	_, roleID, _ := client.GetVerifyConfig(guildID)

	return &m.Response{
		Type:    m.AddRoleResponse,
		GuildID: guildID,
		UserID:  data.Message.Author.ID,
		RoleID:  roleID,
	}, nil
}
