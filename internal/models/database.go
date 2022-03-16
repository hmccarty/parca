package models

import (
	c "github.com/hmccarty/arc-assistant/internal/services/config"
)

/*
	user:
		<userid>:
			- username: str
			- balance: float

	verify:
		<guildid>:
			- domain: str
			- role: roleid
			<userid>:
				- code: int

	arcdle:
		<userid>:
			- channel: channelid
			- message: messageid
			- status: int
			- hidden: str
			- visible: str

	daily: [userid]

	backlog: [str]

	calendar:
		<guildid>:
			[
				- channel: [channelid]
				- calendar: [calendarid]
			]

	bounty:
		<guildid>:
			[
				title: str
				user: userid
				guild: guildid
				channel: channelid
				message: messageid
				amt: float
			]
*/

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
