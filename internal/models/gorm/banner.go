package gormModels

import (
	httpModels "github.com/themilchenko/avito_internship-problem_2024/internal/models/http"
	"gorm.io/gorm"
)

type Banner struct {
	gorm.Model
	ID        uint64
	FeatureID uint64
	IsActive  bool
}

type BannerContent struct {
	BannerID uint64 `gorm:"primarykey"`
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

type BannerJoined struct {
	Banner
	BannerContent
}

func (b BannerContent) ToHTTPModel() httpModels.BannerContent {
	return httpModels.BannerContent{
		Title: b.Title,
		Text:  b.Text,
		URL:   b.URL,
	}
}

func (bj BannerJoined) ToHTTPModel() httpModels.Banner {
	return httpModels.Banner{
		BannerID:  bj.ID,
		FeatureID: bj.FeatureID,
		IsActive:  bj.IsActive,
		Content:   bj.BannerContent.ToHTTPModel(),
		CreatedAt: bj.CreatedAt,
		UpdatedAt: bj.UpdatedAt,
	}
}

func (b Banner) ToHTTPModel() httpModels.Banner {
	return httpModels.Banner{
		BannerID:  b.ID,
		FeatureID: b.FeatureID,
		IsActive:  b.IsActive,
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.CreatedAt,
	}
}
