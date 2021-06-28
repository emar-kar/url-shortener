package service

import (
	"fmt"

	"github.com/teris-io/shortid"

	"github.com/emar-kar/urlshortener"
)

// LinkManager represents an interface for database interactions.
// Current project uses Redis as a storage, but it can be changed
// to another one.
type LinkManager interface {
	SetLink(*urlshortener.Link) error
	GetLink(string) (*urlshortener.Link, error)
	LinkExists(string) bool
	Redirect(string) error
}

// Service represents communications with the database and
// implements short URL generator.
type Service struct {
	LinkManager
}

// NewService creates an object with database which implements
// LinkManager interface.
func NewService(db LinkManager) *Service {
	return &Service{
		LinkManager: db,
	}
}

// GenerateShortURL creates unique identifier for the short link.
func (s *Service) GenerateShortURL(host string) (string, error) {
	for {
		hash, err := shortid.Generate()
		if err != nil {
			return "", err
		}
		shortURL := fmt.Sprintf("%s/%s", host, hash)
		if !s.LinkManager.LinkExists(shortURL) {
			return shortURL, nil
		}
	}
}
