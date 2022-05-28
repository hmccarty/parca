package currency

import (
	"fmt"

	m "github.com/hmccarty/parca/internal/models"
)

type SetBalance struct {
	modIDs         []string
	createDbClient func() m.DbClient
}

func NewSetBalanceCommand(modIDs []string, createDbClient func() m.DbClient) m.Command {
	return &SetBalance{
		modIDs:         modIDs,
		createDbClient: createDbClient,
	}
}

func (*SetBalance) Name() string {
	return "setbalance"
}

func (*SetBalance) Description() string {
	return "Sets the balance of a user"
}

func (*SetBalance) Options() []m.CommandOptionMetadata {
	return []m.CommandOptionMetadata{
		{
			Name:        "user",
			Description: "User you would like to change the balance of",
			Type:        m.UserOption,
			Required:    true,
		},
		{
			Name:        "amount",
			Description: "Amount to change balance to",
			Type:        m.FloatOption,
			Required:    true,
		},
	}
}

func (cmd *SetBalance) Run(ctx m.ChatContext) error {
	if len(ctx.Options()) != 2 {
		return m.ErrMissingOptions
	}

	isMod := false
	for _, modID := range cmd.modIDs {
		if ctx.UserID() == modID {
			isMod = true
			break
		}
	}
	if !isMod {
		return ctx.Respond(m.Response{
			Type:        m.AckResponse,
			Description: "Only bot moderators can use this command",
			Color:       m.ColorRed,
		})
	}

	userID, err := ctx.Options()[0].ToUser()
	if err != nil {
		return err
	}

	amount, err := ctx.Options()[1].ToFloat()
	if err != nil {
		return err
	}

	client := cmd.createDbClient()
	client.SetUserBalance(userID, ctx.GuildID(), amount)
	balance, _ := client.GetUserBalance(userID)
	return ctx.Respond(m.Response{
		Type: m.AckResponse,
		Description: fmt.Sprintf("Set <@%s>'s balance to %.2f",
			userID, balance),
		Color: m.ColorGreen,
	})
}
