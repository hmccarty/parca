package discord

import (
	"errors"

	dg "github.com/bwmarrin/discordgo"
	m "github.com/hmccarty/parca/internal/models"
)

func convertToOpt(s *dg.Session, guildID string, interactOpt *dg.ApplicationCommandInteractionDataOption) (m.CommandOption, error) {
	optType := m.CommandOptionType(interactOpt.Type)

	var optValue interface{}
	switch optType {
	case m.UserOption:
		user := interactOpt.UserValue(s)
		if user == nil {
			return m.CommandOption{}, errors.New("User not found")
		}
		optValue = m.User{
			ID:       user.ID,
			Email:    user.Email,
			Username: user.Username,
		}
	case m.RoleOption:
		role := interactOpt.RoleValue(s, guildID)
		if role == nil {
			return m.CommandOption{}, errors.New("Role not found")
		}
		optValue = m.Role{
			ID:   role.ID,
			Name: role.Name,
		}
	default:
		optValue = interactOpt.Value
	}

	return m.CommandOption{
		Name:  interactOpt.Name,
		Type:  optType,
		Value: optValue,
	}, nil
}

func convertToComponentEmoji(emoji m.Emoji) dg.ComponentEmoji {
	switch emoji {
	case m.ThumbsUpEmoji:
		return dg.ComponentEmoji{
			Name: "👍",
		}
	case m.ThumbsDownEmoji:
		return dg.ComponentEmoji{
			Name: "👎",
		}
	}
	return dg.ComponentEmoji{}
}

func convertToComponent(resp m.Response) ([]dg.MessageComponent, error) {
	if resp.Buttons == nil {
		return nil, nil
	}

	actionRow := &dg.ActionsRow{
		Components: []dg.MessageComponent{},
	}

	for _, button := range resp.Buttons {
		var style dg.ButtonStyle
		switch button.Style {
		case m.PrimaryButtonStyle:
			style = dg.PrimaryButton
		case m.SecondaryButtonStyle:
			style = dg.DangerButton
		case m.LinkButtonStyle:
			style = dg.LinkButton
		default:
			style = dg.SecondaryButton
		}

		actionRow.Components = append(actionRow.Components,
			dg.Button{
				Label:    button.Label,
				Emoji:    componentEmojiFromEmoji(button.Emoji),
				Style:    style,
				URL:      button.URL,
				CustomID: button.ReactData,
			})
	}

	return []dg.MessageComponent{actionRow}, nil
}

func messageFromData(message *dg.Message) *m.Message {
	if message == nil {
		return nil
	}
	return &m.Message{
		ID:        message.ID,
		ChannelID: message.ChannelID,
		GuildID:   message.GuildID,
		Content:   message.Content,
		Timestamp: message.Timestamp,
		Author:    userFromData(message.Author),
		Member:    memberFromData(message.Member),
	}
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

func cmdDataFromInteract(interact *dg.Interaction) m.CommandData {
	return m.CommandData{
		GuildID:   interact.GuildID,
		ChannelID: interact.ChannelID,
		User:      userFromData(interact.User),
		Member:    memberFromData(interact.Member),
	}
}
