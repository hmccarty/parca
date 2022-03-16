package models

import (
	c "github.com/hmccarty/parca/internal/services/config"
)

type DbClient interface {
	GetUserBalance(string) (float64, error)
	SetUserBalance(string, string, float64) error
	GetBalancesFromGuild(string) ([]*BalanceEntry, error)
}

type OpenClient func(config *c.Config) (DbClient, error)

type BalanceEntry struct {
	UserID  string
	Balance float64
}
