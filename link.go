package urlshortener

import "time"

// Link structure represents the main object of the service.
type Link struct {
	FullForm   string        `json:"full_url"`
	ShortForm  string        `json:"short_url"`
	Expiration time.Duration `json:"expiration_time"`
	Redirects  uint64        `json:"redirects"`
}
