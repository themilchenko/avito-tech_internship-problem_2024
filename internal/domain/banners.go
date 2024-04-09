package domain

import gormModels "github.com/themilchenko/avito_internship-problem_2024/internal/models/gorm"

type BannersUsecase interface{}

type BannersRepository interface {
	CreateBanner(
		banner gormModels.Banner,
		bannerInfo gormModels.BannerContent,
		bannerTag []gormModels.BannerTagRelation,
	) (uint64, error)
	UpdateBanner(
		banner gormModels.Banner,
		bannerInfo gormModels.BannerContent,
		bannerTag []gormModels.BannerTagRelation,
	) error
	DeleteBanner(bannerID uint64) error
	GetUserBanner(userID, tagID, featureID uint64) (gormModels.Banner, error)
	GetBanners(tagID, featureID, limit, offset uint64) ([]gormModels.Banner, error)
}
