package currency

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"

	m "github.com/hmccarty/parca/internal/models"
)

const (
	pollTitle         = "Poll: %s"
	pollDesc          = "Vote using reactions below\n Current tally: %d - %d"
	pollUpReactData   = "poll-%s-up"
	pollDownReactData = "poll-%s-down"
)

type Poll struct {
	createDbClient func() m.DbClient
}

func NewPollCommand(createDbClient func() m.DbClient) m.Command {
	return &Poll{
		createDbClient: createDbClient,
	}
}

func (*Poll) Name() string {
	return "poll"
}

func (*Poll) Description() string {
	return "Initiates a democratic process"
}

func (*Poll) Options() []m.CommandOption {
	return []m.CommandOption{
		{
			Name:        "title",
			Description: "Name of your poll",
			Type:        m.StringOption,
			Required:    true,
		},
	}
}

func (command *Poll) Run(_ m.CommandData, opts []m.CommandOption) m.Response {
	title := opts[0].Value.(string)
	invalid, err := regexp.MatchString(`[^?!0-9A-Za-z ]`, title)
	if invalid || err != nil {
		return m.Response{
			Type:        m.MessageResponse,
			Description: "Invalid title, please only use alpha-numeric characters",
		}
	}

	pollID := fmt.Sprintf("%d", rand.Intn(100000))

	client := command.createDbClient()
	err = client.CreatePoll(title, pollID)
	if err != nil {
		if err == m.ErrorPollIDAlreadyExists {
			for err == m.ErrorPollIDAlreadyExists {
				pollID = fmt.Sprintf("%d", rand.Intn(100000))
				err = client.CreatePoll(title, pollID)
			}
			if err != nil {
				return m.Response{
					Type:        m.MessageResponse,
					Description: "Failed to create poll please try again later",
				}
			}
		} else {
			return m.Response{
				Type:        m.MessageResponse,
				Description: "Failed to create poll please try again later",
			}
		}
	}

	buttons := []m.ResponseButton{
		{
			Style:     m.PrimaryButtonStyle,
			Emoji:     m.ThumbsUpEmoji,
			ReactData: fmt.Sprintf(pollUpReactData, pollID),
		},
		{
			Style:     m.SecondaryButtonStyle,
			Emoji:     m.ThumbsDownEmoji,
			ReactData: fmt.Sprintf(pollDownReactData, pollID),
		},
	}

	return m.Response{
		Type:        m.MessageResponse,
		Title:       fmt.Sprintf(pollTitle, title),
		Description: fmt.Sprintf(pollDesc, 0, 0),
		Buttons:     buttons,
	}
}

func (command *Poll) HandleReaction(data m.CommandData, reactData string) m.Response {
	reactComp := strings.Split(reactData, "-")
	pollID := reactComp[1]
	vote := reactComp[2] == "up"

	var userID string
	if data.User != nil {
		userID = data.User.ID
	} else {
		userID = data.Member.User.ID
	}

	client := command.createDbClient()
	err := client.AddPollVote(vote, pollID, userID)
	if err != nil {
		fmt.Println(err)
	}

	yesCnt, noCnt, err := client.GetPollVote(pollID)
	if err != nil {
		fmt.Println(err)
	}

	title, err := client.GetPollTitle(pollID)
	if err != nil {
		fmt.Println(err)
	}

	return m.Response{
		Type:        m.MessageEditResponse,
		Title:       fmt.Sprintf(pollTitle, title),
		Description: fmt.Sprintf(pollDesc, yesCnt, noCnt),
	}
}
