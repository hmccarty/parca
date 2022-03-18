package discord

import (
	"log"

	dg "github.com/bwmarrin/discordgo"
	m "github.com/hmccarty/parca/internal/models"
)

type CommandHandler func(*dg.Session, *dg.InteractionCreate)

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

func handlerFromCommand(command m.Command) CommandHandler {
	return func(s *dg.Session, i *dg.InteractionCreate) {
		data, _ := dataFromInteraction(i.Interaction)
		appData := i.ApplicationCommandData()
		options := make([]m.CommandOption, len(appData.Options))
		for i, v := range appData.Options {
			option, err := optionFromInteraction(v)
			if err != nil {
				log.Println(err)
			}
			options[i] = option
		}

		resp := command.Run(data, options)
		s.InteractionRespond(i.Interaction, &dg.InteractionResponse{
			Type: dg.InteractionResponseChannelMessageWithSource,
			Data: &dg.InteractionResponseData{
				Embeds: []*dg.MessageEmbed{
					{
						Title:       resp.Title,
						Description: resp.Description,
						URL:         resp.URL,
					},
				},
			},
		})
	}
}
