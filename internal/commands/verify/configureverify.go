package verify

import (
	"regexp"

	m "github.com/hmccarty/parca/internal/models"
)

type ConfigureVerify struct {
	modIDs         []string
	createDbClient func() m.DbClient
}

func NewConfigureVerifyCommand(modIDs []string, createDbClient func() m.DbClient) m.Command {
	return &ConfigureVerify{
		modIDs:         modIDs,
		createDbClient: createDbClient,
	}
}

func (*ConfigureVerify) Name() string {
	return "configureverify"
}

func (*ConfigureVerify) Description() string {
	return "Configures verification for a server"
}

func (*ConfigureVerify) Options() []m.CommandOptionMetadata {
	return []m.CommandOptionMetadata{
		{
			Name:        "domain",
			Description: "Email domain required by users",
			Type:        m.StringOption,
			Required:    true,
		},
		{
			Name:        "verified-role",
			Description: "Role to represent valid domain holder",
			Type:        m.RoleOption,
			Required:    true,
		},
	}
}

func (cmd *ConfigureVerify) Run(ctx m.CommandContext) error {
	if len(ctx.Options()) != 2 {
		return m.ErrMissingOptions
	}

	isMod := false
	for _, modID := range cmd.modIDs {
		if ctx.UserID() == modID {
			isMod = true
			break
		}
	}
	if !isMod {
		return ctx.Respond(m.Response{
			Type:        m.MessageResponse,
			Description: "Only bot moderators can use this command",
			Color:       m.ColorRed,
		})
	}

	domain, err := ctx.Options()[0].ToString()
	if err != nil {
		return err
	}

	valid, err := regexp.MatchString(`\b[0-9A-Za-z]+\.[0-9A-Za-z]+\b`, domain)
	if err != nil || !valid {
		return ctx.Respond(m.Response{
			Type:        m.MessageResponse,
			Description: "Invalid domain, should be URL-like (e.g. `purdue.edu`)",
			Color:       m.ColorRed,
		})
	}

	roleID, err := ctx.Options()[1].ToRole()
	if err != nil {
		return err
	}

	client := cmd.createDbClient()
	err = client.AddVerifyConfig(domain, roleID, ctx.GuildID())
	if err != nil {
		return ctx.Respond(m.Response{
			Type:        m.MessageResponse,
			Description: "Failed to update verification configuration",
			Color:       m.ColorRed,
		})
	}

	return ctx.Respond(m.Response{
		Type:        m.MessageResponse,
		Description: "Configuration updated",
		Color:       m.ColorGreen,
	})
}
