package repository

import (
	"fmt"
	"gitlab/live/be-live-api/dto"
	"gitlab/live/be-live-api/model"
	"gitlab/live/be-live-api/utils"

	"gorm.io/gorm"
)

type StreamRepository struct {
	db *gorm.DB
}

func newStreamRepository(db *gorm.DB) *StreamRepository {
	return &StreamRepository{
		db: db,
	}
}

func (s *StreamRepository) PaginateStreamStatisticsData(page, limit int, cond *dto.StatisticsQuery) (*utils.PaginationModel[model.StreamAnalytics], error) {

	var query = s.db.Model(model.StreamAnalytics{}).Preload("Stream")
	if cond != nil && cond.Sort != "" && cond.SortBy != "" {
		if cond.SortBy != "duration" {
			if cond.SortBy == "title" {
				query = query.Joins("LEFT JOIN streams st ON st.id = stream_analytics.stream_id")
				query = query.Order(fmt.Sprintf("st.%s %s", cond.SortBy, cond.Sort))
			} else {

				query = query.Order(fmt.Sprintf("stream_analytics.%s %s", cond.SortBy, cond.Sort))
			}
		}
	} else {
		query = query.Order(fmt.Sprintf("stream_analytics.%s %s", "created_at", "DESC"))
	}
	pagination, err := utils.CreatePage[model.StreamAnalytics](query, page, limit)
	if err != nil {
		return nil, err
	}
	return utils.Create(pagination, page, limit)
}

func (s *StreamRepository) GetStatisticsTotalStream() (int64, int64, error) {
	var activeStream, totalStream int64

	if err := s.db.Model(model.Stream{}).Where("status=?", model.STARTED).Count(&activeStream).Error; err != nil {
		return 0, 0, err
	}
	if err := s.db.Model(model.Stream{}).Count(&totalStream).Error; err != nil {
		return 0, 0, err
	}
	return totalStream, activeStream, nil
}
