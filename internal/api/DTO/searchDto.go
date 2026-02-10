package DTO

type SearchDto struct {
	Title string  `json:"title"`
	Url   string  `json:"url"`
	Score float64 `json:"score"`
}
