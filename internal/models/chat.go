package models

type ChatClient interface{}

type ChatContext interface {
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

type ChatMessage struct {
	IsDM     bool
	ID       string
	Content  string
	Reaction string
	Values   []string
}
