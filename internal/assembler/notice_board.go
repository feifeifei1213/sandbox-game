package assembler

import "time"

type PlayerNoticeBoard struct {
	PinnedNotice *PlayerNoticeItem  `json:"pinnedNotice"`
	RecentList   []PlayerNoticeItem `json:"recentList"`
}

type PlayerNoticeItem struct {
	ID          int64     `json:"id"`
	Kind        string    `json:"kind"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Pinned      bool      `json:"pinned"`
	PublishedAt time.Time `json:"publishedAt"`
	YearNo      *int      `json:"yearNo,omitempty"`
	StageCode   *string   `json:"stageCode,omitempty"`
	Amount      *float64  `json:"amount,omitempty"`
}
