package general

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"

	m "github.com/hmccarty/parca/internal/models"
)

const (
	bountySubmitReactData = "bounty-submit-%s"
	bountyAcceptReactData = "bounty-accept-%s"
	bountyDenyReactData   = "bounty-deny-%s"
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

func (cmd *Bounty) Run(ctx m.CommandContext) error {
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

		_ = link

		bountyID := fmt.Sprintf("%d", rand.Intn(100000))

		client := cmd.createDbClient()
		err = client.CreateBounty(desc, bountyID)
		if err != nil {
			if err == m.ErrorBountyIDAlreadyExists {
				for err == m.ErrorBountyIDAlreadyExists {
					bountyID = fmt.Sprintf("%d", rand.Intn(100000))
					err = client.CreateBounty(desc, bountyID)
				}
				if err != nil {
					return ctx.Respond(m.Response{
						Type:        m.MessageResponse,
						Description: "Failed to create bounty please try again later",
					})
				}
			} else {
				return ctx.Respond(m.Response{
					Type:        m.MessageResponse,
					Description: "Failed to create bounty please try again later",
				})
			}
		}

		button := m.ResponseButton{
			Style:     m.PrimaryButtonStyle,
			Label:     "Submit",
			ReactData: fmt.Sprintf(bountySubmitReactData, bountyID),
		}

		return ctx.Respond(m.Response{
			Type:        m.MessageResponse,
			Title:       fmt.Sprintf("Bounty: %s", title),
			Description: fmt.Sprintf("%s \n", desc),
			Buttons:     []m.ResponseButton{button},
		})
	} else if ctx.Message() != nil {
		msg := strings.Split(ctx.Message().Reaction, "-")
		reactType := msg[1]
		userID := msg[2]

		switch reactType {
		case "submit":
			err := ctx.Respond(m.Response{
				Type:        m.DMAuthorResponse,
				UserID:      userID,
				Description: fmt.Sprintf("Did <@%s> do it?", ctx.UserID()),
				Buttons: []m.ResponseButton{
					{
						Style:     m.PrimaryButtonStyle,
						Label:     "Yes",
						ReactData: fmt.Sprintf(bountyAcceptReactData, ctx.Message().ID),
					},
					{
						Style:     m.SecondaryButtonStyle,
						Label:     "No",
						ReactData: fmt.Sprintf(bountyDenyReactData, ctx.Message().ID),
					},
				},
			})
			if err != nil {
				return err
			}

			fmt.Println("test")

		case "accept":

		case "deny":
		}

		return nil
	}

	return errors.New("invalid command context")
}
