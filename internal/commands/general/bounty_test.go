package general

import (
	"testing"

	m "github.com/hmccarty/parca/internal/models"
	"github.com/hmccarty/parca/internal/services/mock"
)

func createMockDbClient() m.DbClient {
	return &mock.MockDbClient{}
}

func TestBountyBasic(t *testing.T) {
	ctx := mock.MockChatContext{}
	ctx.SetOptions([]m.CommandOption{
		{
			Metadata: m.CommandOptionMetadata{
				Type: m.StringOption,
				Name: "title",
			},
			Value: "Basic Test",
		},
		{
			Metadata: m.CommandOptionMetadata{
				Type: m.StringOption,
				Name: "description",
			},
			Value: "Description",
		},
	})

	bounty := NewBountyCommand(5.0, createMockDbClient)
	err := bounty.Run(&ctx)
	if err != nil {
		t.Errorf("Encountered unexpected error: %v", err)
	}

	response := ctx.GetResponse()
	if response.Type != m.AckResponse {
		t.Errorf("Responded with incorrect error type")
	}
}
