package rkmcp

import (
	"time"
)

// SearchResult is the JSON object returned by search tools.
type SearchResult struct {
	Recalls    []Recall `json:"recalls"`
	Offset     int      `json:"offset"`
	NextOffset *int     `json:"nextOffset,omitempty"`
}

// Recalls is the historical wrapper; new responses also include offset fields.
type Recalls struct {
	Recalls []Recall `json:"recalls"`
}

type Recall struct {
	ID     string `json:"id"`
	Source string `json:"source"`

	Title                string `json:"title"`
	Description          string `json:"description"`
	DescriptionTruncated bool   `json:"descriptionTruncated,omitempty"`
	URL                  string `json:"url"`

	PublishedOn time.Time  `json:"publishedOn"`
	Extracted   *Extracted `json:"extracted,omitempty"`
}

type Extracted struct {
	Locations         []string           `json:"locations,omitempty"`
	Establishments    []string           `json:"establishments,omitempty"`
	Stores            []string           `json:"stores,omitempty"`
	Products          []ExtractedProduct `json:"products,omitempty"`
	ProductsTruncated bool               `json:"productsTruncated,omitempty"`
	Images            []string           `json:"images,omitempty"`
	Contact           *ExtractedContact  `json:"contact,omitempty"`
}

type ExtractedProduct struct {
	Name         string   `json:"productName,omitempty"`
	Packaging    string   `json:"packaging,omitempty"`
	LotCodes     []string `json:"lotCodes,omitempty"`
	ModelNumbers []string `json:"modelNumbers,omitempty"`
	UPCs         []string `json:"upcs,omitempty"`
	UseByDates   []string `json:"useByDates,omitempty"`
}

type ExtractedContact struct {
	Phones   []string `json:"phones,omitempty"`
	Emails   []string `json:"emails,omitempty"`
	Websites []string `json:"websites,omitempty"`
}

// SearchOptions are optional filters for search_product_recalls.
type SearchOptions struct {
	Query    string
	Source   string // cpsc, fdafoodsafety, FDAMedWatch, usda, nhtsa
	Since    string // YYYY-MM-DD
	Until    string // YYYY-MM-DD
	Location string
	Offset   int
	Limit    int
}

// IdentifierOptions are arguments for search_recalls_by_identifier.
type IdentifierOptions struct {
	UPC         string
	LotCode     string
	ModelNumber string
	ProductName string
	VIN         string
	Offset      int
	Limit       int
}
