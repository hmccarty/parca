package verify

import (
	"fmt"
	"math/rand"
	"regexp"

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

func (*Verify) Options() []m.CommandOptionMetadata {
	return []m.CommandOptionMetadata{
		{
			Name:        "email",
			Description: "Email with the domain required by the server",
			Type:        m.StringOption,
			Required:    true,
		},
	}
}

func (cmd *Verify) Run(ctx m.ChatContext) error {
	if len(ctx.Options()) != 1 {
		return m.ErrMissingOptions
	}

	email, err := ctx.Options()[0].ToString()
	if err != nil {
		return err
	}

	client := cmd.createDbClient()
	domain, _, err := client.GetVerifyConfig(ctx.GuildID())
	if err != nil {
		return ctx.Respond(m.Response{
			Type:        m.AckResponse,
			IsEphemeral: true,
			Description: "Verification not configured for server, contact server admin",
		})
	}

	validEmailPattern := fmt.Sprintf(`\b[0-9A-Za-z]+@%s\b`, domain)
	isValidEmail, err := regexp.MatchString(validEmailPattern, email)
	if err != nil {
		return ctx.Respond(m.Response{
			Type:        m.AckResponse,
			IsEphemeral: true,
			Description: "Failed to validate email",
		})
	} else if !isValidEmail {
		invalidMsg := fmt.Sprintf("Invalid email, ensure you use an email with a `%s` domain",
			domain)
		return ctx.Respond(m.Response{
			Type:        m.AckResponse,
			IsEphemeral: true,
			Description: invalidMsg,
		})
	}

	code := fmt.Sprintf("%d", rand.Intn(6000-1000)+1000)
	err = cmd.emailClient.SendEmail(email, "Discord Server Verification", code)
	if err != nil {
		return ctx.Respond(m.Response{
			Type:        m.AckResponse,
			IsEphemeral: true,
			Description: fmt.Sprintf("Failed to send email: %s", err),
		})
	}

	err = client.AddVerifyCode(code, ctx.UserID(), ctx.GuildID())
	if err != nil {
		return ctx.Respond(m.Response{
			Type:        m.AckResponse,
			IsEphemeral: true,
			Description: "Failed to save code, try again later",
		})
	}

	err = ctx.Respond(m.Response{
		Type:        m.AckResponse,
		Description: "Check your DMs",
		IsEphemeral: true,
	})
	if err != nil {
		return err
	}

	return ctx.Respond(m.Response{
		Type:        m.DMResponse,
		Description: "Respond with the code sent to your email here",
	})
}
