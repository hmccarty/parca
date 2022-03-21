package discord

import (
	"fmt"

	dg "github.com/bwmarrin/discordgo"
	m "github.com/hmccarty/parca/internal/models"
)

type DiscordEventContext struct {
	session *dg.Session

	guildID   string
	userID    string
	channelID string
	messageID string

	message *m.ChatMessage
}

func (c *DiscordEventContext) Respond(resp m.Response) error {
	switch resp.Type {
	case m.MessageResponse:
		c.session.ChannelMessageSendEmbed(c.channelID,
			&dg.MessageEmbed{
				Title:       resp.Title,
				Description: resp.Description,
				URL:         resp.URL,
				Color:       resp.Color,
			})
	case m.AddRoleResponse:
		err := c.session.GuildMemberRoleAdd(resp.GuildID,
			resp.UserID, resp.RoleID)
		if err != nil {
			fmt.Println(err)
		}

		c.session.ChannelMessageSendEmbed(c.channelID,
			&dg.MessageEmbed{
				Title:       resp.Title,
				Description: resp.Description,
				URL:         resp.URL,
				Color:       resp.Color,
			})
	case m.RemoveRoleResponse:
		err := c.session.GuildMemberRoleRemove(resp.GuildID,
			resp.UserID, resp.RoleID)
		if err != nil {
			fmt.Println(err)
		}

		c.session.ChannelMessageSendEmbed(c.channelID,
			&dg.MessageEmbed{
				Title:       resp.Title,
				Description: resp.Description,
				URL:         resp.URL,
				Color:       resp.Color,
			})
	}
	return nil
}

func (c *DiscordEventContext) GuildID() string {
	return c.guildID
}

func (c *DiscordEventContext) UserID() string {
	return c.userID
}

func (c *DiscordEventContext) ChannelID() string {
	return c.channelID
}

func (c *DiscordEventContext) MessageID() string {
	return c.messageID
}

func (c *DiscordEventContext) Message() *m.ChatMessage {
	return c.message
}

func (c *DiscordEventContext) GetGuildNameFromID(id string) (string, error) {
	guild, err := c.session.State.Guild(id)
	if err != nil || guild == nil {
		guild, err = c.session.Guild(id)
		if err != nil || guild == nil {
			return "", err
		}
	}
	return guild.Name, nil
}

func (c *DiscordEventContext) GetUserNameFromIDs(userID, _ string) (string, error) {
	user, err := c.session.User(userID)
	if err != nil || user == nil {
		return "", err
	}
	return user.Username, nil
}

func (c *DiscordEventContext) GetChannelNameFromIDs(channelID, _ string) (string, error) {
	channel, err := c.session.State.Channel(channelID)
	if err != nil || channel == nil {
		channel, err = c.session.Channel(channelID)
		if err != nil || channel == nil {
			return "", err
		}
	}
	return channel.Name, nil
}

func (c *DiscordEventContext) GetRoleNameFromIDs(roleID, guildID string) (string, error) {
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

func (c *DiscordEventContext) IsChannelDM(channelID, _ string) (bool, error) {
	channel, err := c.session.State.Channel(channelID)
	if err != nil || channel == nil {
		channel, err = c.session.Channel(channelID)
		if err != nil || channel == nil {
			return false, err
		}
	}
	return channel.Type == dg.ChannelTypeDM, nil
}
