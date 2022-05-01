package redis

import (
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
)

const TokenCollection = "token"

type Client struct {
	Pool *redis.Pool
}

func NewClient() *Client {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "6379"
	}

	addr := host + ":" + port
	pool := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		// Dial or DialContext must be set. When both are set, DialContext takes precedence over Dial.
		Dial: func() (redis.Conn, error) { return redis.Dial("tcp", addr) },
	}
	return &Client{Pool: pool}
}

func (c *Client) CloseSession() {
	_ = c.Pool.Close()
}

func (c *Client) AddToken(token string, expiry string) error {
	conn := c.Pool.Get()
	defer conn.Close()

	_, err := conn.Do("HSET", TokenCollection, token, expiry)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) UpdateToken(token string, expiry string) error {
	conn := c.Pool.Get()
	defer conn.Close()

	_, err := conn.Do("HSET", TokenCollection, token, expiry)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) TokenExists(token string) (bool, error) {
	conn := c.Pool.Get()
	defer conn.Close()

	res, err := redis.Bool(conn.Do("HEXISTS", TokenCollection, token))
	if err != nil {
		return false, err
	}

	return res, nil
}

func (c *Client) ExpiryByToken(token string) (string, error) {
	conn := c.Pool.Get()
	defer conn.Close()

	expiry, err := redis.String(conn.Do("HGET", TokenCollection, token))
	if err != nil {
		return "", err
	}

	return expiry, nil
}

func (c *Client) AllTokens() (map[string]string, error) {
	conn := c.Pool.Get()
	defer conn.Close()

	res, err := redis.Strings(conn.Do("HGETALL", TokenCollection))
	if err != nil {
		return nil, err
	}

	i := 0
	m := make(map[string]string)
	for i < len(res) {
		m[res[i]] = res[i+1]
		i += 2
	}

	return m, nil
}
