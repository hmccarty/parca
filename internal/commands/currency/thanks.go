package currency

import (
	"fmt"

	m "github.com/hmccarty/parca/internal/models"
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

func (*Thanks) Options() []m.CommandOptionMetadata {
	return []m.CommandOptionMetadata{
		{
			Name:        "receiver",
			Description: "User you would like to thank",
			Type:        m.UserOption,
			Required:    true,
		},
	}
}

func (cmd *Thanks) Run(ctx m.ChatContext) error {
	if len(ctx.Options()) == 2 {
		return m.ErrMissingOptions
	}

	client := cmd.createDbClient()

	senderID := ctx.UserID()
	senderBalance, err := client.GetUserBalance(senderID)
	if err != nil {
		return ctx.Respond(m.Response{
			Type:        m.AckResponse,
			Description: fmt.Sprintf("Failed to get balance of <@%s>", senderID),
			Color:       m.ColorRed,
		})
	} else if senderBalance < txnFee {
		return ctx.Respond(m.Response{
			Type: m.AckResponse,
			Description: fmt.Sprintf("Insufficient funds, you have %.2f coins and %.2f are required",
				senderBalance, txnFee),
			Color: m.ColorRed,
		})
	}

	receiverID, err := ctx.Options()[0].ToUser()
	if err != nil {
		return err
	}

	if senderID == receiverID {
		return ctx.Respond(m.Response{
			Type:        m.AckResponse,
			Description: "You can't thank yourself",
			Color:       m.ColorRed,
		})
	}

	receiverBalance, err := client.GetUserBalance(receiverID)
	if err != nil {
		return ctx.Respond(m.Response{
			Type:        m.AckResponse,
			Description: fmt.Sprintf("Failed to get balance of <@%s>", receiverID),
			Color:       m.ColorRed,
		})
	}

	client.SetUserBalance(senderID, ctx.GuildID(), senderBalance-txnFee)
	client.SetUserBalance(receiverID, ctx.GuildID(), receiverBalance+amount)
	return ctx.Respond(m.Response{
		Type: m.AckResponse,
		Description: fmt.Sprintf("<@%s> thanked <@%s> with %.2f ARC coins",
			senderID, receiverID, amount),
		Color: m.ColorGreen,
	})
}
