package model

type Search struct {
	OriginTable        string
	ID                 string
	Title              string
	HighlightedTitle   string
	Content            string
	HighlightedContent string
	Rank               float64
}
