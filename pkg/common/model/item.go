package model

import pb "github.com/emmaLP/gs-software-onboarding/pkg/grpc/proto"

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

func PItemToItem(item *pb.Item) Item {
	return Item{
		ID:        int(item.Id),
		Type:      item.Type,
		Text:      item.Text,
		URL:       item.Url,
		Score:     int(item.Score),
		Title:     item.Title,
		Time:      item.Time,
		CreatedBy: item.CreatedBy,
		Dead:      item.Dead,
		Deleted:   item.Deleted,
	}
}

func ItemToPItem(item Item) *pb.Item {
	return &pb.Item{
		Id:        int32(item.ID),
		Type:      item.Type,
		Text:      item.Text,
		Url:       item.URL,
		Score:     int64(item.Score),
		Title:     item.Title,
		Time:      item.Time,
		CreatedBy: item.CreatedBy,
		Dead:      item.Dead,
		Deleted:   item.Deleted,
	}
}
