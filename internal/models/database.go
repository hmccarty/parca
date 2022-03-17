package models

import (
	c "github.com/hmccarty/parca/internal/services/config"
)

type DbClient interface {
	// Currency
	GetUserBalance(userid string) (float64, error)
	SetUserBalance(userid, guildid string, amt float64) error
	GetBalancesFromGuild(guildid string) ([]*BalanceEntry, error)

	// Calendar
	AddCalendar(calendarid, channelid, guildid string) error
	GetCalendars(channelid, guildid string) ([]string, error)
	HasCalendar(calendarid, channelid, guildid string) (bool, error)
	RemoveCalendar(calendarid, channelid, guildid string) error
}

type OpenClient func(config *c.Config) (DbClient, error)

type BalanceEntry struct {
	UserID  string
	Balance float64
}
