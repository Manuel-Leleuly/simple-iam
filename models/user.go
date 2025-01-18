package models

import (
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRequest struct {
	Name
	Username string `gorm:"column:username;not null" json:"username"`
	Email    string `gorm:"column:email;not null" json:"email"`
	Password string `gorm:"column:password;not null" json:"password"`
}

type UserUpdateRequest struct {
	Name
	Username string `gorm:"column:username;not null" json:"username"`
}

type User struct {
	Id string `gorm:"primary_key;column:id;not null;<-create" json:"id"`
	Name
	Username string `gorm:"column:username;not null" json:"username"`
	Email    string `gorm:"column:email;not null" json:"email"`
	Password string `gorm:"column:password;not null" json:"-"`
	TimeRecord
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate(db *gorm.DB) error {
	u.Id = strings.ReplaceAll(uuid.New().String(), "-", "")
	return nil
}
