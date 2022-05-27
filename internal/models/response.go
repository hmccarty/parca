package models

type Response struct {
	Type      ResponseType
	GuildID   string
	ChannelID string
	UserID    string
	MessageID string
	Color     int
	IsForm    bool

	// Message responses
	IsEphemeral bool
	URL         string
	Title       string
	Description string
	Buttons     []ResponseButton

	// Form responses
	CustomID string
	Inputs   []ResponseInput

	// Role responses
	RoleID string
}

type ResponseType uint8

const (
	MessageResponse     ResponseType = 0
	AckResponse         ResponseType = 1
	DMResponse          ResponseType = 2
	AddRoleResponse     ResponseType = 3
	RemoveRoleResponse  ResponseType = 4
	MessageEditResponse ResponseType = 5
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

type ResponseInput struct {
	Style    ResponseInputStyle
	Label    string
	Required bool
	CustomID string
}

type ResponseInputStyle uint8

const (
	ShortInputStyle     ResponseInputStyle = 1
	ParagraphInputStyle ResponseInputStyle = 2
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
