package bannersRepository

import (
	gormModels "github.com/themilchenko/avito_internship-problem_2024/internal/models/gorm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Postgres struct {
	DB *gorm.DB
}

func NewPostgres(url string) (*Postgres, error) {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(
		&gormModels.Banner{},
		&gormModels.BannerContent{},
		&gormModels.BannerTagRelation{},
	)

	return &Postgres{
		DB: db,
	}, nil
}
