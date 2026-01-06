package gota

import (
	"time"

	"gorm.io/gorm"
)

type Model struct {
	Id        string `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
