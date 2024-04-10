package bannersRepository

import (
	gormModels "github.com/themilchenko/avito_internship-problem_2024/internal/models/gorm"
	httpModels "github.com/themilchenko/avito_internship-problem_2024/internal/models/http"
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
	if err := tx.Create(&banner).Scan(&recBanner).Error; err != nil {
		return 0, err
	}

	if err := tx.Create(bannerInfo).Error; err != nil {
		return 0, err
	}

	for _, tag := range bannerTag {
		if err := tx.Create(tag).Error; err != nil {
			return 0, err
		}
	}

	return recBanner.ID, tx.Commit().Error
}

func (db *Postgres) UpdateBanner(banner gormModels.Banner) error {
	return db.DB.Model(&gormModels.Banner{ID: banner.ID}).Updates(banner).Error
}

func (db *Postgres) UpdateBannerInfo(bannerInfo gormModels.BannerContent) error {
	return db.DB.Model(&gormModels.BannerContent{BannerID: bannerInfo.BannerID}).
		Updates(bannerInfo).
		Error
}

func (db *Postgres) UpdateBannerTransactional(
	banner gormModels.Banner,
	BannerContent gormModels.BannerContent,
	tagsToDel, tagsToAdd []uint64,
) error {
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Model(&gormModels.Banner{ID: banner.ID}).Updates(banner).Error; err != nil {
		return err
	}

	if err := tx.Model(&gormModels.BannerContent{BannerID: BannerContent.BannerID}).
		Updates(BannerContent).
		Error; err != nil {
		return err
	}

	for _, t := range tagsToDel {
		if err := tx.Unscoped().Delete(&gormModels.BannerTagRelation{}).Where("tag_id = ?", t).Error; err != nil {
			return err
		}
	}

	for _, t := range tagsToAdd {
		if err := tx.Create(&gormModels.BannerTagRelation{BannerID: banner.ID, TagID: t}).Error; err != nil {
			return err
		}
	}

	return tx.Commit().Error
}

func (db *Postgres) DeleteTagByID(tagID uint64) error {
	return db.DB.Unscoped().Delete(&gormModels.BannerTagRelation{}).Where("tag_id = ?", tagID).Error
}

func (db *Postgres) CreateTag(tag gormModels.BannerTagRelation) error {
	return db.DB.Create(&tag).Error
}

func (db *Postgres) DeleteBanner(bannerID uint64) error {
	return db.DB.Unscoped().Delete(&gormModels.Banner{}, "id = ?", bannerID).Error
}

func (db *Postgres) GetBannerTags(bannerID uint64) ([]uint64, error) {
	var recievedTags []uint64
	if err := db.DB.Model(&gormModels.BannerTagRelation{}).
		Where("banner_id = ?", bannerID).
		Scan(&recievedTags).
		Error; err != nil {
		return []uint64{}, err
	}
	return recievedTags, nil
}

func (db *Postgres) GetUserBanner(tagID, featureID uint64) (gormModels.BannerContent, error) {
	var recievedUsrBanner gormModels.BannerContent
	if err := db.DB.Model(&gormModels.Banner{}).
		InnerJoins("JOIN banner_contents ON banners.id = banner_contents.banner_id").
		InnerJoins("JOIN banner_tag_relations ON banners.id = banner_tag_relations.banner_id").
		Select("banner.id, banner_contents.title, banner_contents.text, banner_contents.url").
		Where("banners.featureID = ? AND banner_tag_relations.tag_id = ?", featureID, tagID).
		First(&recievedUsrBanner).Error; err != nil {
		return gormModels.BannerContent{}, err
	}
	return recievedUsrBanner, nil
}

func (db *Postgres) GetBanners(
	tagID, featureID, limit, offset uint64,
) ([]httpModels.Banner, error) {
	var recievedBanners []httpModels.Banner
	if err := db.DB.Model(&gormModels.Banner{}).
		InnerJoins("JOIN banner_contents ON banners.id = banner_contents.banner_id").
		InnerJoins("JOIN banner_tag_relations ON banners.id = banner_tag_relations.banner_id").
		Select("banner.id, banner.feature_id, banner.is_active, banner_contents.title, banner_contents.text, banner_contents.url").
		Where("banners.featureID = ? AND banner_tag_relations.tag_id = ?", featureID, tagID).
		Find(&recievedBanners).
		Limit(int(limit)).
		Offset(int(offset)).
		Error; err != nil {
		return []httpModels.Banner{}, err
	}
	return recievedBanners, nil
}
