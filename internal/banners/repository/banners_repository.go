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

func (db *Postgres) CreateBanner(
	banner gormModels.Banner,
	bannerInfo gormModels.BannerContent,
	bannerTag []gormModels.BannerTagRelation,
) (uint64, error) {
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var recBanner gormModels.Banner
	if err := db.DB.Create(&banner).Scan(&recBanner).Error; err != nil {
		return 0, err
	}

	if err := db.DB.Create(bannerInfo).Error; err != nil {
		return 0, err
	}

	for _, tag := range bannerTag {
		if err := db.DB.Create(tag).Error; err != nil {
			return 0, err
		}
	}

	return recBanner.ID, tx.Commit().Error
}

func (db *Postgres) UpdateBanner(
	banner gormModels.Banner,
	bannerInfo gormModels.BannerContent,
	bannerTag []gormModels.BannerTagRelation,
) error {
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := db.DB.Model(&gormModels.Banner{ID: banner.ID}).
		Updates(banner).
		Error; err != nil {
		return err
	}

	if err := db.DB.Model(&gormModels.BannerContent{BannerID: banner.ID}).
		Updates(bannerInfo).
		Error; err != nil {
		return err
	}

	// TODO: Incorrect updating. Rewrite!!!
	for _, tag := range bannerTag {
		if err := db.DB.Model(&gormModels.BannerTagRelation{BannerID: banner.ID}).
			Updates(tag).
			Error; err != nil {
			return err
		}
	}

	return tx.Commit().Error
}

func (db *Postgres) DeleteBanner(bannerID uint64) error {
	return db.DB.Unscoped().Delete(&gormModels.Banner{}, "id = ?", bannerID).Error
}

func (db *Postgres) GetUserBanner(userID, tagID, featureID uint64) (gormModels.Banner, error) {
	return gormModels.Banner{}, nil
}

func (db *Postgres) GetBanners(
	tagID, featureID, limit, offset uint64,
) ([]gormModels.Banner, error) {
	return nil, nil
}
