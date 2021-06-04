package service

import (
	"fmt"

	"github.com/teris-io/shortid"

	"github.com/emar-kar/urlshortener"
)

type Database interface {
	Set(*urlshortener.Link) error
	Get(string) (*urlshortener.Link, error)
	Exist(string) bool
	Redirect(string) error
}

type Service struct {
	Database
}

func NewService(db Database) *Service {
	return &Service{
		Database: db,
	}
}

func (s *Service) GenerateShortURL(host string) (string, error) {
	for {
		hash, err := shortid.Generate()
		if err != nil {
			return "", err
		}
		shortURL := fmt.Sprintf("%s/%s", host, hash)
		if !s.Database.Exist(shortURL) {
			return shortURL, nil
		}
	}
}
