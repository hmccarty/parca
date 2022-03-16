package currency

import (
	"fmt"

	m "github.com/hmccarty/arc-assistant/internal/models"
)

type Leaderboard struct {
	createDbClient func() m.DbClient
}

func NewLeaderboardCommand(createDbClient func() m.DbClient) m.Command {
	return &Leaderboard{
		createDbClient: createDbClient,
	}
}

func (_ *Leaderboard) Name() string {
	return "leaderboard"
}

func (_ *Leaderboard) Description() string {
	return "Prints richest users on server"
}

func (_ *Leaderboard) Options() []m.CommandOption {
	return []m.CommandOption{}
}

func (command *Leaderboard) Run(data m.CommandData, _ []m.CommandOption) string {
	client := command.createDbClient()
	balances, err := client.GetBalancesFromGuild(data.GuildID)
	if err != nil {
		return "Failed"
	}

	var msg string = ""
	for i, balance := range balances {
		msg += fmt.Sprintf("%d. <@%s> has %.2f coins\n",
			i+1, balance.UserID, balance.Balance)
	}
	return msg
}
