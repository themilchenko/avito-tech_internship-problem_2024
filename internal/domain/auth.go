package domain

import (
	gormModels "github.com/themilchenko/avito_internship-problem_2024/internal/models/gorm"
	httpModels "github.com/themilchenko/avito_internship-problem_2024/internal/models/http"
)

type AuthUsecase interface {
	SignUp(user httpModels.User) (string, uint64, error)
	Login(user httpModels.User) (string, uint64, error)
	Logout(sessionID string) error
	Auth(sessionID string) (uint64, error)
	GetUserBySessionID(sessionID string) (httpModels.User, error)
}

type AuthRepository interface {
	CreateUser(user gormModels.User) (uint64, error)
	CreateSession(session gormModels.Session) (string, error)
	GetUserBySessionID(sessionID string) (gormModels.User, error)
	DeleteBySessionID(sessionID string) error
	GetUserByUsername(username string) (gormModels.User, error)
}
