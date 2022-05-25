package discord

import (
	"errors"
	"fmt"

	dg "github.com/bwmarrin/discordgo"
	m "github.com/hmccarty/parca/internal/models"
)

var (
	ErrOptionTypeNotSupported = errors.New("option type is not supported")
)

func cmdToApp(cmd m.Command) (*dg.ApplicationCommand, error) {
	var opts []*dg.ApplicationCommandOption
	for _, cmdOpt := range cmd.Options() {
		opt, err := cmdOptMetadataToAppOpt(cmdOpt)
		if err != nil {
			return nil, err
		}
		opts = append(opts, opt)
	}

	return &dg.ApplicationCommand{
		Name:        cmd.Name(),
		Description: cmd.Description(),
		Options:     opts,
	}, nil
}

func cmdOptMetadataToAppOpt(cmdOptMetadata m.CommandOptionMetadata) (*dg.ApplicationCommandOption, error) {
	optType, err := cmdOptTypeToDiscordOpt(cmdOptMetadata.Type)
	if err != nil {
		return nil, err
	}

	return &dg.ApplicationCommandOption{
		Type:        optType,
		Name:        cmdOptMetadata.Name,
		Description: cmdOptMetadata.Description,
		Required:    cmdOptMetadata.Required,
	}, nil
}

func cmdOptTypeToDiscordOpt(optType m.CommandOptionType) (dg.ApplicationCommandOptionType, error) {
	switch optType {
	case m.StringOption:
		return dg.ApplicationCommandOptionString, nil
	case m.IntegerOption:
		return dg.ApplicationCommandOptionInteger, nil
	case m.FloatOption:
		return dg.ApplicationCommandOptionNumber, nil
	case m.BooleanOption:
		return dg.ApplicationCommandOptionBoolean, nil
	case m.UserOption:
		return dg.ApplicationCommandOptionUser, nil
	case m.ChannelOption:
		return dg.ApplicationCommandOptionChannel, nil
	case m.RoleOption:
		return dg.ApplicationCommandOptionRole, nil
	default:
		return 0, ErrOptionTypeNotSupported
	}
}

func discordOptTypeToCmdOpt(optType dg.ApplicationCommandOptionType) (m.CommandOptionType, error) {
	switch optType {
	case dg.ApplicationCommandOptionString:
		return m.StringOption, nil
	case dg.ApplicationCommandOptionInteger:
		return m.IntegerOption, nil
	case dg.ApplicationCommandOptionNumber:
		return m.FloatOption, nil
	case dg.ApplicationCommandOptionBoolean:
		return m.BooleanOption, nil
	case dg.ApplicationCommandOptionUser:
		return m.UserOption, nil
	case dg.ApplicationCommandOptionChannel:
		return m.ChannelOption, nil
	case dg.ApplicationCommandOptionRole:
		return m.RoleOption, nil
	default:
		return 0, ErrOptionTypeNotSupported
	}
}

func emojiToComponentEmoji(emoji m.Emoji) dg.ComponentEmoji {
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

func buttonsToComponent(buttons []m.ResponseButton) (dg.MessageComponent, error) {
	if buttons == nil {
		fmt.Println("true")
		return nil, nil
	}

	actionRow := &dg.ActionsRow{
		Components: []dg.MessageComponent{},
	}

	for _, button := range buttons {
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
				Emoji:    emojiToComponentEmoji(button.Emoji),
				Style:    style,
				URL:      button.URL,
				CustomID: button.ReactData,
			})
	}

	return actionRow, nil
}
