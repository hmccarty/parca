package models

import (
	c "github.com/hmccarty/parca/internal/services/config"
)

type DbClient interface {
	// Currency
	GetUserBalance(userID string) (float64, error)
	SetUserBalance(userID, guildID string, amt float64) error
	GetBalancesFromGuild(guildID string) ([]*BalanceEntry, error)

	// Calendar
	AddCalendar(calendarID, channelID, guildID string) error
	GetCalendars(channelID, guildID string) ([]string, error)
	HasCalendar(calendarID, channelID, guildID string) (bool, error)
	RemoveCalendar(calendarID, channelID, guildID string) error

	// Verification
	AddVerifyConfig(domain, roleID, guildID string) error
	AddVerifyCode(code, userID, guildID string) error
	GetVerifyCode(userID, guildID string) (string, error)
}

type OpenClient func(config *c.Config) (DbClient, error)

type BalanceEntry struct {
	UserID  string
	Balance float64
}
