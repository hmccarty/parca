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

func (event *VerifyOnMessageEvent) Handle(data m.EventData) (m.Response, error) {

	return m.Response{}, nil
}
