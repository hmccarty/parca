package redis

import (
	"errors"
	"fmt"

	"context"
	"strconv"

	"github.com/go-redis/redis/v8"
	m "github.com/hmccarty/parca/internal/models"
	c "github.com/hmccarty/parca/internal/services/config"
)

const (
	// Currency
	guildBalanceSortedKey = "guild:%s:balances"
	userBalanceKey        = "user:%s:balance"
	userGuildsKey         = "user:%s:guilds"

	// Calendar
	calendarKey = "guild:%s:channel:%s:calendar"

	// Verification
	verifyConfigKey = "verify:guild:%s"
	verifyCodeKey   = "verify:user:%s"

	// General
	pollTitleKey   = "poll:%s:title"
	pollYesVoteKey = "poll:%s:yes"
	pollNoVoteKey  = "poll:%s:no"

	bountyKey = "bounty:%s"
)

func OpenRedisClient(config *c.Config) m.DbClient {
	return &RedisClient{
		client: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		}),
	}
}

type RedisClient struct {
	client *redis.Client
}

func (r *RedisClient) registerUser(userID string, guildID string) error {
	ctx := context.Background()
	key := fmt.Sprintf(userGuildsKey, userID)
	return r.client.SAdd(ctx, key, guildID).Err()
}

func (r *RedisClient) GetUserBalance(userID string) (float64, error) {
	ctx := context.Background()
	key := fmt.Sprintf(userBalanceKey, userID)
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	return strconv.ParseFloat(val, 64)
}

func (r *RedisClient) SetUserBalance(userID string, guildID string, balance float64) error {
	ctx := context.Background()

	key := fmt.Sprintf(userBalanceKey, userID)
	err := r.client.Set(ctx, key, balance, 0).Err()
	if err != nil {
		return err
	}

	r.registerUser(userID, guildID)

	key = fmt.Sprintf(userGuildsKey, userID)
	var guilds []string
	guilds, err = r.client.SMembers(ctx, key).Result()
	if err != nil {
		return nil
	}

	for _, guild := range guilds {
		key = fmt.Sprintf(guildBalanceSortedKey, guild)
		member := &redis.Z{
			Score:  balance,
			Member: userID,
		}
		err = r.client.ZAdd(ctx, key, member).Err()
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *RedisClient) GetBalancesFromGuild(guildID string) ([]*m.BalanceEntry, error) {
	ctx := context.Background()

	key := fmt.Sprintf(guildBalanceSortedKey, guildID)
	vals, err := r.client.ZRevRangeByScoreWithScores(ctx, key, &redis.ZRangeBy{
		Min:   "-inf",
		Max:   "+inf",
		Count: 10,
	}).Result()
	if err != nil {
		return nil, err
	}

	balances := make([]*m.BalanceEntry, len(vals))
	for i, v := range vals {
		balances[i] = &m.BalanceEntry{
			UserID:  v.Member.(string),
			Balance: v.Score,
		}
	}

	return balances, nil
}

func (r *RedisClient) AddCalendar(calendarID, channelID, guildID string) error {
	ctx := context.Background()

	key := fmt.Sprintf(calendarKey, guildID, channelID)
	return r.client.SAdd(ctx, key, calendarID).Err()
}

func (r *RedisClient) GetCalendars(channelID, guildID string) ([]string, error) {
	ctx := context.Background()

	key := fmt.Sprintf(calendarKey, guildID, channelID)
	return r.client.SMembers(ctx, key).Result()
}

func (r *RedisClient) HasCalendar(calendarID, channelID, guildID string) (bool, error) {
	ctx := context.Background()

	key := fmt.Sprintf(calendarKey, guildID, channelID)
	return r.client.SIsMember(ctx, key, calendarID).Result()
}

func (r *RedisClient) RemoveCalendar(calendarID, channelID, guildID string) error {
	ctx := context.Background()

	key := fmt.Sprintf(calendarKey, guildID, channelID)
	return r.client.SRem(ctx, key, calendarID).Err()
}

func (r *RedisClient) AddVerifyConfig(domain, roleID, guildID string) error {
	ctx := context.Background()

	key := fmt.Sprintf(verifyConfigKey, guildID)
	return r.client.HSet(ctx, key, "domain", domain, "roleID", roleID).Err()
}

func (r *RedisClient) GetVerifyConfig(guildID string) (string, string, error) {
	ctx := context.Background()

	key := fmt.Sprintf(verifyConfigKey, guildID)
	value, err := r.client.HMGet(ctx, key, "domain", "roleID").Result()
	if err != nil {
		return "", "", err
	} else if value[0] == nil || value[1] == nil {
		return "", "", errors.New("failed to collect from redis")
	}

	return value[0].(string), value[1].(string), err
}

func (r *RedisClient) AddVerifyCode(code, userID, guildID string) error {
	ctx := context.Background()

	key := fmt.Sprintf(verifyCodeKey, userID)
	return r.client.HSet(ctx, key, "code", code, "guildID", guildID).Err()
}

func (r *RedisClient) GetVerifyCode(userID string) (string, string, error) {
	ctx := context.Background()

	key := fmt.Sprintf(verifyCodeKey, userID)
	value, err := r.client.HMGet(ctx, key, "code", "guildID").Result()

	if err != nil {
		return "", "", err
	} else if len(value) != 2 {
		return "", "", errors.New("failed to collect from redis")
	} else if value[0] == nil || value[1] == nil {
		return "", "", nil
	}
	return value[0].(string), value[1].(string), err
}

func (r *RedisClient) CreatePoll(pollTitle, pollID string) error {
	ctx := context.Background()

	titleKey := fmt.Sprintf(pollTitleKey, pollID)
	_, err := r.client.Get(ctx, titleKey).Result()
	if err == nil {
		return m.ErrorPollIDAlreadyExists
	}
	r.client.Set(ctx, titleKey, pollTitle, 0)

	return nil
}

func (r *RedisClient) GetPollTitle(pollID string) (string, error) {
	ctx := context.Background()

	titleKey := fmt.Sprintf(pollTitleKey, pollID)
	return r.client.Get(ctx, titleKey).Result()
}

func (r *RedisClient) AddPollVote(vote bool, pollID, userID string) error {
	ctx := context.Background()

	var newVoteKey string
	var oldVoteKey string
	if vote == true {
		newVoteKey = fmt.Sprintf(pollYesVoteKey, pollID)
		oldVoteKey = fmt.Sprintf(pollNoVoteKey, pollID)
	} else {
		newVoteKey = fmt.Sprintf(pollNoVoteKey, pollID)
		oldVoteKey = fmt.Sprintf(pollYesVoteKey, pollID)
	}

	alreadyVoted, err := r.client.SIsMember(ctx, newVoteKey, userID).Result()
	if err != nil {
		return err
	} else if alreadyVoted {
		return m.ErrorUserAlreadyVoted
	}

	hasVoted, err := r.client.SIsMember(ctx, oldVoteKey, userID).Result()
	if err != nil {
		return err
	} else if hasVoted {
		wasRemoved, err := r.client.SRem(ctx, oldVoteKey, userID).Result()
		if err != nil {
			return err
		} else if wasRemoved != 1 {
			return m.ErrorUnableToRemoveVoter
		}
	}
	return r.client.SAdd(ctx, newVoteKey, userID).Err()
}

func (r *RedisClient) GetPollVote(pollID string) (int, int, error) {
	ctx := context.Background()

	yesKey := fmt.Sprintf(pollYesVoteKey, pollID)
	yesUsers, err := r.client.SMembers(ctx, yesKey).Result()
	if err != nil {
		return 0, 0, m.ErrorPollIDDoesntExist
	}

	noKey := fmt.Sprintf(pollNoVoteKey, pollID)
	noUsers, err := r.client.SMembers(ctx, noKey).Result()
	if err != nil {
		return 0, 0, err
	}

	return len(yesUsers), len(noUsers), nil
}

func (r *RedisClient) CreateBounty(bountyID, title, desc, link string) error {
	ctx := context.Background()

	key := fmt.Sprintf(bountyKey, bountyID)
	if r.client.HLen(ctx, key).Val() > 0 {
		return m.ErrorBountyIDAlreadyExists
	}

	return r.client.HMSet(ctx, key, "title", title, "desc", desc, "link", link).Err()
}

func (r *RedisClient) GetBounty(bountyID string) (string, string, string, error) {
	ctx := context.Background()

	key := fmt.Sprintf(bountyKey, bountyID)
	values, err := r.client.HMGet(ctx, key,
		"title", "desc", "link").Result()
	if err != nil {
		return "", "", "", err
	} else if len(values) != 3 {
		return "", "", "", m.ErrorBountyIDMissingValues
	}

	return values[0].(string), values[1].(string), values[2].(string), nil
}
