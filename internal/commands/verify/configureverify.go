package verify

import (
	m "github.com/hmccarty/parca/internal/models"
)

type ConfigureVerify struct {
	createDbClient func() m.DbClient
}

func NewConfigureVerifyCommand(createDbClient func() m.DbClient) m.Command {
	return &ConfigureVerify{
		createDbClient: createDbClient,
	}
}

func (*ConfigureVerify) Name() string {
	return "configureverify"
}

func (*ConfigureVerify) Description() string {
	return "Configures verification for a server"
}

func (*ConfigureVerify) Options() []m.CommandOption {
	return []m.CommandOption{
		{
			Name:     "domain",
			Type:     m.StringOption,
			Required: true,
		},
		{
			Name:     "verified-role",
			Type:     m.RoleOption,
			Required: true,
		},
	}
}

func (command *ConfigureVerify) Run(data m.CommandData, opts []m.CommandOption) m.Response {
	if len(opts) != 2 {
		return m.Response{
			Description: "Missing arguments",
		}
	}

	// TODO: verify permissions

	var domain string
	var roleID string
	for _, option := range opts {
		switch option.Name {
		case "domain":
			domain = option.Value.(string)
		case "verified-role":
			roleID = option.Value.(string)
		}
	}

	// TODO: verify domain

	client := command.createDbClient()
	err := client.AddVerifyConfig(domain, roleID, data.GuildID)
	if err != nil {
		return m.Response{
			Description: "Failed to update verification configuration",
		}
	}

	return m.Response{
		Description: "Configuration updated",
	}
}

func (*ConfigureVerify) HandleReaction(data m.CommandData, reaction string) m.Response {
	return m.Response{
		Description: "Not expecting a reaction",
	}
}
