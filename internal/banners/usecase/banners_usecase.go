package bannersUsecase

import (
	"errors"
	"go/token"
	"slices"

	"github.com/themilchenko/avito_internship-problem_2024/internal/domain"
	gormModels "github.com/themilchenko/avito_internship-problem_2024/internal/models/gorm"
	httpModels "github.com/themilchenko/avito_internship-problem_2024/internal/models/http"
	"gorm.io/gorm"
)

type bannersUsecase struct {
	bannersRepository domain.BannersRepository
	authRepository    domain.AuthRepository
}

func NewBannersUsecase(b domain.BannersRepository, a domain.AuthRepository) bannersUsecase {
	return bannersUsecase{
		bannersRepository: b,
		authRepository:    a,
	}
}

func (u bannersUsecase) GetUserBanner(
	tagID, featureID uint64,
	useLastVersion bool,
) (httpModels.BannerContent, error) {
	bannerContent, err := u.bannersRepository.GetUserBanner(tagID, featureID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return httpModels.BannerContent{}, domain.ErrNotFound
		}
		return httpModels.BannerContent{}, domain.ErrInternal
	}
	return bannerContent.ToHTTPModel(), nil
}

func (u bannersUsecase) GetBanners(
	tagID, featureID, limit, offset uint64,
) ([]httpModels.Banner, error) {
	banners, err := u.bannersRepository.GetBanners(tagID, featureID, limit, offset)
	if err != nil {
		return []httpModels.Banner{}, err
	}

	for _, b := range banners {
		b.TagsIDs, err = u.bannersRepository.GetBannerTags(b.BannerID)
		if err != nil {
			return nil, err
		}
	}
	return banners, nil
}

func (u bannersUsecase) CreateBanner(banner httpModels.Banner) (uint64, error) {
	bannerTags := make([]gormModels.BannerTagRelation, len(banner.TagsIDs))
	for i, t := range banner.TagsIDs {
		bannerTags[i] = gormModels.BannerTagRelation{
			BannerID: banner.BannerID,
			TagID:    t,
		}
	}

	bannerID, err := u.bannersRepository.CreateBanner(gormModels.Banner{
		ID:        banner.BannerID,
		FeatureID: banner.FeatureID,
	}, gormModels.BannerContent{
		Title: banner.Content.Title,
		Text:  banner.Content.Text,
		URL:   banner.Content.URL,
	}, bannerTags)
	if err != nil {
		return 0, err
	}

	return bannerID, nil
}

func (u bannersUsecase) UpdateBannerByID(banner httpModels.Banner) error {
	tags, err := u.bannersRepository.GetBannerTags(banner.BannerID)
	if err != nil {
		return err
	}

	tagsToAdd := make([]uint64, 0)
	for _, t := range banner.TagsIDs {
		var i int
		if i = slices.Index(tags, t); i == -1 {
			tagsToAdd = append(tagsToAdd, t)
			continue
		}

		tags = append(tags[:i], tags[i+1:]...)
	}

	if err := u.bannersRepository.UpdateBannerTransactional(gormModels.Banner{
		ID:        banner.BannerID,
		FeatureID: banner.FeatureID,
	}, gormModels.BannerContent{
		Title: banner.Content.Title,
		Text:  banner.Content.Text,
		URL:   banner.Content.URL,
	}, tags, tagsToAdd); err != nil {
		return err
	}

	return nil
}

func (u bannersUsecase) DeleteBannerByID(bannerID uint64) error {
	return u.DeleteBannerByID(bannerID)
}
