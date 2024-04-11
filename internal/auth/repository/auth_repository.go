package authRepository

import (
	"log"
	"os"
	"time"

	gormModels "github.com/themilchenko/avito_internship-problem_2024/internal/models/gorm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Postgres struct {
	DB *gorm.DB
}

func NewPostgres(url string) (*Postgres, error) {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  logger.Info,
				IgnoreRecordNotFoundError: false,
				ParameterizedQueries:      true,
				Colorful:                  false,
			},
		),
	})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(
		&gormModels.User{},
		&gormModels.Session{},
	)

	return &Postgres{
		DB: db,
	}, nil
}

func (db *Postgres) CreateUser(user gormModels.User) (uint64, error) {
	var recUser gormModels.User
	if err := db.DB.Create(&user).Scan(&recUser).Error; err != nil {
		return 0, err
	}
	return recUser.ID, nil
}

func (db *Postgres) CreateSession(session gormModels.Session) (string, error) {
	if err := db.DB.Create(&session).Error; err != nil {
		return "", err
	}
	return session.Value, nil
}

func (db *Postgres) GetUserBySessionID(sessionID string) (gormModels.User, error) {
	var recievedUser gormModels.User
	if err := db.DB.
		Joins("JOIN sessions ON users.id = sessions.user_id").
		Where("sessions.value = ?", sessionID).
		Select("users.id, users.username, users.password, users.role").
		First(&recievedUser).Scan(&recievedUser).Error; err != nil {
		return gormModels.User{}, err
	}
	return recievedUser, nil
}

func (db *Postgres) DeleteBySessionID(sessionID string) error {
	return db.DB.Unscoped().Delete(&gormModels.Session{}, "value = ?", sessionID).Error
}

func (db *Postgres) GetUserByUsername(username string) (gormModels.User, error) {
	var recievedUser gormModels.User
	if err := db.DB.Model(&gormModels.User{
		Username: username,
	}).Scan(&recievedUser).Error; err != nil {
		return gormModels.User{}, err
	}
	return recievedUser, nil
}
