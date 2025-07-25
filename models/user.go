package models

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

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
	// this is to prevent changing the id of test user
	if u.Id == "" {
		u.Id = strings.ReplaceAll(uuid.New().String(), "-", "")
	}
	return nil
}

// only for request body. NOT saved to DB
type UserRequest struct {
	NameRequest
	Username string `json:"username" validate:"required,username"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"password"`
}

type UserUpdateRequest struct {
	FirstName string `json:"first_name" validate:"omitempty,name"`
	LastName  string `json:"last_name" validate:"omitempty,name"`
	Username  string `json:"username" validate:"omitempty,username"`
}

// response
type UserResponse struct {
	Id string `json:"id"`
	Name
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
