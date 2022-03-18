package models

import "time"

type Message struct {
	ID        string
	ChannelID string
	GuildID   string
	Content   string
	Timestamp time.Time
	Author    *User
	Member    *Member
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
