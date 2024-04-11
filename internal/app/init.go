package app

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	httpAuth "github.com/themilchenko/avito_internship-problem_2024/internal/auth/delivery"
	authMiddleware "github.com/themilchenko/avito_internship-problem_2024/internal/auth/delivery/middleware"
	authRepository "github.com/themilchenko/avito_internship-problem_2024/internal/auth/repository"
	authUsecase "github.com/themilchenko/avito_internship-problem_2024/internal/auth/usecase"
	httpBanners "github.com/themilchenko/avito_internship-problem_2024/internal/banners/delivery"
	bannersRepository "github.com/themilchenko/avito_internship-problem_2024/internal/banners/repository"
	bannersUsecase "github.com/themilchenko/avito_internship-problem_2024/internal/banners/usecase"
	"github.com/themilchenko/avito_internship-problem_2024/internal/config"
	"github.com/themilchenko/avito_internship-problem_2024/internal/domain"
	"github.com/themilchenko/avito_internship-problem_2024/internal/utils/crypto"
	"github.com/themilchenko/avito_internship-problem_2024/pkg/logger"
)

type server struct {
	server *echo.Echo
	config *config.Config

	authUsecase    domain.AuthUsecase
	bannersUsecase domain.BannersUsecase

	authHandler    httpAuth.AuthHandler
	bannersHandler httpBanners.BannersHandler

	authMiddleware *authMiddleware.AuthMiddleware
}

func NewServer(e *echo.Echo, c *config.Config) *server {
	return &server{
		server: e,
		config: c,
	}
}

func (s *server) Start() error {
	s.makeUsecases()
	s.makeMiddlewares()
	s.makeHandlers()
	s.makeRouter()
	s.makeEchoLogger()
	return s.server.Start(
		s.config.Server.Address + ":" + strconv.FormatUint(s.config.Server.Port, 10),
	)
}

func (s *server) GetEchoInstance() *echo.Echo {
	return s.server
}

func (s *server) makeHandlers() {
	s.authHandler = httpAuth.NewAuthHandler(s.authUsecase, s.config.CookieSettings)
	s.bannersHandler = httpBanners.NewBannersHandler(s.bannersUsecase)
}

func (s *server) makeUsecases() {
	pgParams := s.config.FormatDbAddr()

	authDB, err := authRepository.NewPostgres(pgParams)
	if err != nil {
		s.server.Logger.Fatal(err)
	}
	bannersDB, err := bannersRepository.NewPostgres(pgParams)
	if err != nil {
		s.server.Logger.Fatal(err)
	}

	s.authUsecase = authUsecase.NewAuthUsecase(authDB, s.config.CookieSettings, crypto.HashPassword)
	s.bannersUsecase = bannersUsecase.NewBannersUsecase(bannersDB, authDB)
}

func (s *server) makeRouter() {
	v1 := s.server.Group("/api")
	v1.Use(logger.Middleware())
	v1.Use(middleware.Secure())

	v1.POST("/login", s.authHandler.Login)
	v1.POST("/signup", s.authHandler.SignUp)
	v1.DELETE("/logout", s.authHandler.Logout, s.authMiddleware.LoginRequired)
	v1.GET("/auth", s.authHandler.Auth)

	v1.GET("/user_banner", s.bannersHandler.GetUserBanner, s.authMiddleware.LoginRequired)

	b := v1.Group("/banner", s.authMiddleware.AdminRequiured)
	b.GET("", s.bannersHandler.GetBanners)
	b.POST("", s.bannersHandler.CreateBanner)

	bid := b.Group("/:id")
	bid.PATCH("", s.bannersHandler.PatchBanner)
	bid.DELETE("", s.bannersHandler.DeleteBanner)
}

func (s *server) makeMiddlewares() {
	s.authMiddleware = authMiddleware.NewAuthMiddleware(s.authUsecase)
}

func (s *server) makeEchoLogger() {
	s.server.Logger = logger.GetInstance()
	s.server.Logger.SetLevel(logger.ToLevel(s.config.LoggerLvl))
	s.server.Logger.Info("server started")
}
