package discord

import (
	dg "github.com/bwmarrin/discordgo"
	m "github.com/hmccarty/parca/internal/models"
)

func getEmbed(resp m.Response) *dg.MessageEmbed {
	return &dg.MessageEmbed{
		Title:       resp.Title,
		Description: resp.Description,
		URL:         resp.URL,
		Color:       resp.Color,
	}
}

func getComponents(resp m.Response) []dg.MessageComponent {
	var components []dg.MessageComponent
	btnComponent, err := buttonsToComponent(resp.Buttons)
	if err != nil {
		return nil
	} else if btnComponent != nil {
		components = []dg.MessageComponent{
			btnComponent,
		}
	}
	return components
}
