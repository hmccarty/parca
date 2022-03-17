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

func (*Balance) Options() []m.CommandOption {
	return []m.CommandOption{
		{
			Name:     "user",
			Type:     m.UserOption,
			Required: false,
		},
	}
}

func (command *Balance) Run(data m.CommandData, opts []m.CommandOption) m.Response {
	var userID string
	if len(opts) == 1 {
		userID = opts[0].Value.(string)
	} else if data.User == nil && data.Member == nil {
		return m.Response{
			Description: "You must be logged in to use this command",
		}
	} else if data.User != nil {
		userID = data.User.ID
	} else {
		userID = data.Member.User.ID
	}

	client := command.createDbClient()
	balance, _ := client.GetUserBalance(userID)
	return m.Response{
		Description: fmt.Sprintf("You have %.2f in your account", balance),
	}
}
