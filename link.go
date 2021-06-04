package urlshortener

import "time"

type Link struct {
	FullForm   string        `json:"full_url"`
	ShortForm  string        `json:"short_url"`
	Expiration time.Duration `json:"expiration_time"`
	Redirects  uint64        `json:"redirects"`
}
