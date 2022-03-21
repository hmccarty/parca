package models

type Response struct {
	Type      ResponseType
	GuildID   string
	ChannelID string
	UserID    string
	Color     int

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
	MessageResponse     ResponseType = 0
	DMAuthorResponse    ResponseType = 1
	AddRoleResponse     ResponseType = 2
	RemoveRoleResponse  ResponseType = 3
	MessageEditResponse ResponseType = 4
)

type ResponseButton struct {
	Style     ResponseButtonStyle
	Label     string
	Emoji     Emoji
	ReactData string
	URL       string
}

type ResponseButtonStyle uint8

const (
	PrimaryButtonStyle   ResponseButtonStyle = 1
	SecondaryButtonStyle ResponseButtonStyle = 2
	LinkButtonStyle      ResponseButtonStyle = 3
)

type Emoji uint8

const (
	ThumbsUpEmoji   Emoji = 1
	ThumbsDownEmoji Emoji = 2
)

const (
	ColorRed   = 0xc41010
	ColorGreen = 0x207002
)
