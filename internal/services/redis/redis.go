package redis

import (
	"fmt"

	"context"
	"strconv"

	"github.com/go-redis/redis/v8"
	m "github.com/hmccarty/arc-assistant/internal/models"
	c "github.com/hmccarty/arc-assistant/internal/services/config"
)

const (
	guildBalanceSortedKey = "guild:%s:balances"
	userBalanceKey        = "user:%s:balance"
	userGuildsKey         = "user:%s:guilds"
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
	if err != nil {
		return 0, err
	}

	balance, err := strconv.ParseFloat(val, 64)
	return balance, nil
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
