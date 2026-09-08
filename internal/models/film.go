package models

type Film struct {
	ID        uint   `gorm:"primaryKey"`
	TMDBID    int    `gorm:"uniqueIndex;not null"`
	Titre     string `gorm:"not null"`
	Annee     int
	PosterURL string
	Synopsis  string
}
