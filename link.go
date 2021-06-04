package urlshortener

import "time"

type Link struct {
	FullForm   string
	ShortForm  string
	Expiration time.Duration
	Redirects  uint64
}
