package httpBanners

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/themilchenko/avito_internship-problem_2024/internal/domain"
	httpModels "github.com/themilchenko/avito_internship-problem_2024/internal/models/http"
)

type BannersHandler struct {
	bannersUsecase domain.BannersUsecase
}

func NewBannersHandler(b domain.BannersUsecase) BannersHandler {
	return BannersHandler{
		bannersUsecase: b,
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
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, usrBanner)
}

func (h BannersHandler) GetBanners(c echo.Context) error {
	tagID, err := strconv.ParseUint(c.QueryParam("tag_id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	featureID, err := strconv.ParseUint(c.QueryParam("feature_id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	limit, err := strconv.ParseUint(c.QueryParam("limit"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	offset, err := strconv.ParseUint(c.QueryParam("offset"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	banners, err := h.bannersUsecase.GetBanners(tagID, featureID, limit, offset)
	if err != nil {
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
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusNoContent, httpModels.EmptyStruct{})
}
