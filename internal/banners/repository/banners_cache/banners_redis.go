package bannersCache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	httpModels "github.com/themilchenko/avito_internship-problem_2024/internal/models/http"
)

const (
	minToDelete = 10
)

type BannersCacheRedis struct {
	cache *redis.Client
}

func NewBannersCacheRedis(addr, password string) (*BannersCacheRedis, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		return &BannersCacheRedis{}, err
	}

	return &BannersCacheRedis{
		cache: rdb,
	}, nil
}

func (c *BannersCacheRedis) Get(tagID, featureID uint64) (httpModels.BannerContent, error) {
	bannerJSON, err := c.cache.Get(context.Background(), fmt.Sprintf("%d:%d", featureID, tagID)).
		Result()
	if err != nil {
		return httpModels.BannerContent{}, err
	}

	var resContent httpModels.BannerContent
	if err := json.Unmarshal([]byte(bannerJSON), &resContent); err != nil {
		return httpModels.BannerContent{}, err
	}
	return resContent, nil
}

func (c *BannersCacheRedis) Add(
	tagID, featureID uint64,
	bannerContent httpModels.BannerContent,
) error {
	bannerJSON, err := json.Marshal(bannerContent)
	if err != nil {
		return err
	}

	// TODO: delete magic number
	return c.cache.Set(context.Background(), fmt.Sprintf("%d:%d", featureID, tagID), string(bannerJSON), minToDelete*time.Minute).
		Err()
}

func (c *BannersCacheRedis) Delete(tagID, featureID uint64) error {
	return c.cache.Del(context.Background(), fmt.Sprintf("%d:%d", featureID, tagID)).Err()
}

func (c *BannersCacheRedis) DeleteByFeatureID(featureID uint64) error {
	iter := c.cache.Scan(context.Background(), 0, fmt.Sprintf("%d:*", featureID), 0).Iterator()
	for iter.Next(context.Background()) {
		if err := c.cache.Del(context.Background(), iter.Val()).Err(); err != nil {
			return err
		}
	}
	if err := iter.Err(); err != nil {
		return err
	}
	return nil
}
