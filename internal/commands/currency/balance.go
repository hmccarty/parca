package currency

import (
	"fmt"

	m "github.com/hmccarty/parca/internal/models"
)

type Balance struct {
	createDbClient func() m.DbClient
}

func NewBalanceCommand(createDbClient func() m.DbClient) m.Command {
	return &Balance{
		createDbClient: createDbClient,
	}
}

func (*Balance) Name() string {
	return "balance"
}

func (*Balance) Description() string {
	return "Gets the balance of a user"
}

func (*Balance) Options() []m.CommandOptionMetadata {
	return []m.CommandOptionMetadata{
		{
			Name:        "user",
			Description: "Balance you would like to check",
			Type:        m.UserOption,
			Required:    false,
		},
	}
}

func (cmd *Balance) Run(ctx m.CommandContext) error {
	var userID, userPrompt string
	var err error
	if len(ctx.Options()) == 1 {
		userID, err = ctx.Options()[0].ToUser()
		if err != nil {
			return err
		}

		userName, err := ctx.GetUserNameFromIDs(userID, ctx.GuildID())
		if err != nil {
			return err
		}
		userPrompt = userName + " owns"
	} else {
		userID = ctx.UserID()
		userPrompt = "You own"
	}

	client := cmd.createDbClient()
	balance, _ := client.GetUserBalance(userID)
	return ctx.Respond(m.Response{
		Type: m.MessageResponse,
		Description: fmt.Sprintf("%s **%.2f** ARC coins",
			userPrompt, balance),
	})
}
