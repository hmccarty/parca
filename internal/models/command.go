package models

import "errors"

var (
	ErrWrongOptionType = errors.New("option cannot be casted this type")
)

type Command interface {
	Name() string
	Description() string
	Options() []CommandOptionMetadata
	Run(ctx CommandContext) error
}

type CommandContext interface {
	Respond(Response) error
	GuildID() string
	UserID() string
	ChannelID() string
	Options() []CommandOption
	Message() *ChatMessage

	GetGuildNameFromID(guildID string) (string, error)
	GetUserNameFromIDs(userID, guildID string) (string, error)
	GetChannelNameFromIDs(channelID, guildID string) (string, error)
	GetRoleNameFromIDs(roleID, guildID string) (string, error)
}

type CommandOptionMetadata struct {
	Type        CommandOptionType
	Name        string
	Description string
	Required    bool
}

type CommandOption struct {
	Metadata CommandOptionMetadata
	Value    interface{}
}

func (c CommandOption) ToString() (string, error) {
	if c.Metadata.Type != StringOption {
		return "", ErrWrongOptionType
	}
	return c.Value.(string), nil
}

func (c CommandOption) ToInteger() (int64, error) {
	if c.Metadata.Type != IntegerOption {
		return 0, ErrWrongOptionType
	}
	return c.Value.(int64), nil
}

func (c CommandOption) ToFloat() (float64, error) {
	if c.Metadata.Type != FloatOption {
		return 0, ErrWrongOptionType
	}
	return c.Value.(float64), nil
}

func (c CommandOption) ToBoolean() (bool, error) {
	if c.Metadata.Type != BooleanOption {
		return false, ErrWrongOptionType
	}
	return c.Value.(bool), nil
}

func (c CommandOption) ToUser() (string, error) {
	if c.Metadata.Type != UserOption {
		return "", ErrWrongOptionType
	}
	return c.Value.(string), nil
}

func (c CommandOption) ToChannel() (string, error) {
	if c.Metadata.Type != ChannelOption {
		return "", ErrWrongOptionType
	}
	return c.Value.(string), nil
}

func (c CommandOption) ToRole() (string, error) {
	if c.Metadata.Type != RoleOption {
		return "", ErrWrongOptionType
	}
	return c.Value.(string), nil
}

type CommandOptionType uint8

const (
	StringOption  CommandOptionType = 1
	IntegerOption CommandOptionType = 2
	FloatOption   CommandOptionType = 3
	BooleanOption CommandOptionType = 4
	UserOption    CommandOptionType = 5
	ChannelOption CommandOptionType = 6
	RoleOption    CommandOptionType = 7
)

var (
	ErrMissingOptions = errors.New("mising required options")
)
