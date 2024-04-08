package bannersUsecase

import "github.com/themilchenko/avito_internship-problem_2024/internal/domain"

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
