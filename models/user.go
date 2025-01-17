package models

import (
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	Id         string     `gorm:"primary_key;column:id;not null;<-create" json:"id"`
	Name       Name       `gorm:"embedded"`
	Username   string     `gorm:"column:username;not null" json:"username"`
	Email      string     `gorm:"column:email;not null" json:"email"`
	Password   string     `gorm:"column:password;not null" json:"-"`
	TimeRecord TimeRecord `gorm:"embedded"`
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate(db *gorm.DB) error {
	u.Id = strings.ReplaceAll(uuid.New().String(), "-", "")
	return nil
}
