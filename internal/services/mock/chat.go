package mock

import (
	m "github.com/hmccarty/parca/internal/models"
)

type MockChatContext struct {
	contextType m.ChatContextType
	response    m.Response
	guildID     string
	userID      string
	channelID   string
	uniqueID    string
	options     []m.CommandOption
	message     *m.ChatMessage

	GuildNames   map[string]string
	UserNames    map[string]string
	ChannelNames map[string]string
	RoleNames    map[string]string

	ErrGetGuildName   error
	ErrGetUserName    error
	ErrGetChannelName error
	ErrGetRoleName    error
}

func (c *MockChatContext) SetType(contextType m.ChatContextType) {
	c.contextType = contextType
}

func (c *MockChatContext) Type() m.ChatContextType {
	return c.contextType
}

func (c *MockChatContext) Respond(response m.Response) error {
	c.response = response
	return nil
}

func (c *MockChatContext) GetResponse() m.Response {
	return c.response
}

func (c *MockChatContext) SetGuildID(guildID string) {
	c.guildID = guildID
}

func (c *MockChatContext) GuildID() string {
	return c.guildID
}

func (c *MockChatContext) SetUserID(userID string) {
	c.userID = userID
}

func (c *MockChatContext) UserID() string {
	return c.userID
}

func (c *MockChatContext) SetChannelID(channelID string) {
	c.channelID = channelID
}

func (c *MockChatContext) ChannelID() string {
	return c.channelID
}

func (c *MockChatContext) SetUniqueID(uniqueID string) {
	c.uniqueID = uniqueID
}

func (c *MockChatContext) UniqueID() string {
	return c.uniqueID
}

func (c *MockChatContext) SetOptions(options []m.CommandOption) {
	c.options = options
}

func (c *MockChatContext) Options() []m.CommandOption {
	return c.options
}

func (c *MockChatContext) SetMessage(message *m.ChatMessage) {
	c.message = message
}

func (c *MockChatContext) Message() *m.ChatMessage {
	return c.message
}

func (c *MockChatContext) GetGuildNameFromID(guildID string) (string, error) {
	if name, ok := c.GuildNames[guildID]; ok {
		return name, nil
	}
	return "", c.ErrGetGuildName
}

func (c *MockChatContext) GetUserNameFromIDs(userID, guildID string) (string, error) {
	if name, ok := c.UserNames[userID]; ok {
		return name, nil
	}
	return "", c.ErrGetUserName
}

func (c *MockChatContext) GetChannelNameFromIDs(channelID, guildID string) (string, error) {
	if name, ok := c.ChannelNames[channelID]; ok {
		return name, nil
	}
	return "", c.ErrGetChannelName
}

func (c *MockChatContext) GetRoleNameFromIDs(roleID, guildID string) (string, error) {
	if name, ok := c.RoleNames[roleID]; ok {
		return name, nil
	}
	return "", c.ErrGetRoleName
}
