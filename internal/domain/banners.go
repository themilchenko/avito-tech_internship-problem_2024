package domain

import (
	gormModels "github.com/themilchenko/avito_internship-problem_2024/internal/models/gorm"
	httpModels "github.com/themilchenko/avito_internship-problem_2024/internal/models/http"
)

type BannersUsecase interface {
	GetUserBanner(tagID, featureID uint64, useLastVersion bool) (httpModels.BannerContent, error)
	GetBanners(tagID, featureID, limit, offset uint64) ([]httpModels.Banner, error)
	CreateBanner(banner httpModels.Banner) (uint64, error)
	UpdateBannerByID(banner httpModels.Banner) error
	DeleteBannerByID(bannerID uint64) error
	GetBannerByFeatureAndTag(featureID, tagID uint64) (httpModels.Banner, error)
}

type BannersRepository interface {
	CreateBanner(
		banner gormModels.Banner,
		bannerInfo gormModels.BannerContent,
		bannerTag []gormModels.BannerTagRelation,
	) (uint64, error)
	UpdateBanner(banner gormModels.Banner) error
	UpdateBannerInfo(bannerInfo gormModels.BannerContent) error
	UpdateBannerTransactional(
		banner gormModels.Banner,
		BannerContent gormModels.BannerContent,
		tagsToDel, tagsToAdd []uint64,
	) error
	DeleteTagByID(tagID uint64) error
	CreateTag(tag gormModels.BannerTagRelation) error
	GetBannerTags(bannerID uint64) ([]gormModels.BannerTagRelation, error)
	DeleteBanner(bannerID uint64) error
	GetUserBanner(tagID, featureID uint64) (gormModels.BannerContent, error)
	GetBanners(tagID, featureID, limit, offset uint64) ([]gormModels.BannerJoined, error)
	GetBannerByID(bannerID uint64) (gormModels.Banner, error)
	GetBannerByFeatureAndTag(featureID, tagID uint64) (gormModels.Banner, error)
}

type BannersCache interface {
	Get(tagID, featureID uint64) (httpModels.BannerContent, error)
	Add(tagID, featureID uint64, bannerContent httpModels.BannerContent) error
	Delete(tagID, featureID uint64) error
	DeleteByFeatureID(featureID uint64) error
}
