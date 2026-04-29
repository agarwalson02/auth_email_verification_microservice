package redis

import (
	"auth_service/config"
	"context"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

type Client struct {
	client *goredis.Client
}

func NewRedisClient(cfg *config.Config) (*Client,error) {

	client := goredis.NewClient(&goredis.Options{
		Addr:     cfg.Redis.Host + ":" + cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	return &Client{client: client},nil
}

func (c *Client) SetSession(ctx context.Context, sessionID, userID string, ttl time.Duration) error {
	return c.client.Set(ctx, sessionID, userID, ttl).Err()
}

func (c *Client) GetSession(ctx context.Context, sessionID string) (string, error) {
	return c.client.Get(ctx, sessionID).Result()
}

func (c *Client) DeleteSession(ctx context.Context, sessionID string) error {
	return c.client.Del(ctx, sessionID).Err()
}

func (c *Client) Close() error {
	return c.client.Close()
}
