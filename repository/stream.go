package repository

import (
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

func (s *StreamRepository) PaginateStreamStatisticsData(page, limit int) (*utils.PaginationModel[model.StreamAnalytics], error) {

	var query = s.db.Model(model.StreamAnalytics{}).Preload("Stream")
	pagination, err := utils.CreatePage[model.StreamAnalytics](query, page, limit)
	if err != nil {
		return nil, err
	}
	return utils.Create(pagination, page, limit)
}
