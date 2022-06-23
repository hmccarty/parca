package verify

import (
	"fmt"
	"math/rand"
	"regexp"
	"time"

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
	opt := ctx.Options()[0]
	if opt.Metadata.Name == "email" {
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

		seed := rand.NewSource(time.Now().UnixNano())
		code := fmt.Sprintf("%d", rand.New(seed).Intn(6000-1000)+1000)

		var subject string
		serverName, err := ctx.GetGuildNameFromID(ctx.GuildID())
		if err != nil {
			subject = fmt.Sprintf("%s Server Verification Code", serverName)
		} else {
			subject = "Server Verification Code"
		}

		err = cmd.emailClient.SendEmail(email, subject, code)
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

		return ctx.Respond(m.Response{
			Type:        m.AckResponse,
			IsForm:      true,
			Title:       "Server Verification",
			Description: "Check the entered email for the verification code.",
			CustomID:    fmt.Sprintf("verify-%s-%s", ctx.GuildID(), ctx.UserID()),
			Inputs: []m.ResponseInput{
				{
					Style:    m.ShortInputStyle,
					Label:    "Verification Code",
					Required: true,
					CustomID: "verify-code",
				},
			},
		})
	} else if opt.Metadata.Name == "verify-code" {
		input, err := opt.ToString()
		if err != nil {
			return err
		}

		client := cmd.createDbClient()
		code, guildID, err := client.GetVerifyCode(ctx.UserID())
		if err != nil {
			return err
		}

		if code == "" {
			return nil
		} else if input != code {
			return ctx.Respond(m.Response{
				Type:        m.AckResponse,
				IsEphemeral: true,
				Description: "Invalid code",
			})
		}

		_, roleID, err := client.GetVerifyConfig(guildID)
		if err != nil {
			return err
		}

		err = ctx.Respond(m.Response{
			Type:    m.AddRoleResponse,
			GuildID: guildID,
			UserID:  ctx.UserID(),
			RoleID:  roleID,
		})
		if err != nil {
			return err
		}

		guildName, err := ctx.GetGuildNameFromID(guildID)
		if err != nil {
			return err
		}

		return ctx.Respond(m.Response{
			Type:        m.AckResponse,
			IsEphemeral: true,
			Description: fmt.Sprintf("You have been verified on %s", guildName),
		})
	} else {
		fmt.Println(opt.Metadata.Name)
		return m.ErrMissingOptions
	}
}
