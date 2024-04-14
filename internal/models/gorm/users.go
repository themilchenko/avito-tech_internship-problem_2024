package gormModels

import (
	"time"

	httpModels "github.com/themilchenko/avito_internship-problem_2024/internal/models/http"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       uint64
	Username string `gorm:"unique"`
	Password string
	Role     string
}

type Session struct {
	Value      string `gorm:"primarykey"`
	UserID     uint64
	ExpireDate time.Time
	User       User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
}

func (u User) ToHTTPModel() httpModels.User {
	return httpModels.User{
		Username: u.Username,
		Password: u.Password,
		Role:     u.Role,
	}
}
