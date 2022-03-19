package currency

import (
	"fmt"
	"strings"

	m "github.com/hmccarty/parca/internal/models"
)

const (
	roleReactData = "rolemenu-%s"
)

type RoleMenu struct {
	createDbClient func() m.DbClient
}

func NewRoleMenuCommand(createDbClient func() m.DbClient) m.Command {
	return &RoleMenu{
		createDbClient: createDbClient,
	}
}

func (*RoleMenu) Name() string {
	return "rolemenu"
}

func (*RoleMenu) Description() string {
	return "Creates a role menu (max 5 roles per menu)"
}

func (*RoleMenu) Options() []m.CommandOption {
	return []m.CommandOption{
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

func (command *RoleMenu) Run(data m.CommandData, opts []m.CommandOption) m.Response {
	title := opts[0].Value.(string)

	buttons := make([]m.ResponseButton, len(opts[1:]))
	for i, v := range opts[1:] {
		role := v.Value.(m.Role)
		buttons[i] = m.ResponseButton{
			Style:     m.EmojiButtonStyle,
			Label:     role.Name,
			ReactData: fmt.Sprintf(roleReactData, role.ID),
		}
	}

	return m.Response{
		Type:        m.MessageResponse,
		Title:       fmt.Sprintf("%s Role Menu", title),
		Description: "Click the buttons below to add role",
		Buttons:     buttons,
	}
}

func (command *RoleMenu) HandleReaction(data m.CommandData, reaction string) m.Response {
	roleID := strings.Split(reaction, "-")[1]
	var userID string
	if data.User != nil {
		userID = data.User.ID
	} else {
		userID = data.Member.User.ID
	}

	return m.Response{
		Type:    m.AddRoleResponse,
		GuildID: data.GuildID,
		UserID:  userID,
		RoleID:  roleID,
	}
}
