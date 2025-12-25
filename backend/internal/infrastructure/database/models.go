package database

import (
	"gorm.io/gorm"
)

type WriterModel struct {
	ID        uint64 `gorm:"primaryKey"`
	Name      string `gorm:"type:varchar(255);not null"`
	BirthYear int    `gorm:"not null"`
	DeathYear *int
	Bio       *string `gorm:"type:text"`
}

func (WriterModel) TableName() string {
	return "writers"
}

type WorkModel struct {
	ID       uint64 `gorm:"primaryKey"`
	Title    string `gorm:"type:varchar(255);not null"`
	AuthorID uint64 `gorm:"not null;index"`
}

func (WorkModel) TableName() string {
	return "works"
}

type OpinionModel struct {
	WriterID      uint64  `gorm:"primaryKey"`
	WorkID        uint64  `gorm:"primaryKey"`
	Sentiment     bool    `gorm:"not null"`
	Quote         string  `gorm:"type:text;not null"`
	Source        string  `gorm:"type:varchar(255);not null"`
	Page          *string `gorm:"type:varchar(100)"`
	StatementYear *int
}

func (OpinionModel) TableName() string {
	return "opinions"
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&WriterModel{},
		&WorkModel{},
		&OpinionModel{},
	)
}
