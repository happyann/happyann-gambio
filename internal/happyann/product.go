package happyann

import (
	"time"
)

const DateFormat = "2006-01-02T15:04:05.000Z"

type CategoryData struct {
	Id           int64     `json:"id"`
	ParentId     int64     `json:"parent_id"`
	Title        string    `json:"title"`
	DateAdded    time.Time `json:"date_added"`
	LastModified time.Time `json:"last_modified"`
}

type ProductData struct {
	Id           int64          `json:"id"`
	Source       string         `json:"shop"`
	Title        string         `json:"title"`
	Url          string         `json:"url"`
	Images       []string       `json:"images"`
	PricePerUnit int            `json:"price_per_unit"`
	PriceUnit    string         `json:"unit_type"`
	UnitCount    float32        `json:"unit_count"`
	ProductTypes []string       `json:"product_type"`
	Description  string         `json:"description"`
	Categories   []CategoryData `json:"categories"`
	IsNew        bool           `json:"is_new"`
	IsRemoved    bool           `json:"is_removed"`
	IsActive     bool           `json:"is_active"`
	DateAdded    time.Time      `json:"date_added"`
	LastModified time.Time      `json:"last_modified"`
}
