package redis

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/go-redis/redis/v8"

	"github.com/emar-kar/urlshortener"
)

type LinkType string

func (lt *LinkType) String() string {
	return string(*lt)
}

func (lt *LinkType) TrimShort() string {
	return strings.TrimSuffix(lt.String(), ":short")
}

func (lt *LinkType) TrimFull() string {
	return strings.TrimSuffix(lt.String(), ":full")
}

func (lt *LinkType) Short() string {
	return lt.String() + ":short"
}

func (lt *LinkType) Full() string {
	return lt.String() + ":full"
}

// Redis keys:
// url:short
// url:full

type DB struct {
	Client *redis.Client
}

func NewDB() *DB {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return &DB{Client: rdb}
}

func (db *DB) Get(shortURL string) (*urlshortener.Link, error) {
	link := &urlshortener.Link{ShortForm: shortURL}

	redisShortLink := LinkType(shortURL)
	ctx := context.Background()
	redisFullLink, err := db.Client.Get(ctx, redisShortLink.Short()).Result()
	if err != nil {
		return nil, fmt.Errorf("cannot get %s full URL: %w", shortURL, err)
	}
	fullLink := LinkType(redisFullLink)
	link.FullForm = fullLink.TrimFull()

	exp, err := db.Client.TTL(ctx, redisShortLink.Short()).Result()
	if err != nil {
		return nil, fmt.Errorf("cannot get %s TTL: %w", shortURL, err)
	}
	link.Expiration = exp

	redirects, err := db.Client.Get(ctx, redisFullLink).Uint64()
	if err != nil {
		return nil, fmt.Errorf("cannot get %s redirects: %w", fullLink, err)
	}
	link.Redirects = redirects

	return link, nil
}

func (db *DB) Set(link *urlshortener.Link) error {
	ctx := context.Background()
	redisShortLink := LinkType(link.ShortForm)
	redisFullLink := LinkType(link.FullForm)
	if err := db.Client.Set(ctx, redisShortLink.Short(), redisFullLink.Full(), link.Expiration).Err(); err != nil {
		return err
	}

	if err := db.Client.Set(ctx, redisFullLink.Full(), 0, link.Expiration).Err(); err != nil {
		return err
	}

	return nil
}

func (db *DB) Exist(shortURL string) bool {
	redisShortLink := LinkType(shortURL)
	if _, err := db.Client.Get(context.Background(), redisShortLink.Short()).Result(); errors.Is(err, redis.Nil) {
		return false
	}
	return true
}

func (db *DB) Redirect(fullURL string) error {
	redisFullLink := LinkType(fullURL)
	ctx := context.Background()
	red, err := db.Client.Get(ctx, redisFullLink.Full()).Uint64()
	if err != nil {
		return err
	}

	exp, err := db.Client.TTL(ctx, redisFullLink.Full()).Result()
	if err != nil {
		return err
	}

	red++
	if err := db.Client.Set(ctx, redisFullLink.Full(), red, exp).Err(); err != nil {
		return err
	}

	return nil
}
