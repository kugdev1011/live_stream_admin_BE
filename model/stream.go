package model

import "time"

type StreamStatus string

const (
	STARTED  StreamStatus = "started"
	ENDED    StreamStatus = "ended"
	UPCOMING StreamStatus = "upcoming"
)

type LikeEmoteType string

type StreamType string

const (
	CAMERASTREAM   StreamType = "camera"
	SOFTWARESTREAM StreamType = "software" // like obs
)

const (
	LikeEmoteTypeLike    LikeEmoteType = "like"
	LikeEmoteTypeDislike LikeEmoteType = "dislike"
	LikeEmoteTypeLaugh   LikeEmoteType = "laugh"
	LikeEmoteTypeSad     LikeEmoteType = "sad"
	LikeEmoteTypeWow     LikeEmoteType = "wow"
	LikeEmoteTypeHeart   LikeEmoteType = "heart"
)

type Stream struct {
	StreamID    uint         `gorm:"primaryKey"`
	UserID      uint         `gorm:"not null"`
	Title       string       `gorm:"type:varchar(100);not null"`
	Description string       `gorm:"type:text"`
	Status      StreamStatus `gorm:"type:varchar(50);not null"`
	// StreamURL    string       `gorm:"type:text;not null"`
	StreamToken  string     `gorm:"type:text;not null"` // generated from streaming server
	StreamKey    string     `gorm:"type:text;not null"` // generated from web
	StreamType   StreamType `gorm:"type:varchar(50);not null"`
	ThumbnailURL string     `gorm:"type:text;not null"`
	StartedAt    time.Time  `gorm:"default:CURRENT_TIMESTAMP"`
	EndedAt      time.Time  `gorm:"default:CURRENT_TIMESTAMP"`
	User         User       `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

type Notification struct {
	NotificationID uint      `gorm:"primaryKey;autoIncrement"`
	UserID         uint      `gorm:"not null"`
	StreamID       uint      `gorm:"not null"`
	Content        string    `gorm:"type:text;not null"`
	CreatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	Stream         Stream    `gorm:"foreignKey:StreamID;constraint:OnDelete:CASCADE"`
	User           User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

// type Chat struct {
// 	MessageID uint      `gorm:"primaryKey;autoIncrement"`
// 	StreamID  uint      `gorm:"not null"`
// 	UserID    uint      `gorm:"not null"`
// 	Message   string    `gorm:"type:text;not null"`
// 	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
// 	User      User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
// 	Stream    Stream    `gorm:"foreignKey:StreamID;constraint:OnDelete:CASCADE"`
// }

type Subscription struct {
	SubscriptionID uint      `gorm:"primaryKey;autoIncrement"`
	SubscriberID   uint      `gorm:"not null"`
	StreamerID     uint      `gorm:"not null"`
	CreatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	// StartDate      time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	// EndDate        time.Time `gorm:"not null"`
	// AutoRenew      bool      `gorm:"not null"`
	Subscriber User `gorm:"foreignKey:SubscriberID;constraint:OnDelete:CASCADE"`
	Streamer   User `gorm:"foreignKey:StreamerID;constraint:OnDelete:CASCADE"`
}

type StreamAnalytics struct {
	AnalyticsID uint      `gorm:"primaryKey;autoIncrement"`
	StreamID    uint      `gorm:"not null"`
	Views       uint      `gorm:"not null"`
	Likes       uint      `gorm:"not null"`
	Comments    uint      `gorm:"not null"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	Stream      Stream    `gorm:"foreignKey:StreamID;constraint:OnDelete:CASCADE"`
}

type Like struct {
	LikesID   uint          `gorm:"primaryKey;autoIncrement"`
	UserID    uint          `gorm:"not null"`
	StreamID  uint          `gorm:"not null"`
	LikeEmote LikeEmoteType `gorm:"type:varchar(50);not null"`
	CreatedAt time.Time     `gorm:"default:CURRENT_TIMESTAMP"`
	Stream    Stream        `gorm:"foreignKey:StreamID;constraint:OnDelete:CASCADE"`
	User      User          `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

type Comment struct {
	CommentsID uint      `gorm:"primaryKey;autoIncrement"`
	UserID     uint      `gorm:"not null"`
	StreamID   uint      `gorm:"not null"`
	Comment    string    `gorm:"type:text;not null"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	Stream     Stream    `gorm:"foreignKey:StreamID;constraint:OnDelete:CASCADE"`
	User       User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}
type Share struct {
	ShareID   uint      `gorm:"primaryKey;autoIncrement"`
	UserID    uint      `gorm:"not null"`
	StreamID  uint      `gorm:"not null"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	Stream    Stream    `gorm:"foreignKey:StreamID;constraint:OnDelete:CASCADE"`
	User      User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}
