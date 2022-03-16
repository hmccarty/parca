package currency

import (
	"fmt"

	m "github.com/hmccarty/arc-assistant/internal/models"
)

type GetBalance struct {
	createDbClient func() m.DbClient
}

func NewGetBalanceCommand(createDbClient func() m.DbClient) m.Command {
	return &GetBalance{
		createDbClient: createDbClient,
	}
}

func (_ *GetBalance) Name() string {
	return "getbalance"
}

func (_ *GetBalance) Description() string {
	return "Gets the balance of a user"
}

func (_ *GetBalance) Options() []m.CommandOption {
	return []m.CommandOption{
		{
			Name:     "user",
			Type:     m.UserOption,
			Required: true,
		},
	}
}

func (command *GetBalance) Run(data m.CommandData, opts []m.CommandOption) string {
	if len(opts) != 1 {
		return "Invalid number of options"
	}

	userID := opts[0].Value.(string)
	client := command.createDbClient()
	balance, _ := client.GetUserBalance(userID)
	return fmt.Sprintf("You have %.2f in your account", balance)
}
