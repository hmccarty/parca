package discord

import (
	dg "github.com/bwmarrin/discordgo"
	m "github.com/hmccarty/arc-assistant/internal/models"
)

func appFromCommand(command m.Command) *dg.ApplicationCommand {
	var appOptions []*dg.ApplicationCommandOption = nil
	if len(command.Options()) > 0 {
		appOptions = make([]*dg.ApplicationCommandOption, len(command.Options()))
		for i, v := range command.Options() {
			appOptions[i] = &dg.ApplicationCommandOption{
				Type:        dg.ApplicationCommandOptionType(v.Type),
				Name:        v.Name,
				Required:    v.Required,
				Description: "Description",
			}
		}
	}

	return &dg.ApplicationCommand{
		Name:        command.Name(),
		Description: command.Description(),
		Options:     appOptions,
	}
}

func optionFromInteraction(interactionOption *dg.ApplicationCommandInteractionDataOption) (m.CommandOption, error) {
	option := m.CommandOption{
		Name:  interactionOption.Name,
		Type:  m.CommandOptionType(interactionOption.Type),
		Value: interactionOption.Value,
	}
	return option, nil
}

func userFromData(user *dg.User) *m.User {
	if user == nil {
		return nil
	}
	return &m.User{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
	}
}

func memberFromData(member *dg.Member) *m.Member {
	if member == nil {
		return nil
	}
	return &m.Member{
		GuildID: member.GuildID,
		User:    userFromData(member.User),
		Roles:   member.Roles,
	}
}

func roleFromData(role *dg.Role) *m.Role {
	if role == nil {
		return nil
	}
	return &m.Role{
		ID:   role.ID,
		Name: role.Name,
	}
}

func dataFromInteraction(interaction *dg.Interaction) (m.CommandData, error) {
	data := m.CommandData{
		GuildID:   interaction.GuildID,
		ChannelID: interaction.ChannelID,
		User:      userFromData(interaction.User),
		Member:    memberFromData(interaction.Member),
	}
	return data, nil
}
