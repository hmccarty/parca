package discord

import (
	"errors"
	"fmt"

	dg "github.com/bwmarrin/discordgo"
	m "github.com/hmccarty/parca/internal/models"
)

var (
	ErrNoUser              = errors.New("no user in interaction")
	ErrUnsupportedResponse = errors.New("unsupported response")
)

func createCmdCtx(
	session *dg.Session,
	interact *dg.Interaction,
	cmd m.Command) (m.CommandContext, error) {

	var userID string
	if interact.User != nil {
		userID = interact.User.ID
	} else if interact.Member != nil {
		userID = interact.Member.User.ID
	} else {
		return nil, ErrNoUser
	}

	var opts []m.CommandOption
	var message *m.ChatMessage
	switch interact.Type {
	case dg.InteractionApplicationCommand:
		message = nil

		appInteractData := interact.ApplicationCommandData()
		for _, opt := range appInteractData.Options {
			optType, err := discordOptTypeToCmdOpt(opt.Type)
			if err != nil {
				return nil, err
			}

			opts = append(opts, m.CommandOption{
				Metadata: m.CommandOptionMetadata{
					Type: optType,
					Name: opt.Name,
				},
				Value: opt.Value,
			})
		}
	case dg.InteractionMessageComponent:
		messageID := ""
		content := ""
		if interact.Message != nil {
			messageID = interact.Message.ID
			content = interact.Message.Content
		}

		msgComponentData := interact.MessageComponentData()
		message = &m.ChatMessage{
			ID:       messageID,
			Content:  content,
			Reaction: msgComponentData.CustomID,
			Values:   msgComponentData.Values,
		}
		opts = nil
	}

	return &DiscordCmdContext{
		session:   session,
		interact:  interact,
		guildID:   interact.GuildID,
		userID:    userID,
		channelID: interact.ChannelID,
		options:   opts,
		message:   message,
	}, nil
}

type DiscordCmdContext struct {
	session  *dg.Session
	interact *dg.Interaction

	guildID   string
	userID    string
	channelID string
	messageID string

	options []m.CommandOption
	message *m.ChatMessage
}

func (c *DiscordCmdContext) Respond(resp m.Response) error {
	switch resp.Type {
	case m.MessageResponse:
		// TODO: Enable more complex components
		var components []dg.MessageComponent
		btnComponent, err := buttonsToComponent(resp.Buttons)
		if err != nil {
			return err
		} else if btnComponent != nil {
			components = []dg.MessageComponent{
				btnComponent,
			}
		}

		interactResp := &dg.InteractionResponse{
			Type: dg.InteractionResponseChannelMessageWithSource,
			Data: &dg.InteractionResponseData{
				Embeds: []*dg.MessageEmbed{
					{
						Title:       resp.Title,
						Description: resp.Description,
						URL:         resp.URL,
						Color:       resp.Color,
					},
				},
				Components: components,
			},
		}

		err = c.session.InteractionRespond(c.interact, interactResp)
		if err != nil {
			return err
		}
	case m.DMAuthorResponse:
		dmChannel, err := c.session.UserChannelCreate(c.userID)
		if err != nil {
			return err
		}

		var components []dg.MessageComponent
		btnComponent, err := buttonsToComponent(resp.Buttons)
		if err != nil {
			return err
		} else if btnComponent != nil {
			components = []dg.MessageComponent{
				btnComponent,
			}
		}

		c.session.ChannelMessageSendComplex(dmChannel.ID,
			&dg.MessageSend{
				Embeds: []*dg.MessageEmbed{
					{
						Title:       resp.Title,
						Description: resp.Description,
						URL:         resp.URL,
						Color:       resp.Color,
					},
				},
				Components: components,
			},
		)

		interactResp := &dg.InteractionResponse{
			Type: dg.InteractionResponseChannelMessageWithSource,
			Data: &dg.InteractionResponseData{
				Content: "Check your DMs",
				Flags:   uint64(dg.MessageFlagsEphemeral),
			},
		}

		err = c.session.InteractionRespond(c.interact, interactResp)
		if err != nil {
			return err
		}

	case m.AddRoleResponse:
		err := c.session.GuildMemberRoleAdd(resp.GuildID, resp.UserID, resp.RoleID)
		if err != nil {
			return err
		}

		interactResp := &dg.InteractionResponse{
			Type: dg.InteractionResponseChannelMessageWithSource,
			Data: &dg.InteractionResponseData{
				Content: "Role updated",
				Flags:   uint64(dg.MessageFlagsEphemeral),
			},
		}

		err = c.session.InteractionRespond(c.interact, interactResp)
		if err != nil {
			return err
		}
	case m.MessageEditResponse:
		interactResp := &dg.InteractionResponse{
			Type: dg.InteractionResponseUpdateMessage,
			Data: &dg.InteractionResponseData{
				Embeds: []*dg.MessageEmbed{
					{
						Title:       resp.Title,
						Description: resp.Description,
						URL:         resp.URL,
						Color:       resp.Color,
					},
				},
			},
		}

		err := c.session.InteractionRespond(c.interact, interactResp)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("%v: type: %d", ErrUnsupportedResponse, resp.Type)
	}
	return nil
}

func (c *DiscordCmdContext) GuildID() string {
	return c.guildID
}

func (c *DiscordCmdContext) UserID() string {
	return c.userID
}

func (c *DiscordCmdContext) ChannelID() string {
	return c.channelID
}

func (c *DiscordCmdContext) MessageID() string {
	return c.messageID
}

func (c *DiscordCmdContext) Options() []m.CommandOption {
	return c.options
}

func (c *DiscordCmdContext) Message() *m.ChatMessage {
	return c.message
}

func (c *DiscordCmdContext) GetGuildNameFromID(id string) (string, error) {
	guild, err := c.session.State.Guild(id)
	if err != nil || guild == nil {
		guild, err = c.session.Guild(id)
		if err != nil || guild == nil {
			return "", err
		}
	}
	return guild.Name, nil
}

func (c *DiscordCmdContext) GetUserNameFromIDs(userID, _ string) (string, error) {
	user, err := c.session.User(userID)
	if err != nil || user == nil {
		return "", err
	}
	return user.Username, nil
}

func (c *DiscordCmdContext) GetChannelNameFromIDs(channelID, _ string) (string, error) {
	channel, err := c.session.State.Channel(channelID)
	if err != nil || channel == nil {
		channel, err = c.session.Channel(channelID)
		if err != nil || channel == nil {
			return "", err
		}
	}
	return channel.Name, nil
}

func (c *DiscordCmdContext) GetRoleNameFromIDs(roleID, guildID string) (string, error) {
	role, err := c.session.State.Role(guildID, roleID)
	if err != nil || role == nil {
		roles, err := c.session.GuildRoles(guildID)
		for _, r := range roles {
			if r.ID == roleID {
				role = r
				break
			}
		}

		if err != nil || role == nil {
			return "", err
		}
	}
	return role.Name, nil
}
