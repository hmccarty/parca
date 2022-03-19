package currency

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

func (*Status) Options() []m.CommandOption {
	return []m.CommandOption{}
}

func (*Status) Run(_ m.CommandData, _ []m.CommandOption) m.Response {
	return m.Response{
		Type:        m.MessageResponse,
		Description: "Never been better",
	}
}

func (*Status) HandleReaction(_ m.CommandData, _ string) m.Response {
	return m.Response{
		Type:        m.MessageResponse,
		Description: "Not expecting a reaction",
	}
}
