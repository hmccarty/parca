package discord

import (
	dg "github.com/bwmarrin/discordgo"
	m "github.com/hmccarty/parca/internal/models"
)

func getMessage(resp m.Response) *dg.MessageSend {
	if resp.Title == "" {
		return &dg.MessageSend{
			Content:    resp.Description,
			Components: getComponents(resp),
		}
	} else {
		return &dg.MessageSend{
			Embeds: []*dg.MessageEmbed{
				getEmbed(resp),
			},
			Components: getComponents(resp),
		}
	}
}

func getMessageEdit(resp m.Response) *dg.MessageEdit {
	if resp.Title == "" {
		return &dg.MessageEdit{
			ID:         resp.MessageID,
			Channel:    resp.ChannelID,
			Content:    &resp.Description,
			Components: getComponents(resp),
		}
	} else {
		return &dg.MessageEdit{
			ID:      resp.MessageID,
			Channel: resp.ChannelID,
			Embeds: []*dg.MessageEmbed{
				getEmbed(resp),
			},
			Components: getComponents(resp),
		}
	}
}

func getInteraction(resp m.Response) *dg.InteractionResponse {
	if resp.IsForm {
		return &dg.InteractionResponse{
			Type: dg.InteractionResponseModal,
			Data: &dg.InteractionResponseData{
				Flags:      uint64(getFlags(resp)),
				Components: getComponents(resp),
				CustomID:   resp.CustomID,
				Title:      resp.Title,
			},
		}
	} else if resp.Title == "" {
		return &dg.InteractionResponse{
			Type: dg.InteractionResponseChannelMessageWithSource,
			Data: &dg.InteractionResponseData{
				Flags:      uint64(getFlags(resp)),
				Content:    resp.Description,
				Components: getComponents(resp),
			},
		}
	} else {
		return &dg.InteractionResponse{
			Type: dg.InteractionResponseChannelMessageWithSource,
			Data: &dg.InteractionResponseData{
				Flags: uint64(getFlags(resp)),
				Embeds: []*dg.MessageEmbed{
					getEmbed(resp),
				},
				Components: getComponents(resp),
			},
		}
	}
}

func getEmbed(resp m.Response) *dg.MessageEmbed {
	return &dg.MessageEmbed{
		Title:       resp.Title,
		Description: resp.Description,
		URL:         resp.URL,
		Color:       resp.Color,
	}
}

func getComponents(resp m.Response) []dg.MessageComponent {
	components := []dg.MessageComponent{}

	for _, v := range resp.Inputs {
		inputComponent, err := inputToComponent(v)
		if err != nil {
			continue
		}
		components = append(components, inputComponent)
	}

	btnComponent, err := buttonsToComponent(resp.Buttons)
	if err != nil {
		return []dg.MessageComponent{}
	} else if btnComponent != nil {
		components = append(components, btnComponent)
	}

	return components
}

func getFlags(resp m.Response) dg.MessageFlags {
	var flags dg.MessageFlags

	if resp.IsEphemeral {
		flags |= dg.MessageFlagsEphemeral
	}

	// TODO: Increase flag options

	return flags
}
