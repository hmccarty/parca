package general

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

func (*Poll) Options() []m.CommandOptionMetadata {
	return []m.CommandOptionMetadata{
		{
			Name:        "title",
			Description: "Name of your poll",
			Type:        m.StringOption,
			Required:    true,
		},
	}
}

func (cmd *Poll) Run(ctx m.ChatContext) error {
	if ctx.Options() != nil {
		opts := ctx.Options()
		title, err := opts[0].ToString()
		if err != nil {
			return err
		}

		invalid, err := regexp.MatchString(`[^?!0-9A-Za-z ]`, title)
		if invalid || err != nil {
			return ctx.Respond(m.Response{
				Type:        m.MessageResponse,
				Description: "Invalid title, please only use alpha-numeric characters",
			})
		}

		pollID := fmt.Sprintf("%d", rand.Intn(100000))

		client := cmd.createDbClient()
		err = client.CreatePoll(title, pollID)
		if err != nil {
			if err == m.ErrorPollIDAlreadyExists {
				for err == m.ErrorPollIDAlreadyExists {
					pollID = fmt.Sprintf("%d", rand.Intn(100000))
					err = client.CreatePoll(title, pollID)
				}
				if err != nil {
					return ctx.Respond(m.Response{
						Type:        m.MessageResponse,
						Description: "Failed to create poll please try again later",
					})
				}
			} else {
				return ctx.Respond(m.Response{
					Type:        m.MessageResponse,
					Description: "Failed to create poll please try again later",
				})
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

		return ctx.Respond(m.Response{
			Type:        m.MessageResponse,
			Title:       fmt.Sprintf(pollTitle, title),
			Description: fmt.Sprintf(pollDesc, 0, 0),
			Buttons:     buttons,
		})
	} else {
		msg := ctx.Message()
		reactComp := strings.Split(msg.Reaction, "-")
		pollID := reactComp[1]
		vote := reactComp[2] == "up"

		client := cmd.createDbClient()
		err := client.AddPollVote(vote, pollID, ctx.UserID())
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

		color := 0
		if yesCnt > noCnt {
			color = m.ColorGreen
		} else if yesCnt < noCnt {
			color = m.ColorRed
		}

		return ctx.Respond(m.Response{
			Type:        m.MessageEditResponse,
			Title:       fmt.Sprintf(pollTitle, title),
			Description: fmt.Sprintf(pollDesc, yesCnt, noCnt),
			Color:       color,
		})
	}
}
