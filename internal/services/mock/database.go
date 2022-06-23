package mock

import (
	m "github.com/hmccarty/parca/internal/models"
)

type VerifyAttempt struct {
	Code    string
	GuildID string
}

type Poll struct {
	Title  string
	PartyA int
	PartyB int
}

type Bounty struct {
	Title   string
	Desc    string
	Link    string
	Claimed bool
}

type MockDbClient struct {
	// Mock data
	UserBalances          map[string]float64
	GuildSortedUsers      []string
	UserGuilds            []string
	GuildChannelCalendars map[string]map[string][]string
	GuildVerifyDomain     map[string]string
	GuildVerifyRole       map[string]string
	UserVerifyAttempt     map[string]VerifyAttempt
	Polls                 map[string]Poll
	Bounties              map[string]Bounty

	// Mock errors
	ErrGetUserBalance       error
	ErrSetUserBalance       error
	ErrGetBalancesFromGuild error
	ErrAddCalendar          error
	ErrGetCalendars         error
	ErrHasCalendar          error
	ErrRemoveCalendar       error
	ErrGetVerifyConfig      error
	ErrAddVerifyCode        error
	ErrGetVerifyCode        error
	ErrCreatePoll           error
	ErrGetPollTitle         error
	ErrAddPollVote          error
	ErrGetPollVote          error
	ErrCreateBounty         error
	ErrGetBounty            error
}

func (c *MockDbClient) GetUserBalance(userID string) (float64, error) {
	if balance, ok := c.UserBalances[userID]; ok {
		return balance, nil
	} else {
		return 0.0, c.ErrGetUserBalance
	}
}

func (c *MockDbClient) SetUserBalance(userID, guildID string, amt float64) error {
	if _, ok := c.UserBalances[userID]; ok {
		c.UserBalances[userID] = amt
		return nil
	} else {
		return c.ErrSetUserBalance
	}
}

func (c *MockDbClient) GetBalancesFromGuild(guildID string) ([]*m.BalanceEntry, error) {
	var balances []*m.BalanceEntry
	for _, userID := range c.GuildSortedUsers {
		if balance, ok := c.UserBalances[userID]; ok {
			balances = append(balances, &m.BalanceEntry{
				UserID: userID, Balance: balance,
			})
		} else {
			return nil, c.ErrGetBalancesFromGuild
		}
	}
	return balances, nil
}

func (c *MockDbClient) AddCalendar(calendarID, channelID, guildID string) error {
	if guildChannels, ok := c.GuildChannelCalendars[guildID]; ok {
		if calendarIDs, ok := guildChannels[channelID]; ok {
			c.GuildChannelCalendars[guildID][channelID] = append(calendarIDs, calendarID)
			return nil
		}
	}
	return c.ErrAddCalendar
}

func (c *MockDbClient) GetCalendars(channelID, guildID string) ([]string, error) {
	if guildChannels, ok := c.GuildChannelCalendars[guildID]; ok {
		if calendarIDs, ok := guildChannels[channelID]; ok {
			return calendarIDs, nil
		}
	}
	return nil, c.ErrGetCalendars
}

func (c *MockDbClient) HasCalendar(calendarID, channelID, guildID string) (bool, error) {
	if guildChannels, ok := c.GuildChannelCalendars[guildID]; ok {
		if calendarIDs, ok := guildChannels[channelID]; ok {
			for _, v := range calendarIDs {
				if v == calendarID {
					return true, nil
				}
			}
			return false, nil
		}
	}
	return false, c.ErrHasCalendar
}

func (c *MockDbClient) RemoveCalendar(calendarID, channelID, guildID string) error {
	if guildChannels, ok := c.GuildChannelCalendars[guildID]; ok {
		if calendarIDs, ok := guildChannels[channelID]; ok {
			for i, v := range calendarIDs {
				if v == calendarID {
					c.GuildChannelCalendars[guildID][channelID] = append(calendarIDs[:i], calendarIDs[i+1:]...)
					return nil
				}
			}
		}
	}
	return c.ErrRemoveCalendar
}

func (c *MockDbClient) AddVerifyConfig(domain, roleID, guildID string) error {
	c.GuildVerifyDomain[guildID] = domain
	c.GuildVerifyRole[guildID] = roleID
	return nil
}

func (c *MockDbClient) GetVerifyConfig(guildID string) (string, string, error) {
	domain, ok := c.GuildVerifyDomain[guildID]
	if !ok {
		return "", "", c.ErrGetVerifyConfig
	}

	roleID, ok := c.GuildVerifyRole[guildID]
	if !ok {
		return "", "", c.ErrGetVerifyConfig
	}

	return domain, roleID, nil
}

func (c *MockDbClient) AddVerifyCode(code, userID, guildID string) error {
	c.UserVerifyAttempt[userID] = VerifyAttempt{
		GuildID: guildID, Code: code,
	}
	return nil
}

func (c *MockDbClient) GetVerifyCode(userID string) (string, string, error) {
	if verifyAttempt, ok := c.UserVerifyAttempt[userID]; ok {
		return verifyAttempt.GuildID, verifyAttempt.Code, nil
	}
	return "", "", c.ErrGetVerifyCode
}

func (c *MockDbClient) CreatePoll(pollTitle, pollID string) error {
	if _, ok := c.Polls[pollID]; ok {
		return c.ErrCreatePoll
	}
	c.Polls[pollID] = Poll{
		Title: pollTitle, PartyA: 0, PartyB: 0,
	}
	return nil
}

func (c *MockDbClient) GetPollTitle(pollID string) (string, error) {
	if poll, ok := c.Polls[pollID]; ok {
		return poll.Title, nil
	}
	return "", c.ErrGetPollTitle
}

func (c *MockDbClient) AddPollVote(vote bool, pollID, userID string) error {
	if poll, ok := c.Polls[pollID]; ok {
		if vote {
			c.Polls[pollID] = Poll{
				Title: poll.Title, PartyA: poll.PartyA + 1, PartyB: poll.PartyB,
			}
		} else {
			c.Polls[pollID] = Poll{
				Title: poll.Title, PartyA: poll.PartyA, PartyB: poll.PartyB + 1,
			}
		}
		return nil
	} else {
		return c.ErrAddPollVote
	}
}

func (c *MockDbClient) GetPollVote(pollID string) (int, int, error) {
	if poll, ok := c.Polls[pollID]; ok {
		return poll.PartyA, poll.PartyB, nil
	}
	return 0, 0, c.ErrGetPollVote
}

func (c *MockDbClient) CreateBounty(bountyID, title, desc, link string) error {
	if _, ok := c.Bounties[bountyID]; ok {
		return c.ErrCreateBounty
	}
	c.Bounties[bountyID] = Bounty{
		Title: title, Desc: desc, Link: link,
	}
	return nil
}

func (c *MockDbClient) SetBountyAsClaimed(bountyID string) error {
	if bounty, ok := c.Bounties[bountyID]; ok {
		bounty.Claimed = true
		return nil
	} else {
		return c.ErrCreateBounty
	}
}

func (c *MockDbClient) WasBountyClaimed(bountyID string) (bool, error) {
	if bounty, ok := c.Bounties[bountyID]; ok {
		return bounty.Claimed, nil
	} else {
		return false, nil
	}
}

func (c *MockDbClient) GetBounty(bountyID string) (string, string, string, error) {
	if bounty, ok := c.Bounties[bountyID]; ok {
		return bounty.Title, bounty.Desc, bounty.Link, nil
	}
	return "", "", "", c.ErrGetBounty
}
