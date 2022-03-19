package currency

import (
	"fmt"

	m "github.com/hmccarty/parca/internal/models"
)

type Pay struct {
	createDbClient func() m.DbClient
}

func NewPayCommand(createDbClient func() m.DbClient) m.Command {
	return &Pay{
		createDbClient: createDbClient,
	}
}

func (*Pay) Name() string {
	return "pay"
}

func (*Pay) Description() string {
	return "Sends a token of gratitude to another user (with txn fee)"
}

func (*Pay) Options() []m.CommandOption {
	return []m.CommandOption{
		{
			Name:     "receiver",
			Type:     m.UserOption,
			Required: true,
		},
		{
			Name:     "amount",
			Type:     m.NumberOption,
			Required: true,
		},
	}
}

func (command *Pay) Run(data m.CommandData, opts []m.CommandOption) m.Response {
	if len(opts) != 2 {
		return m.Response{
			Description: "Invalid number of options",
		}
	}

	client := command.createDbClient()

	var senderID string
	if data.User == nil && data.Member == nil {
		return m.Response{
			Description: "You must be logged in to use this command",
		}
	} else if data.User != nil {
		senderID = data.User.ID
	} else {
		senderID = data.Member.User.ID
	}

	senderBalance, err := client.GetUserBalance(senderID)
	amount := opts[1].Value.(float64)
	if err != nil {
		return m.Response{
			Description: fmt.Sprintf("Failed to get balance of <@%s>", senderID),
		}
	} else if senderBalance < amount {
		return m.Response{
			Description: fmt.Sprintf("Insufficient funds, you have %.2f coins and %.2f are required",
				senderBalance, amount),
		}
	}

	receiverID := opts[0].Value.(string)
	if senderID == receiverID {
		return m.Response{
			Description: "You can't thank yourself",
		}
	}

	receiverBalance, err := client.GetUserBalance(receiverID)
	if err != nil {
		return m.Response{
			Description: fmt.Sprintf("Failed to get balance of <@%s>", receiverID),
		}
	}

	client.SetUserBalance(senderID, data.GuildID, senderBalance-amount)
	client.SetUserBalance(receiverID, data.GuildID, receiverBalance+amount)
	return m.Response{
		Description: fmt.Sprintf("Paid <@%s> %.2f ARC coins", receiverID, amount),
	}
}

func (*Pay) HandleReaction(data m.CommandData, reaction string) m.Response {
	return m.Response{
		Description: "Not expecting a reaction",
	}
}
