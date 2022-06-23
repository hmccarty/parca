package general

import (
	"fmt"
	"testing"

	m "github.com/hmccarty/parca/internal/models"
	"github.com/hmccarty/parca/internal/services/mock"
)

func createMockDbClient() m.DbClient {
	return &mock.MockDbClient{}
}

func TestBountyBasic(t *testing.T) {
	bounty := NewBountyCommand(5.0, createMockDbClient)
	ctx := mock.MockChatContext{}
	bounty.Run(&ctx)
	fmt.Println(ctx.GetResponse())
}
