package gormModels

import "gorm.io/gorm"

type Banner struct {
	gorm.Model
	ID        uint64
	FeatureID uint64
	IsActive  bool
}

type BannerContent struct {
	BannerID uint64
	Title    string
	Text     string
	URL      string
	Banner   Banner `gorm:"foreignKey:BannerID;constraint:OnDelete:CASCADE;"`
}

type BannerTagRelation struct {
	BannerID uint64 `gorm:"uniqueIndex:idx_banner_tag"`
	TagID    uint64 `gorm:"uniqueIndex:idx_banner_tag"`
	Banner   Banner `gorm:"foreignKey:BannerID;constraint:OnDelete:CASCADE;"`
}
