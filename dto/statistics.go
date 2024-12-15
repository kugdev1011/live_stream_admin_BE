package dto

import "time"

type LiveStreamRespDTO struct {
	Title       string     `json:"title,omitempty"`
	Description string     `json:"description,omitempty"`
	VideoSize   int64      `json:"video_size,omitempty"`
	Likes       uint       `json:"likes,omitempty"`
	Viewers     uint       `json:"viewers,omitempty"`
	Comments    uint       `json:"comments,omitempty"`
	Duration    int64      `json:"duration,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
}

type StatisticsTotalLiveStreamDTO struct {
	ActiveLiveStreams uint `json:"active_live_streams"`
	TotalLiveStreams  uint `json:"total_live_streams"`
}

type StatisticsQuery struct {
	SortBy string `query:"sort_by" validate:"omitempty,oneof=title created_at views likes comments video_size duration stream_id id"`
	Sort   string `query:"sort" validate:"omitempty,oneof=DESC ASC"`
}
