package httpBanners

import "github.com/themilchenko/avito_internship-problem_2024/internal/domain"

type BannersHandler struct {
	bannersUsecase domain.BannersUsecase
}

func NewBannersHandler(b domain.BannersUsecase) BannersHandler {
	return BannersHandler{
		bannersUsecase: b,
	}
}
