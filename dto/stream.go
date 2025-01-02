package dto

import "gitlab/live/be-live-api/model"

type StreamRequest struct {
	Title             string             `json:"title" form:"title" validate:"required"`
	Description       string             `json:"description" form:"description" validate:"required"`
	UserID            uint               `json:"user_id" form:"user_id" validate:"required"`
	VideoFileName     string             `json:"-" form:"-"`
	ThumbnailFileName string             `json:"-" form:"-"`
	Status            model.StreamStatus `json:"status" form:"status" validate:"omitempty,oneof=pending started ended upcoming"`
	ScheduledAt       string             `json:"scheduled_at" form:"scheduled_at" validate:"required,datetime=2006-01-02 15:04:05.999 -0700"` //expect in utc
	CategoryIDs       []uint             `json:"category_ids" form:"category_ids" validate:"required,max=3,dive,required"`
}

type StreamAndStreamScheduleDto struct {
	Stream         *model.Stream         `json:"-"`
	ScheduleStream *model.ScheduleStream `json:"-" gorm:"-"`
}
