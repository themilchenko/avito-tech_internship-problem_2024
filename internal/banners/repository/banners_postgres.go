package bannersRepository

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

	limit, offset uint64
}

func NewPostgres(url string, l, o uint64) (*Postgres, error) {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  logger.Warn,
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
		&gormModels.Banner{},
		&gormModels.BannerContent{},
		&gormModels.BannerTagRelation{},
	)

	return &Postgres{
		DB:     db,
		limit:  l,
		offset: o,
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

	bannerInfo.BannerID = recBanner.ID
	if err := tx.Create(&bannerInfo).Error; err != nil {
		return 0, err
	}

	for _, tag := range bannerTag {
		tag.BannerID = recBanner.ID
		if err := tx.Create(&tag).Error; err != nil {
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
		if err := tx.Unscoped().Delete(&gormModels.BannerTagRelation{}, "tag_id = ?", t).Error; err != nil {
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

func (db *Postgres) GetBannerTags(bannerID uint64) ([]gormModels.BannerTagRelation, error) {
	recievedTags := make([]gormModels.BannerTagRelation, 0)
	if err := db.DB.Model(&gormModels.BannerTagRelation{}).
		Where("banner_id = ?", bannerID).
		Scan(&recievedTags).
		Error; err != nil {
		return []gormModels.BannerTagRelation{}, err
	}
	return recievedTags, nil
}

func (db *Postgres) GetUserBanner(tagID, featureID uint64) (gormModels.BannerContent, error) {
	var recievedUsrBanner gormModels.BannerContent
	if err := db.DB.Model(&gormModels.BannerContent{}).Distinct().
		InnerJoins("JOIN banners ON banners.id = banner_contents.banner_id").
		InnerJoins("JOIN banner_tag_relations ON banners.id = banner_tag_relations.banner_id").
		Select("banner_contents.banner_id, banner_contents.title, banner_contents.text, banner_contents.url").
		Where("banners.feature_id = ? AND banner_tag_relations.tag_id = ?", featureID, tagID).
		First(&recievedUsrBanner).Error; err != nil {
		return gormModels.BannerContent{}, err
	}
	return recievedUsrBanner, nil
}

func (db *Postgres) GetBanners(
	tagID, featureID, limit, offset uint64,
) ([]gormModels.BannerJoined, error) {
	whereCondition := ""
	var args []interface{}
	if tagID != 0 {
		whereCondition += "banner_tag_relations.tag_id = ?"
		args = append(args, tagID)
	}
	if featureID != 0 {
		if len(whereCondition) != 0 {
			whereCondition += " AND banners.feature_id = ?"
		} else {
			whereCondition += "banners.feature_id = ?"
		}
		args = append(args, featureID)
	}

	if limit == 0 {
		limit = db.limit
	}
	if offset == 0 {
		offset = db.offset
	}

	var recievedBanners []gormModels.BannerJoined
	if err := db.DB.Model(&gormModels.Banner{}).Distinct().
		InnerJoins("JOIN banner_contents ON banners.id = banner_contents.banner_id").
		InnerJoins("JOIN banner_tag_relations ON banners.id = banner_tag_relations.banner_id").
		Select("banners.id, banners.feature_id, banners.is_active, banners.created_at, banners.updated_at,"+
			" banner_contents.title, banner_contents.text, banner_contents.url").
		Where(whereCondition, args...).
		Limit(int(limit)).
		Offset(int(offset - 1)).
		Find(&recievedBanners).
		Error; err != nil {
		return []gormModels.BannerJoined{}, err
	}
	return recievedBanners, nil
}

func (db *Postgres) GetBannerByID(bannerID uint64) (gormModels.Banner, error) {
	var resBanner gormModels.Banner
	if err := db.DB.Model(&gormModels.Banner{}).Where("id = ?", bannerID).First(&resBanner).Error; err != nil {
		return gormModels.Banner{}, err
	}
	return resBanner, nil
}

func (db *Postgres) GetBannerByFeatureAndTag(featureID, tagID uint64) (gormModels.Banner, error) {
	var recievedBanner gormModels.Banner
	if err := db.DB.Model(&gormModels.Banner{}).Distinct().
		InnerJoins("JOIN banner_tag_relations ON banners.id = banner_tag_relations.banner_id").
		Select("banners.*").
		Where("banners.feature_id = ? AND banner_tag_relations.tag_id = ?", featureID, tagID).
		First(&recievedBanner).Error; err != nil {
		return gormModels.Banner{}, err
	}
	return recievedBanner, nil
}
