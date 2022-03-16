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
	SubCommandGroupOption                   = 2
	StringOption                            = 3
	IntegerOption                           = 4
	BooleanOption                           = 5
	UserOption                              = 6
	ChannelOption                           = 7
	RoleOption                              = 8
	MentionableOption                       = 9
	NumberOption                            = 10
	AttachmentOption                        = 11
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
