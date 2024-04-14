package httpModels

import "time"

type BannerID struct {
	ID uint64 `json:"banner_id"`
}

type BannerContent struct {
	Title string `json:"title"`
	Text  string `json:"text"`
	URL   string `json:"url"`
}

type Banner struct {
	BannerID  uint64        `json:"banner_id"`
	TagsIDs   []uint64      `json:"tag_ids"`
	FeatureID uint64        `json:"feature_id"`
	Content   BannerContent `json:"content"`
	IsActive  bool          `json:"is_active"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}
