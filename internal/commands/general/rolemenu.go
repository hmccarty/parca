package general

import (
	"errors"
	"fmt"
	"strings"

	m "github.com/hmccarty/parca/internal/models"
)

const (
	roleReactData = "rolemenu-%s"
)

type RoleMenu struct{}

func NewRoleMenuCommand() m.Command {
	return &RoleMenu{}
}

func (*RoleMenu) Name() string {
	return "rolemenu"
}

func (*RoleMenu) Description() string {
	return "Creates a role menu (max 5 roles per menu)"
}

func (*RoleMenu) Options() []m.CommandOptionMetadata {
	return []m.CommandOptionMetadata{
		{
			Name:        "title",
			Description: "Title of role menu",
			Type:        m.StringOption,
			Required:    true,
		},
		{
			Name:        "role1",
			Description: "Role to displayed on menu",
			Type:        m.RoleOption,
			Required:    true,
		},
		{
			Name:        "role2",
			Description: "Role to displayed on menu",
			Type:        m.RoleOption,
			Required:    false,
		},
		{
			Name:        "role3",
			Description: "Role to displayed on menu",
			Type:        m.RoleOption,
			Required:    false,
		},
		{
			Name:        "role4",
			Description: "Role to displayed on menu",
			Type:        m.RoleOption,
			Required:    false,
		},
		{
			Name:        "role5",
			Description: "Role to displayed on menu",
			Type:        m.RoleOption,
			Required:    false,
		},
	}
}

func (*RoleMenu) Run(ctx m.CommandContext) error {
	if ctx.Options() != nil {
		// If command is called
		opts := ctx.Options()
		title, err := opts[0].ToString()
		if err != nil {
			return err
		}

		buttons := make([]m.ResponseButton, len(opts[1:]))
		for i, opt := range opts[1:] {
			roleID, err := opt.ToRole()
			if err != nil {
				return err
			}

			roleName, err := ctx.GetRoleNameFromIDs(roleID, ctx.GuildID())

			buttons[i] = m.ResponseButton{
				Style:     m.PrimaryButtonStyle,
				Label:     roleName,
				ReactData: fmt.Sprintf(roleReactData, roleID),
			}
		}

		return ctx.Respond(m.Response{
			Type:        m.MessageResponse,
			Title:       fmt.Sprintf("%s Role Menu", title),
			Description: "Click the buttons below to add role",
			Buttons:     buttons,
		})
	} else if ctx.Message() != nil {
		// If a command response has a reaction
		msg := ctx.Message()
		roleID := strings.Split(msg.Reaction, "-")[1]

		return ctx.Respond(m.Response{
			Type:    m.AddRoleResponse,
			GuildID: ctx.GuildID(),
			UserID:  ctx.UserID(),
			RoleID:  roleID,
		})
	}

	return errors.New("invalid command context")
}
