package blog

import "time"

type Article struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Summary     string    `json:"summary"`
	Body        string    `json:"body"`
	Authors     []string  `json:"authors"`
	PublishedAt time.Time `json:"published_at"`
	Tags        []string  `json:"tags"`
}
