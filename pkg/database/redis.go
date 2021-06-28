package database

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

// DB structure represents a client to communicate with Redis.
type DB struct {
	Client *redis.Client
}

// NewDB connects to Redis and returns DB structure.
func NewDB(redisURL string) (*DB, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}
	rdb := redis.NewClient(opt)

	return &DB{Client: rdb}, nil
}

// GetLink gets link information from the database. Returns link object.
func (db *DB) GetLink(shortURL string) (*urlshortener.Link, error) {
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

// SetLink adds link information into the database.
// Creates two entries:
// 	key: short link - value: full link
//  key: full link  - value: amount of redirects
func (db *DB) SetLink(link *urlshortener.Link) error {
	ctx := context.Background()
	if err := db.Client.Set(
		ctx,
		link.ShortForm+":short",
		link.FullForm+":full",
		link.Expiration,
	).Err(); err != nil {
		return err
	}

	if err := db.Client.Set(
		ctx, link.FullForm+":full",
		0,
		link.Expiration,
	).Err(); err != nil {
		return err
	}

	return nil
}

// LinkExists checks if short link is already exists in the database.
func (db *DB) LinkExists(link string) bool {
	if _, err := db.Client.Get(
		context.Background(),
		link+":short",
	).Result(); errors.Is(err, redis.Nil) {
		return false
	}
	return true
}

// Redirect increments if user has used the short link.
func (db *DB) Redirect(link string) error {
	if _, err := db.Client.Incr(
		context.Background(),
		link+":full",
	).Result(); err != nil {
		return err
	}

	return nil
}
