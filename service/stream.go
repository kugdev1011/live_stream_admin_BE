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

func (s *StreamService) GetStatisticsTotalLiveStreamData() (*dto.StatisticsTotalLiveStreamDTO, error) {
	total, active, err := s.repo.Stream.GetStatisticsTotalStream()
	if err != nil {
		return nil, err
	}
	return &dto.StatisticsTotalLiveStreamDTO{ActiveLiveStreams: uint(active), TotalLiveStreams: uint(total)}, err
}

func (s *StreamService) GetStreamAnalyticsData(page, limit int, req *dto.StatisticsQuery) (*utils.PaginationModel[dto.LiveStreamRespDTO], error) {
	pagination, err := s.repo.Stream.PaginateStreamStatisticsData(page, limit, req)
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
		live_stream_dto.VideoSize = int64(v.VideoSize)
		live_stream_dto.Viewers = v.Views
		live_stream_dto.CreatedAt = &v.CreatedAt

		if v.Stream.EndedAt.Valid && v.Stream.StartedAt.Valid {
			endAt, errEndAt := utils.ConvertDatetimeToTimestamp(v.Stream.EndedAt.String, utils.DATETIME_LAYOUT)
			startAt, errStartAt := utils.ConvertDatetimeToTimestamp(v.Stream.StartedAt.String, utils.DATETIME_LAYOUT)
			if errEndAt != nil || errStartAt != nil {
				live_stream_dto.Duration = 0
			}
			live_stream_dto.Duration = int64(endAt.Sub(*startAt))
		}
		result.Page = append(result.Page, *live_stream_dto)
	}
	return result, nil
}
