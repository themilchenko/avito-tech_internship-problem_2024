package bannersUsecase

import (
	"errors"
	"slices"

	"github.com/themilchenko/avito_internship-problem_2024/internal/domain"
	gormModels "github.com/themilchenko/avito_internship-problem_2024/internal/models/gorm"
	httpModels "github.com/themilchenko/avito_internship-problem_2024/internal/models/http"
	"gorm.io/gorm"
)

type bannersUsecase struct {
	bannersRepository domain.BannersRepository
	bannersCache      domain.BannersCache
}

func NewBannersUsecase(b domain.BannersRepository, c domain.BannersCache) bannersUsecase {
	return bannersUsecase{
		bannersRepository: b,
		bannersCache:      c,
	}
}

func (u bannersUsecase) GetUserBanner(
	tagID, featureID uint64,
	useLastVersion bool,
) (httpModels.BannerContent, error) {
	if !useLastVersion {
		bannerContent, err := u.bannersCache.Get(tagID, featureID)
		if err != nil && err.Error() != domain.ErrRedisNotFound.Error() {
			return httpModels.BannerContent{}, err
		}
		if err == nil {
			return bannerContent, nil
		}
	}

	// Go to postgres to get actual data
	bannerContent, err := u.bannersRepository.GetUserBanner(tagID, featureID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return httpModels.BannerContent{}, domain.ErrNotFound
		}
		return httpModels.BannerContent{}, domain.ErrInternal
	}

	// Check existion in cache
	tags, err := u.bannersRepository.GetBannerTags(bannerContent.BannerID)
	if err != nil {
		return httpModels.BannerContent{}, err
	}
	noCacheData := true
	for _, t := range tags {
		_, err = u.bannersCache.Get(t.TagID, featureID)
		if err != nil && err.Error() != domain.ErrRedisNotFound.Error() {
			return httpModels.BannerContent{}, nil
		}
		if err == nil {
			noCacheData = false
			break
		}
	}

	// Only if in redis no data about this banner we push it
	if noCacheData {
		for _, t := range tags {
			if err := u.bannersCache.Add(t.TagID, featureID, bannerContent.ToHTTPModel()); err != nil {
				return httpModels.BannerContent{}, err
			}
		}
	}

	return bannerContent.ToHTTPModel(), nil
}

func (u bannersUsecase) GetBanners(
	tagID, featureID, limit, offset uint64,
) ([]httpModels.Banner, error) {
	banners, err := u.bannersRepository.GetBanners(tagID, featureID, limit, offset)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []httpModels.Banner{}, domain.ErrNotFound
		}
		return []httpModels.Banner{}, err
	}

	httpBanners := make([]httpModels.Banner, len(banners))
	for i, b := range banners {
		httpBanners[i] = b.ToHTTPModel()
		relations, err := u.bannersRepository.GetBannerTags(httpBanners[i].BannerID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return []httpModels.Banner{}, domain.ErrNotFound
			}
			return nil, err
		}

		httpBanners[i].TagsIDs = make([]uint64, len(relations))
		for j, r := range relations {
			httpBanners[i].TagsIDs[j] = r.TagID
		}
	}

	return httpBanners, nil
}

func (u bannersUsecase) CreateBanner(banner httpModels.Banner) (uint64, error) {
	bannerTags := make([]gormModels.BannerTagRelation, len(banner.TagsIDs))
	for i, t := range banner.TagsIDs {
		_, err := u.bannersRepository.GetUserBanner(t, banner.FeatureID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, err
		}
		if err == nil {
			return 0, domain.ErrBannerAlreadyExist
		}
		bannerTags[i] = gormModels.BannerTagRelation{
			BannerID: banner.BannerID,
			TagID:    t,
		}
	}

	bannerID, err := u.bannersRepository.CreateBanner(gormModels.Banner{
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
	if _, err := u.bannersRepository.GetBannerByID(banner.BannerID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrNotFound
		}
		return err
	}

	tags, err := u.bannersRepository.GetBannerTags(banner.BannerID)
	if err != nil {
		return err
	}

	tagIDs := make([]uint64, len(tags))
	for i, t := range tags {
		tagIDs[i] = t.TagID
	}

	tagsToAdd := make([]uint64, 0)
	for _, t := range banner.TagsIDs {
		var i int
		if i = slices.Index(tagIDs, t); i == -1 {
			tagsToAdd = append(tagsToAdd, t)
			continue
		}

		tagIDs = append(tagIDs[:i], tagIDs[i+1:]...)
	}

	if err := u.bannersRepository.UpdateBannerTransactional(
		gormModels.Banner{
			ID:        banner.BannerID,
			FeatureID: banner.FeatureID,
			IsActive:  banner.IsActive,
		}, gormModels.BannerContent{
			BannerID: banner.BannerID,
			Title:    banner.Content.Title,
			Text:     banner.Content.Text,
			URL:      banner.Content.URL,
		},
		tagIDs,
		tagsToAdd,
	); err != nil {
		return err
	}

	return nil
}

func (u bannersUsecase) DeleteBannerByID(bannerID uint64) error {
	if _, err := u.bannersRepository.GetBannerByID(bannerID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrNotFound
		}
		return err
	}
	return u.bannersRepository.DeleteBanner(bannerID)
}

func (u bannersUsecase) GetBannerByFeatureAndTag(
	featureID, tagID uint64,
) (httpModels.Banner, error) {
	banner, err := u.bannersRepository.GetBannerByFeatureAndTag(featureID, tagID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return httpModels.Banner{}, domain.ErrNotFound
		}
		return httpModels.Banner{}, err
	}
	return banner.ToHTTPModel(), nil
}
