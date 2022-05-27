package models

type ChatClient interface{}

type ChatContext interface {
	Type() ChatContextType
	Respond(Response) error
	GuildID() string
	UserID() string
	ChannelID() string
	UniqueID() string
	Options() []CommandOption // TODO: Move to private data
	Message() *ChatMessage    // TODO: Move to private data

	GetGuildNameFromID(guildID string) (string, error)
	GetUserNameFromIDs(userID, guildID string) (string, error)
	GetChannelNameFromIDs(channelID, guildID string) (string, error)
	GetRoleNameFromIDs(roleID, guildID string) (string, error)
}

type ChatContextType uint8

const (
	CommandCall  ChatContextType = 0
	CommandReply ChatContextType = 1
	// TODO: Absorb event contexts using this type
)

type ChatMessage struct {
	IsDM     bool
	ID       string
	Content  string
	Reaction string
	Values   []string
}
