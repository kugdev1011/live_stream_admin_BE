package datasource

import (
	"fmt"
	"gitlab/live/be-live-api/conf"
	"gitlab/live/be-live-api/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func LoadDB() (*gorm.DB, error) {
	dbConfig := conf.GetDatabaseConfig()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=UTC", dbConfig.Host, dbConfig.User, dbConfig.Pass, dbConfig.Name, dbConfig.Port)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	db.AutoMigrate(
		&model.Role{},
		&model.User{},
		&model.AdminLog{},
		&model.BlockedList{},
		&model.Like{},
		&model.Comment{},
		&model.Share{},
		&model.StreamAnalytics{},
		&model.Subscription{},
		&model.Notification{},
		&model.Stream{},
	)

	return db, nil

}
