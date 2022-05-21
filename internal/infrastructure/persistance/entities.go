package persistance

import (
	"time"
)

type AlertDAO struct {
	ChatId        int64  `gorm:"primaryKey;autoIncrement:false"`
	UserId        int64  `gorm:"primaryKey;autoIncrement:false"`
	CoinName      string `gorm:"primaryKey"`
	IsGreaterThan bool   `gorm:"primaryKey"`
	Price         float64
	CreatedAt     time.Time
}

func (s *AlertDAO) TableName() string {
	return "alerts"
}
