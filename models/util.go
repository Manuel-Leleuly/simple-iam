package models

import (
	"time"
)

type Name struct {
	FirstName string `gorm:"column:first_name;not null" json:"first_name"`
	LastName  string `gorm:"column:last_name" json:"last_name"`
}

type TimeRecord struct {
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime;not null;<-create" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime;not null;autoUpdateTime" json:"updated_at"`
	DeletedAt time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}
