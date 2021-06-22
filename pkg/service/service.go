package service

import (
	"fmt"

	"github.com/teris-io/shortid"

	"github.com/emar-kar/urlshortener"
)

type LinkManager interface {
	SetLink(*urlshortener.Link) error
	GetLink(string) (*urlshortener.Link, error)
	LinkExists(string) bool
	Redirect(string) error
}

type Service struct {
	LinkManager
}

func NewService(db LinkManager) *Service {
	return &Service{
		LinkManager: db,
	}
}

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
