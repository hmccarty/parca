package general

import (
	m "github.com/hmccarty/parca/internal/models"
)

type Status struct{}

func NewStatusCommand() m.Command {
	return &Status{}
}

func (*Status) Name() string {
	return "status"
}

func (*Status) Description() string {
	return "Prints status of PARCA"
}

func (*Status) Options() []m.CommandOptionMetadata {
	return []m.CommandOptionMetadata{}
}

func (*Status) Run(ctx m.CommandContext) error {
	return ctx.Respond(m.Response{
		Type:        m.MessageResponse,
		Description: "Better than ever",
	})
}
