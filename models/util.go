package models

import (
	"time"
)

// structs that will be stored in DB
type Name struct {
	FirstName string `gorm:"column:first_name;not null" json:"first_name"`
	LastName  string `gorm:"column:last_name" json:"last_name"`
}

type TimeRecord struct {
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime;not null;<-create" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime;not null;autoUpdateTime" json:"updated_at"`
}

// structs that will NOT be stored in DB
type NameRequest struct {
	FirstName string `json:"first_name" validate:"required,name"`
	LastName  string `json:"last_name" validate:"omitempty,name"`
}
