package dto

import (
	"gitlab/live/be-live-api/model"
	"time"
)

type LiveStreamRespDTO struct {
	Title       string     `json:"title"`
	StreamID    uint       `json:"stream_id"`
	Description string     `json:"description"`
	VideoSize   int64      `json:"video_size"`
	Likes       uint       `json:"likes"`
	Viewers     uint       `json:"viewers"`
	Comments    uint       `json:"comments"`
	Duration    int64      `json:"duration"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
}

type LiveStatQuery struct {
	SortBy  string             `query:"sort_by" validate:"omitempty,oneof=total_viewers likes comments current_viewers"`
	Sort    string             `query:"sort" validate:"omitempty,oneof=DESC ASC"`
	Keyword string             `query:"keyword" validate:"omitempty"`
	Status  model.StreamStatus `query:"status" validate:"omitempty,oneof=pending started upcoming"`
	Page    uint               `query:"page" validate:"required,min=1"`
	Limit   uint               `query:"limit" validate:"required,min=1,max=20"`
}

type LiveStatRespDTO struct {
	StreamID       uint `json:"stream_id"`
	Likes          uint `json:"likes"`
	CurrentViewers uint `json:"current_viewers"`
	TotalViewers   uint `json:"total_viewers"`
	Comments       uint `json:"comments"`
}

type StatisticsTotalLiveStreamDTO struct {
	ActiveLiveStreams uint `json:"active_live_streams"`
	TotalLiveStreams  uint `json:"total_live_streams"`
}

type LiveCurrentViewers struct {
	StreamID uint `json:"stream_id"`
	Viewers  uint `json:"viewers"`
}

type StatisticsQuery struct {
	SortBy  string             `query:"sort_by" validate:"omitempty,oneof=title created_at views likes comments video_size duration stream_id id"`
	Sort    string             `query:"sort" validate:"omitempty,oneof=DESC ASC"`
	From    int64              `query:"from" validate:"omitempty"`
	To      int64              `query:"to" validate:"omitempty"`
	Keyword string             `query:"keyword" validate:"omitempty"`
	Status  model.StreamStatus `query:"status" validate:"omitempty,oneof=pending started ended upcoming"`
	Page    uint               `query:"page" validate:"required,min=1"`
	Limit   uint               `query:"limit" validate:"required,min=1,max=20"`
}

type LiveStreamBroadCastQueryDTO struct {
	SortBy          string               `query:"sort_by" validate:"omitempty,oneof=title started_at ended_at views likes comments video_size duration"`
	Sort            string               `query:"sort" validate:"omitempty,oneof=DESC ASC"`
	Status          []model.StreamStatus `query:"status" validate:"omitempty"`
	Type            model.StreamType     `query:"type" validate:"omitempty,oneof=camera software"`
	Category        string               `query:"category" validate:"omitempty"`
	FromStartedTime int64                `query:"from_started_time" validate:"omitempty"`
	EndStartedTime  int64                `query:"end_started_time" validate:"omitempty"`
	FromEndedTime   int64                `query:"from_ended_time" validate:"omitempty"`
	EndEndedTime    int64                `query:"end_ended_time" validate:"omitempty"`
	Keyword         string               `query:"keyword" validate:"omitempty"`
	Page            uint                 `query:"page" validate:"required,min=1"`
	Limit           uint                 `query:"limit" validate:"required,min=1,max=20"`
}

type LiveStreamBroadCastDTO struct {
	ID                 int                `json:"id,omitempty"`
	Title              string             `json:"title,omitempty"`
	Description        string             `json:"description,omitempty"`
	Status             model.StreamStatus `json:"status,omitempty"`
	PushURL            string             `json:"push_url,omitempty"`      // generated from streaming server
	BroadcastURL       string             `json:"broadcast_url,omitempty"` // generated from web
	StreamType         model.StreamType   `json:"stream_type,omitempty"`
	ThumbnailFileName  string             `json:"thumbnail_file_name,omitempty"`
	StartedAt          *time.Time         `json:"started_at,omitempty"`
	EndedAt            *time.Time         `json:"ended_at,omitempty"`
	User               *UserResponseDTO   `json:"user,omitempty"`
	Category           string             `json:"category,omitempty"`
	LiveStreamAnalytic *LiveStreamRespDTO `json:"live_stream_analytic"`
}
