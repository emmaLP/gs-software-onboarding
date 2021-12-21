package model

// Item represents the API response structure from HackerNew for an item as used for data storage
type Item struct {
	ID        int    `bson:"id" json:"id"`
	Type      string `bson:"type" json:"type"`
	Text      string `bson:"text" json:"text"`
	URL       string `bson:"url" json:"url"`
	Score     int    `bson:"score" json:"score"`
	Title     string `bson:"title" json:"title"`
	Time      int64  `bson:"time" json:"time"`
	CreatedBy string `bson:"by" json:"by"`
	Dead      bool   `bson:"dead" json:"dead"`
	Deleted   bool   `bson:"deleted" json:"deleted"`
}
