package discord

import (
	"fmt"
	"log"
	"strings"

	dg "github.com/bwmarrin/discordgo"
	m "github.com/hmccarty/parca/internal/models"
)

type InteractHandler func(s *dg.Session, i *dg.InteractionCreate)

func getInteractHandler(cmds []m.Command) InteractHandler {
	appInteractHandlers := make(map[string]InteractHandler, len(cmds))
	msgInteractHandlers := make(map[string]InteractHandler, len(cmds))
	for _, cmd := range cmds {
		appInteractHandlers[cmd.Name()] = getAppInteractHandler(cmd)
		msgInteractHandlers[cmd.Name()] = getMsgInteractHandler(cmd)
	}

	return func(s *dg.Session, i *dg.InteractionCreate) {
		switch i.Type {
		case dg.InteractionApplicationCommand:
			if h, ok := appInteractHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}
			// TODO: Handle unexpected command case
		case dg.InteractionMessageComponent:
			key := strings.Split(i.MessageComponentData().CustomID, "-")[0]
			if h, ok := msgInteractHandlers[key]; ok {
				h(s, i)
			}
		}
	}
}

func getAppInteractHandler(cmd m.Command) InteractHandler {
	return func(s *dg.Session, i *dg.InteractionCreate) {
		interactData, err := dataFromInteraction(i.Interaction)
		if err != nil {

		}

		appData := i.ApplicationCommandData()
		options := make([]m.CommandOption, len(appData.Options))
		for i, v := range appData.Options {
			option, err := optionFromInteraction(s, data.GuildID, v)
			if err != nil {
				log.Println(err)
			}
			options[i] = option
		}

		resp := cmd.Run(data, options)
		comp, err := componentsFromResponse(resp)
		if err != nil {
			fmt.Println(err)
		}

		switch resp.Type {
		case m.MessageResponse:
			err := s.InteractionRespond(i.Interaction, &dg.InteractionResponse{
				Type: dg.InteractionResponseChannelMessageWithSource,
				Data: &dg.InteractionResponseData{
					Embeds: []*dg.MessageEmbed{
						{
							Title:       resp.Title,
							Description: resp.Description,
							URL:         resp.URL,
						},
					},

					Components: comp,
				},
			})
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func getMsgInteractHandler() {

}

func getMsgCreateHandler() {

}
