package models

import "errors"

type Event interface {
	GetType() EventType
	Handle(EventContext) error
}

type EventType uint8

const (
	OnMessageCreate EventType = 1
	OnButtonPress   EventType = 2
)

type EventContext interface {
	Respond(Response) error
	GuildID() string
	UserID() string
	ChannelID() string
	Message() *ChatMessage

	GetGuildNameFromID(guildID string) (string, error)
	GetUserNameFromIDs(userID, guildID string) (string, error)
	GetChannelNameFromIDs(channelID, guildID string) (string, error)
	GetRoleNameFromIDs(roleID, guildID string) (string, error)
	IsChannelDM(channelID, guildID string) (bool, error)
}

var (
	ErrMissingData = errors.New("mising required data")
)
