package general

import (
	"fmt"
	"strings"

	m "github.com/hmccarty/parca/internal/models"
)

type Embed struct{}

func NewEmbedCommand() m.Command {
	return &Embed{}
}

func (*Embed) Name() string {
	return "embed"
}

func (*Embed) Description() string {
	return "Echos message in form of an embed"
}

func (*Embed) Options() []m.CommandOptionMetadata {
	return []m.CommandOptionMetadata{
		{
			Name:        "channel",
			Description: "Channel to send message",
			Type:        m.ChannelOption,
			Required:    true,
		},
	}
}

func (*Embed) Run(ctx m.ChatContext) error {
	if len(ctx.Options()) == 1 {
		channelID, err := ctx.Options()[0].ToChannel()
		if err != nil {
			return err
		}

		return ctx.Respond(m.Response{
			Type:     m.AckResponse,
			IsForm:   true,
			Title:    "Embed Message",
			CustomID: fmt.Sprintf("embed-%s", channelID),
			Inputs: []m.ResponseInput{
				{
					Style:    m.ShortInputStyle,
					Label:    "Message Title",
					Required: false,
					CustomID: "embed-title",
				},
				{
					Style:    m.ParagraphInputStyle,
					Label:    "Message Content",
					Required: true,
					CustomID: "embed-content",
				},
			},
		})
	} else {
		title, err := ctx.Options()[0].ToString()
		if err != nil {
			return err
		}

		desc, err := ctx.Options()[1].ToString()
		if err != nil {
			return err
		}

		channelID := strings.Split(ctx.UniqueID(), "-")[1]

		err = ctx.Respond(m.Response{
			Type:        m.MessageResponse,
			Title:       title,
			Description: desc,
			ChannelID:   channelID,
			GuildID:     ctx.GuildID(),
		})
		if err != nil {
			return err
		}

		return ctx.Respond(m.Response{
			Type:        m.AckResponse,
			Description: "Message sent",
			IsEphemeral: true,
		})
	}
}
