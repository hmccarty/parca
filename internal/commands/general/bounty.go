package general

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"

	m "github.com/hmccarty/parca/internal/models"
)

const (
	bountySubmitReactData = "bounty-submit-%s-%s-%s-%s"
	bountyAcceptReactData = "bounty-accept-%s-%s-%s-%s-%s"
	bountyDenyReactData   = "bounty-deny-%s-%s-%s-%s-%s"
)

type Bounty struct {
	createDbClient func() m.DbClient
}

func NewBountyCommand(createDbClient func() m.DbClient) m.Command {
	return &Bounty{
		createDbClient: createDbClient,
	}
}

func (*Bounty) Name() string {
	return "bounty"
}

func (*Bounty) Description() string {
	return "Creates a role menu (max 5 roles per menu)"
}

func (*Bounty) Options() []m.CommandOptionMetadata {
	return []m.CommandOptionMetadata{
		{
			Name:        "title",
			Description: "Title of the bounty",
			Type:        m.StringOption,
			Required:    true,
		},
		{
			Name:        "description",
			Description: "Description of the bounty",
			Type:        m.StringOption,
			Required:    true,
		},
		{
			Name:        "link",
			Description: "Additional link to attach to the bounty",
			Type:        m.StringOption,
			Required:    false,
		},
	}
}

func (cmd *Bounty) Run(ctx m.ChatContext) error {
	if ctx.Options() != nil {
		// If command is called
		opts := ctx.Options()
		title, err := opts[0].ToString()
		if err != nil {
			return err
		}

		desc, err := opts[1].ToString()
		if err != nil {
			return err
		}

		link := ""
		if len(opts) > 2 {
			link, err = opts[2].ToString()
			if err != nil {
				return err
			}
		}

		client := cmd.createDbClient()
		bountyID := fmt.Sprintf("%d", rand.Intn(100000))
		err = client.CreateBounty(bountyID, title, desc, link)
		if err != nil {
			if err == m.ErrorBountyIDAlreadyExists {
				for err == m.ErrorBountyIDAlreadyExists {
					bountyID = fmt.Sprintf("%d", rand.Intn(100000))
					err = client.CreateBounty(bountyID, title, desc, link)
				}
				if err != nil {
					return ctx.Respond(m.Response{
						Type:        m.MessageResponse,
						Description: "Failed to create poll please try again later",
						Color:       m.ColorRed,
					})
				}
			} else {
				return ctx.Respond(m.Response{
					Type:        m.MessageResponse,
					Description: "Failed to create poll please try again later",
					Color:       m.ColorRed,
				})
			}
		}

		button := m.ResponseButton{
			Style: m.PrimaryButtonStyle,
			Label: "Claim",
			ReactData: fmt.Sprintf(bountySubmitReactData, ctx.GuildID(),
				ctx.ChannelID(), ctx.UserID(), bountyID),
		}

		return ctx.Respond(m.Response{
			Type:        m.MessageResponse,
			Title:       fmt.Sprintf("Bounty: %s", title),
			Description: fmt.Sprintf("%s \n", desc),
			Buttons:     []m.ResponseButton{button},
		})
	} else if ctx.Message() != nil {
		msg := strings.Split(ctx.Message().Reaction, "-")
		reactType, reactData := msg[1], msg[2:]

		switch reactType {
		case "submit":
			guildID, channelID, userID, bountyID := reactData[0], reactData[1], reactData[2], reactData[3]

			err := ctx.Respond(m.Response{
				Type:        m.DMResponse,
				UserID:      userID,
				Description: fmt.Sprintf("Did <@%s> complete the bounty 'insert title'?", ctx.UserID()),
				Buttons: []m.ResponseButton{
					{
						Style: m.PrimaryButtonStyle,
						Label: "Yes",
						ReactData: fmt.Sprintf(bountyAcceptReactData, guildID, channelID,
							ctx.UserID(), ctx.Message().ID, bountyID),
					},
					{
						Style: m.SecondaryButtonStyle,
						Label: "No",
						ReactData: fmt.Sprintf(bountyDenyReactData, guildID, channelID,
							ctx.UserID(), ctx.Message().ID, bountyID),
					},
				},
			})
			if err != nil {
				return err
			}

		case "accept":
			guildID, channelID, userID, messageID := reactData[0], reactData[1], reactData[2], reactData[3]
			bountyID := reactData[4]

			client := cmd.createDbClient()
			title, _, _, err := client.GetBounty(bountyID)
			if err != nil {
				return err
			}

			err = ctx.Respond(m.Response{
				Type:        m.MessageEditResponse,
				MessageID:   messageID,
				GuildID:     guildID,
				ChannelID:   channelID,
				Title:       fmt.Sprintf("Bounty: %s", title),
				Description: fmt.Sprintf("Claimed by <@%s>", userID),
			})
			if err != nil {
				return err
			}

			return ctx.Respond(m.Response{
				Type:        m.MessageEditResponse,
				Description: fmt.Sprintf("You confirmed <@%s> as completing bounty: '%s'", userID, title),
			})

		case "deny":
		}

		return nil
	}

	return errors.New("invalid command context")
}