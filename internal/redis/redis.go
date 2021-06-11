package redis

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/go-redis/redis/v8"

	"github.com/emar-kar/urlshortener"
)

// Redis keys:
// url:short
// url:full

type DB struct {
	Client *redis.Client
}

func NewDB() *DB {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	return &DB{Client: rdb}
}

func (db *DB) Get(shortURL string) (*urlshortener.Link, error) {
	link := &urlshortener.Link{ShortForm: shortURL}

	ctx := context.Background()
	redisFullLink, err := db.Client.Get(ctx, shortURL+":short").Result()
	if err != nil {
		return nil, fmt.Errorf("cannot get %s full URL: %w", shortURL, err)
	}
	link.FullForm = strings.TrimSuffix(redisFullLink, ":full")

	exp, err := db.Client.TTL(ctx, shortURL+":short").Result()
	if err != nil {
		return nil, fmt.Errorf("cannot get %s TTL: %w", shortURL, err)
	}
	link.Expiration = exp

	redirects, err := db.Client.Get(ctx, redisFullLink).Uint64()
	if err != nil {
		return nil, fmt.Errorf("cannot get %s redirects: %w", link.FullForm, err)
	}
	link.Redirects = redirects

	return link, nil
}

func (db *DB) Set(link *urlshortener.Link) error {
	ctx := context.Background()
	if err := db.Client.Set(ctx, link.ShortForm+":short", link.FullForm+":full", link.Expiration).Err(); err != nil {
		return err
	}

	if err := db.Client.Set(ctx, link.FullForm+":full", 0, link.Expiration).Err(); err != nil {
		return err
	}

	return nil
}

func (db *DB) Exist(link string) bool {
	if _, err := db.Client.Get(context.Background(), link+":short").Result(); errors.Is(err, redis.Nil) {
		return false
	}
	return true
}

func (db *DB) Redirect(link string) error {
	ctx := context.Background()
	red, err := db.Client.Get(ctx, link+":full").Uint64()
	if err != nil {
		return err
	}

	exp, err := db.Client.TTL(ctx, link+":full").Result()
	if err != nil {
		return err
	}

	red++
	if err := db.Client.Set(ctx, link+":full", red, exp).Err(); err != nil {
		return err
	}

	return nil
}
