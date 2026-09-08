package models

import "time"

type Watchlist struct {
	UserID  uint `gorm:"primaryKey"`
	FilmID  uint `gorm:"primaryKey"`
	AddedAt time.Time

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Film Film `gorm:"foreignKey:FilmID;constraint:OnDelete:CASCADE"`
}
