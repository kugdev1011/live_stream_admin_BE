package service

import (
	"gitlab/live/be-live-api/dto"
	"gitlab/live/be-live-api/repository"
	"gitlab/live/be-live-api/utils"

	"github.com/redis/go-redis/v9"
)

type StreamService struct {
	repo  *repository.Repository
	redis *redis.Client
}

func newStreamService(repo *repository.Repository, redis *redis.Client) *StreamService {
	return &StreamService{
		repo:  repo,
		redis: redis,
	}

}

func (s *StreamService) GetStreamAnalyticsData(page, limit int) (*utils.PaginationModel[dto.LiveStreamRespDTO], error) {
	pagination, err := s.repo.Stream.PaginateStreamStatisticsData(page, limit)
	if err != nil {
		return nil, err
	}

	result := new(utils.PaginationModel[dto.LiveStreamRespDTO])
	result.BasePaginationModel = pagination.BasePaginationModel

	for _, v := range pagination.Page {
		var live_stream_dto = new(dto.LiveStreamRespDTO)
		live_stream_dto.Title = v.Stream.Title
		live_stream_dto.Description = v.Stream.Description
		live_stream_dto.Comments = v.Comments
		live_stream_dto.Likes = v.Likes
		live_stream_dto.VideoSize = utils.ConvertBytes(int64(v.VideoSize))
		live_stream_dto.Viewers = v.Views

		if !v.Stream.EndedAt.Valid && !v.Stream.StartedAt.Valid {
			endAt, _ := utils.ConvertDatetimeToTimestamp(v.Stream.EndedAt.String, utils.DATETIME_LAYOUT)
			startAt, _ := utils.ConvertDatetimeToTimestamp(v.Stream.StartedAt.String, utils.DATETIME_LAYOUT)
			live_stream_dto.Duration = utils.ConvertTimestampToDuration(endAt - startAt)
		}
		result.Page = append(result.Page, *live_stream_dto)
	}
	return result, nil
}
