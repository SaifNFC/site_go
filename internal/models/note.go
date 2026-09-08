package models

import "time"

type Note struct {
	UserID      uint `gorm:"primaryKey"`
	FilmID      uint `gorm:"primaryKey"`
	Note        int  `gorm:"not null"`
	Commentaire string
	CreatedAt   time.Time

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Film Film `gorm:"foreignKey:FilmID;constraint:OnDelete:CASCADE"`
}
