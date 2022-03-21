package currency

import (
	"fmt"

	m "github.com/hmccarty/parca/internal/models"
)

type Leaderboard struct {
	createDbClient func() m.DbClient
}

func NewLeaderboardCommand(createDbClient func() m.DbClient) m.Command {
	return &Leaderboard{
		createDbClient: createDbClient,
	}
}

func (*Leaderboard) Name() string {
	return "leaderboard"
}

func (*Leaderboard) Description() string {
	return "Prints richest users on server"
}

func (*Leaderboard) Options() []m.CommandOptionMetadata {
	return nil
}

func (cmd *Leaderboard) Run(ctx m.CommandContext) error {
	client := cmd.createDbClient()
	balances, err := client.GetBalancesFromGuild(ctx.GuildID())
	if err != nil {
		return ctx.Respond(m.Response{
			Type:        m.MessageResponse,
			Description: "Everybody is broke",
		})
	}

	var msg string = ""
	for i, balance := range balances {
		msg += fmt.Sprintf("%d. <@%s> has %.2f ARC coins\n",
			i+1, balance.UserID, balance.Balance)
	}
	return ctx.Respond(m.Response{
		Type:        m.MessageResponse,
		Title:       "Leaderboard",
		Description: msg,
	})
}
