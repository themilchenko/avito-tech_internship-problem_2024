package domain

import (
	gormModels "github.com/themilchenko/avito_internship-problem_2024/internal/models/gorm"
	httpModels "github.com/themilchenko/avito_internship-problem_2024/internal/models/http"
)

type BannersUsecase interface {
	GetUserBanner(tagID, featureID uint64) (httpModels.BannerContent, error)
	GetBanners(tagID, featureID, limit, offset uint64) ([]httpModels.Banner, error)
	CreateBanner(banner httpModels.Banner) (uint64, error)
	UpdateBannerByID(banner httpModels.Banner) error
	DeleteBannerByID(bannerID uint64) error
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
	GetBannerTags(bannerID uint64) ([]uint64, error)
	DeleteBanner(bannerID uint64) error
	GetUserBanner(tagID, featureID uint64) (gormModels.BannerContent, error)
	GetBanners(tagID, featureID, limit, offset uint64) ([]httpModels.Banner, error)
}
