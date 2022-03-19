package verify

import (
	"fmt"
	"math/rand"

	m "github.com/hmccarty/parca/internal/models"
)

type Verify struct {
	createDbClient func() m.DbClient
	emailClient    m.EmailClient
}

func NewVerifyCommand(createDbClient func() m.DbClient, emailClient m.EmailClient) m.Command {
	return &Verify{
		createDbClient: createDbClient,
		emailClient:    emailClient,
	}
}

func (*Verify) Name() string {
	return "verify"
}

func (*Verify) Description() string {
	return "Prompts server for verification role"
}

func (*Verify) Options() []m.CommandOption {
	return []m.CommandOption{
		{
			Name:     "email",
			Type:     m.StringOption,
			Required: true,
		},
	}
}

func (command *Verify) Run(data m.CommandData, opts []m.CommandOption) m.Response {
	if len(opts) != 1 {
		return m.Response{
			Description: "Missing arguments",
		}
	}

	// TODO: check email matches domain
	email := opts[0].Value.(string)
	code := fmt.Sprintf("%d", rand.Intn(6000-1000)+1000)
	go command.emailClient.SendEmail(email, "Discord Server Verification", code)

	var userID string
	if data.User != nil {
		userID = data.User.ID
	} else {
		userID = data.Member.User.ID
	}

	client := command.createDbClient()
	err := client.AddVerifyCode(code, userID, data.GuildID)
	if err != nil {
		return m.Response{
			Description: "Failed to save code, try again later",
		}
	}

	return m.Response{
		Description: "Check your email for code and respond in DMs",
	}
}

func (*Verify) HandleReaction(data m.CommandData, reaction string) m.Response {
	return m.Response{
		Description: "Not expecting a reaction",
	}
}
