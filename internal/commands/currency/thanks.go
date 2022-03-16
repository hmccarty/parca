package currency

import (
	"fmt"

	m "github.com/hmccarty/arc-assistant/internal/models"
)

const (
	amount float64 = 5.0
	txnFee         = 0.5
)

type Thanks struct {
	createDbClient func() m.DbClient
}

func NewThanksCommand(createDbClient func() m.DbClient) m.Command {
	return &Thanks{
		createDbClient: createDbClient,
	}
}

func (*Thanks) Name() string {
	return "thanks"
}

func (*Thanks) Description() string {
	return "Sends a token of gratitude to another user (with txn fee)"
}

func (*Thanks) Options() []m.CommandOption {
	return []m.CommandOption{
		{
			Name:     "receiver",
			Type:     m.UserOption,
			Required: true,
		},
	}
}

func (command *Thanks) Run(data m.CommandData, opts []m.CommandOption) string {
	if len(opts) != 1 {
		return "Invalid number of options"
	}

	client := command.createDbClient()

	var senderID string
	if data.User == nil && data.Member == nil {
		return "You must be logged in to use this command"
	} else if data.User != nil {
		senderID = data.User.ID
	} else {
		senderID = data.Member.User.ID
	}

	senderBalance, err := client.GetUserBalance(senderID)
	if err != nil {
		return fmt.Sprintf("Failed to get balance of <@%s>", senderID)
	} else if senderBalance < txnFee {
		return fmt.Sprintf("Insufficient funds, you have %.2f coins and %.2f are required",
			senderBalance, txnFee)
	}

	receiverID := opts[0].Value.(string)
	if senderID == receiverID {
		return "You can't thank yourself"
	}

	receiverBalance, err := client.GetUserBalance(receiverID)
	if err != nil {
		return fmt.Sprintf("Failed to get balance of <@%s>", receiverID)
	}

	client.SetUserBalance(senderID, data.GuildID, senderBalance-txnFee)
	client.SetUserBalance(receiverID, data.GuildID, receiverBalance+amount)
	return fmt.Sprintf("Send <@%s> %.2f ARC coins", receiverID, amount)
}
