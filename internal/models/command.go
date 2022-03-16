package models

type Command interface {
	Name() string
	Description() string
	Options() []CommandOption
	Run(data CommandData, opts []CommandOption) string
}

type CommandData struct {
	GuildID   string
	ChannelID string
	User      *User
	Member    *Member
}

type CommandOptionType uint8

const (
	SubCommandOption      CommandOptionType = 1
	SubCommandGroupOption CommandOptionType = 2
	StringOption          CommandOptionType = 3
	IntegerOption         CommandOptionType = 4
	BooleanOption         CommandOptionType = 5
	UserOption            CommandOptionType = 6
	ChannelOption         CommandOptionType = 7
	RoleOption            CommandOptionType = 8
	MentionableOption     CommandOptionType = 9
	NumberOption          CommandOptionType = 10
	AttachmentOption      CommandOptionType = 11
)

type CommandOption struct {
	Name     string
	Type     CommandOptionType
	Required bool
	Value    interface{}
	Options  []*CommandOption
}

type User struct {
	ID       string
	Email    string
	Username string
}

type Member struct {
	GuildID string
	User    *User
	Roles   []string
}

type Role struct {
	ID   string
	Name string
}
