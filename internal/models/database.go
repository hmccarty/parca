package models

import (
	"errors"

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
	GetVerifyConfig(guildID string) (string, string, error)
	AddVerifyCode(code, userID, guildID string) error
	GetVerifyCode(userID string) (string, string, error)

	// General
	CreatePoll(pollTitle, pollID string) error
	GetPollTitle(pollID string) (string, error)
	AddPollVote(vote bool, pollID, userID string) error
	GetPollVote(pollID string) (int, int, error)

	CreateBounty(messageID, channelID, guildID, bountyID string) error
	GetBountyDesc(bountyID string) (string, error)
}

type OpenClient func(config *c.Config) (DbClient, error)

type BalanceEntry struct {
	UserID  string
	Balance float64
}

var (
	ErrorPollIDAlreadyExists = errors.New("poll already exists with id")
	ErrorPollIDDoesntExists  = errors.New("no poll with id")
	ErrorUserAlreadyVoted    = errors.New("user already voted in poll")
	ErrorUnableToRemoveVoter = errors.New("unable to remove voter")

	ErrorBountyIDAlreadyExists = errors.New("bounty already exists with id")
	ErrorBountyIDDoesntExists  = errors.New("no bounty with id")
)
