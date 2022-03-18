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

	// Role responses
	RoleID string
}

type ResponseType uint8

const (
	MessageResponse    ResponseType = 1
	AddRoleResponse    ResponseType = 2
	RemoveRoleResponse ResponseType = 3
)
