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
	return "Sends a token of gratitude to another user (w/o txn fee)"
}

func (*Pay) Options() []m.CommandOptionMetadata {
	return []m.CommandOptionMetadata{
		{
			Name:        "receiver",
			Description: "User you would like to send money to",
			Type:        m.UserOption,
			Required:    true,
		},
		{
			Name:        "amount",
			Description: "Amount you would like to send",
			Type:        m.FloatOption,
			Required:    true,
		},
	}
}

func (cmd *Pay) Run(ctx m.CommandContext) error {
	if len(ctx.Options()) != 2 {
		return m.ErrMissingOptions
	}

	client := cmd.createDbClient()

	senderID := ctx.UserID()
	senderBalance, err := client.GetUserBalance(senderID)

	amount, err := ctx.Options()[1].ToFloat()
	if err != nil {
		return ctx.Respond(m.Response{
			Type:        m.MessageResponse,
			Description: fmt.Sprintf("Failed to get balance of <@%s>", senderID),
			Color:       m.ColorRed,
		})
	} else if senderBalance < amount {
		return ctx.Respond(m.Response{
			Type: m.MessageResponse,
			Description: fmt.Sprintf("Insufficient funds, you have %.2f coins and %.2f are required",
				senderBalance, amount),
			Color: m.ColorRed,
		})
	}

	receiverID, err := ctx.Options()[0].ToUser()
	if err != nil {
		return ctx.Respond(m.Response{
			Type:        m.MessageResponse,
			Description: "Failed to find other user",
			Color:       m.ColorRed,
		})
	} else if senderID == receiverID {
		return ctx.Respond(m.Response{
			Type:        m.MessageResponse,
			Description: "You can't pay yourself",
			Color:       m.ColorRed,
		})
	}

	receiverBalance, err := client.GetUserBalance(receiverID)
	if err != nil {
		return ctx.Respond(m.Response{
			Type:        m.MessageResponse,
			Description: fmt.Sprintf("Failed to get balance of <@%s>", receiverID),
			Color:       m.ColorRed,
		})
	}

	client.SetUserBalance(senderID, ctx.GuildID(), senderBalance-amount)
	client.SetUserBalance(receiverID, ctx.GuildID(), receiverBalance+amount)
	return ctx.Respond(m.Response{
		Type:        m.MessageResponse,
		Description: fmt.Sprintf("Paid <@%s> %.2f ARC coins", receiverID, amount),
		Color:       m.ColorGreen,
	})
}
