package model

// Item represents the API response structure from HackerNew for an item
type Item struct {
	ID        int    `json:"id"`
	Type      string `json:"type"`
	Text      string `json:"text"`
	URL       string `json:"url"`
	Score     int    `json:"score"`
	Title     string `json:"title"`
	Time      int64  `json:"time"`
	CreatedBy string `json:"by"`
	Dead      bool   `json:"dead"`
	Deleted   bool   `json:"deleted"`
}
