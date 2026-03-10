package rkmcp

import (
	"time"
)

type Recalls struct {
	Recalls []Recall `json:"recalls"`
}

type Recall struct {
	ID     string `json:"id"`
	Source string `json:"source"`

	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`

	PublishedOn time.Time `json:"publishedOn"`
}
