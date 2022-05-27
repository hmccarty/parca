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
	ErrNoUserID            = errors.New("no user ID in response")
	ErrNoGuildID           = errors.New("no guild id in response")
	ErrNoChannelID         = errors.New("no channel id in response")
	ErrNoMessageID         = errors.New("no message id in response")
)

func createCtx(
	session *dg.Session,
	interact *dg.Interaction,
	cmd m.Command) (m.ChatContext, error) {

	var userID string
	if interact.User != nil {
		userID = interact.User.ID
	} else if interact.Member != nil {
		userID = interact.Member.User.ID
	} else {
		return nil, ErrNoUser
	}

	var contextType m.ChatContextType
	var opts []m.CommandOption
	var err error
	var uniqueID string
	var message *m.ChatMessage
	switch interact.Type {
	case dg.InteractionApplicationCommand:
		contextType = m.CommandCall
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
	case dg.InteractionModalSubmit:
		contextType = m.CommandReply
		message = nil

		modalData := interact.ModalSubmitData()
		uniqueID = modalData.CustomID
		opts, err = componentsToCmdOpts(modalData.Components)
		if err != nil {
			return nil, err
		}
	case dg.InteractionMessageComponent:
		contextType = m.CommandReply
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

	return &DiscordContext{
		contextType: contextType,
		session:     session,
		interact:    interact,
		guildID:     interact.GuildID,
		userID:      userID,
		channelID:   interact.ChannelID,
		uniqueID:    uniqueID,
		options:     opts,
		message:     message,
	}, nil
}

type DiscordContext struct {
	contextType m.ChatContextType
	session     *dg.Session
	interact    *dg.Interaction

	guildID   string
	userID    string
	channelID string
	messageID string
	uniqueID  string

	options []m.CommandOption
	message *m.ChatMessage
}

func (c *DiscordContext) Type() m.ChatContextType {
	return c.contextType
}

func (c *DiscordContext) Respond(resp m.Response) error {
	switch resp.Type {
	case m.MessageResponse:
		if resp.ChannelID == "" {
			return ErrNoChannelID
		}

		_, err := c.session.ChannelMessageSendComplex(resp.ChannelID, getMessage(resp))
		if err != nil {
			return err
		}

	case m.AckResponse:
		err := c.session.InteractionRespond(c.interact, getInteraction(resp))
		if err != nil {
			return err
		}

	case m.DMResponse:
		var userID string
		if resp.UserID != "" {
			userID = resp.UserID
		} else {
			userID = c.userID
		}

		dmChannel, err := c.session.UserChannelCreate(userID)
		if err != nil {
			return err
		}

		_, err = c.session.ChannelMessageSendComplex(dmChannel.ID, getMessage(resp))
		if err != nil {
			return err
		}

	case m.AddRoleResponse:
		err := c.session.GuildMemberRoleAdd(resp.GuildID, resp.UserID, resp.RoleID)
		if err != nil {
			return err
		}

	case m.MessageEditResponse:
		if resp.MessageID == "" {
			return ErrNoMessageID
		} else if resp.ChannelID == "" {
			return ErrNoChannelID
		}

		_, err := c.session.ChannelMessageEditComplex(
			&dg.MessageEdit{
				Components: getComponents(resp),
				ID:         resp.MessageID,
				Channel:    resp.ChannelID,
				Embeds: []*dg.MessageEmbed{
					getEmbed(resp),
				},
			},
		)
		if err != nil {
			return err
		}

	default:
		return fmt.Errorf("%v: type: %d", ErrUnsupportedResponse, resp.Type)
	}
	return nil
}

func (c *DiscordContext) GuildID() string {
	return c.guildID
}

func (c *DiscordContext) UserID() string {
	return c.userID
}

func (c *DiscordContext) ChannelID() string {
	return c.channelID
}

func (c *DiscordContext) MessageID() string {
	return c.messageID
}

func (c *DiscordContext) UniqueID() string {
	return c.uniqueID
}

func (c *DiscordContext) Options() []m.CommandOption {
	return c.options
}

func (c *DiscordContext) Message() *m.ChatMessage {
	return c.message
}

func (c *DiscordContext) GetGuildNameFromID(id string) (string, error) {
	guild, err := c.session.State.Guild(id)
	if err != nil || guild == nil {
		guild, err = c.session.Guild(id)
		if err != nil || guild == nil {
			return "", err
		}
	}
	return guild.Name, nil
}

func (c *DiscordContext) GetUserNameFromIDs(userID, _ string) (string, error) {
	user, err := c.session.User(userID)
	if err != nil || user == nil {
		return "", err
	}
	return user.Username, nil
}

func (c *DiscordContext) GetChannelNameFromIDs(channelID, _ string) (string, error) {
	channel, err := c.session.State.Channel(channelID)
	if err != nil || channel == nil {
		channel, err = c.session.Channel(channelID)
		if err != nil || channel == nil {
			return "", err
		}
	}
	return channel.Name, nil
}

func (c *DiscordContext) GetRoleNameFromIDs(roleID, guildID string) (string, error) {
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
