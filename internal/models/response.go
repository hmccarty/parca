package models

type Response struct {
	Type      ResponseType
	GuildID   string
	ChannelID string
	UserID    string

	// Message responses
	URL         string
	Title       string
	Description string
	Buttons     []ResponseButton

	// Role responses
	RoleID string
}

type ResponseType uint8

const (
	MessageResponse    ResponseType = 1
	AddRoleResponse    ResponseType = 2
	RemoveRoleResponse ResponseType = 3
)

type ResponseButton struct {
	Style     ResponseButtonStyle
	Label     string
	Emoji     string
	ReactData string
	URL       string
}

type ResponseButtonStyle uint8

const (
	EmojiButtonStyle ResponseButtonStyle = 1
	LinkButtonStyle  ResponseButtonStyle = 2
)
