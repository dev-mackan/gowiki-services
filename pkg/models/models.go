package models

import (
	"strings"
	"time"
)

type Page struct {
	PageId    uint      `json:"page_id"`
	Title     string    `json:"title"`
	LatestRev uint      `json:"latest_rev"`
	CreatedAt time.Time `json:"created_at"`
}

func (p *Page) DisplayTitle() string {
	return strings.ReplaceAll(p.Title, "_", " ")
}

type Revision struct {
	RevId     uint      `json:"rev_id"`
	PageId    uint      `json:"page_id"`
	TextId    uint      `json:"text_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Text struct {
	TextId    uint      `json:"text_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type PageCategory struct {
	CatId     uint      `json:"cat_id"`
	CatTitle  string    `json:"cat_title"`
	CatPages  []Page    `json:"cat_pages"`
	CreatedAt time.Time `json:"created_at"`
}

type PageBundle struct {
	Page     *Page     `json:"page"`
	Revision *Revision `json:"revision"`
	Text     *Text     `json:"text"`
}

type RevisionBundle struct {
	Revision Revision `json:"revision"`
	Text     Text     `json:"text"`
}
