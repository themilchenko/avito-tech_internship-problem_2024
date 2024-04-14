package httpBanners

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/themilchenko/avito_internship-problem_2024/internal/domain"
	httpModels "github.com/themilchenko/avito_internship-problem_2024/internal/models/http"
)

type BannersHandler struct {
	bannersUsecase domain.BannersUsecase
	authUsecase    domain.AuthUsecase
}

func NewBannersHandler(b domain.BannersUsecase, a domain.AuthUsecase) BannersHandler {
	return BannersHandler{
		bannersUsecase: b,
		authUsecase:    a,
	}
}

func (h BannersHandler) GetUserBanner(c echo.Context) error {
	tagID, err := strconv.ParseUint(c.QueryParam("tag_id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	featureID, err := strconv.ParseUint(c.QueryParam("feature_id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	useLastVersion, err := strconv.ParseBool(c.QueryParam("use_last_version"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	usrBanner, err := h.bannersUsecase.GetUserBanner(tagID, featureID, useLastVersion)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, usrBanner)
}

func (h BannersHandler) GetBanners(c echo.Context) error {
	tagStr := c.QueryParam("tag_id")
	var tagID uint64 = 0
	var err error
	if len(tagStr) != 0 {
		tagID, err = strconv.ParseUint(tagStr, 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
	}

	featureStr := c.QueryParam("feature_id")
	var featureID uint64 = 0
	if len(featureStr) != 0 {
		featureID, err = strconv.ParseUint(featureStr, 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
	}

	limitStr := c.QueryParam("feature_id")
	var limit uint64 = 0
	if len(limitStr) != 0 {
		limit, err = strconv.ParseUint(limitStr, 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
	}

	offsetStr := c.QueryParam("offset")
	var offset uint64 = 0
	if len(offsetStr) != 0 {
		offset, err = strconv.ParseUint(offsetStr, 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
	}

	banners, err := h.bannersUsecase.GetBanners(tagID, featureID, limit, offset)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, banners)
}

func (h BannersHandler) CreateBanner(c echo.Context) error {
	var recievedBanner httpModels.Banner
	if err := c.Bind(&recievedBanner); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	bannerID, err := h.bannersUsecase.CreateBanner(recievedBanner)
	if err != nil {
		if errors.Is(err, domain.ErrBannerAlreadyExist) {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, httpModels.BannerID{ID: bannerID})
}

func (h BannersHandler) PatchBanner(c echo.Context) error {
	bannerID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	var recievedBanner httpModels.Banner
	if err := c.Bind(&recievedBanner); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	recievedBanner.BannerID = bannerID

	if err := h.bannersUsecase.UpdateBannerByID(recievedBanner); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, httpModels.EmptyStruct{})
}

func (h BannersHandler) DeleteBanner(c echo.Context) error {
	bannerID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err := h.bannersUsecase.DeleteBannerByID(bannerID); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusNoContent)
}
