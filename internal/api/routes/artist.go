package routes

import (
	"errors"
	"net/http"

	"github.com/dv1x3r/chinook-go/internal/api/services"

	"github.com/dv1x3r/w2go/w2"
	"github.com/labstack/echo/v5"
)

func (h *Handler) GetArtistGrid(c *echo.Context) error {
	req, err := w2.ParseGetGridRequest(c.QueryParam("request"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, w2.NewErrorResponse(err.Error()))
	}
	res, err := h.svc.GetArtistGrid(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, w2.NewErrorResponse(err.Error()))
	}
	return c.JSON(http.StatusOK, res)
}

func (h *Handler) GetArtistDropdown(c *echo.Context) error {
	req, err := w2.ParseGetDropdownRequest(c.QueryParam("request"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, w2.NewErrorResponse(err.Error()))
	}
	res, err := h.svc.GetArtistDropdown(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, w2.NewErrorResponse(err.Error()))
	}
	return c.JSON(http.StatusOK, res)
}

func (h *Handler) PostArtistForm(c *echo.Context) error {
	req, err := w2.ParseSaveFormRequest[services.Artist](c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, w2.NewErrorResponse(err.Error()))
	}
	if req.Record.ID == 0 {
		req.Record.ID, err = h.svc.CreateArtist(c.Request().Context(), req)
	} else {
		err = h.svc.UpdateArtist(c.Request().Context(), req)
	}
	if errors.Is(err, services.ErrArtistExists) {
		return c.JSON(http.StatusConflict, w2.NewErrorResponse(err.Error()))
	} else if errors.Is(err, services.ErrArtistNotFound) {
		return c.JSON(http.StatusNotFound, w2.NewErrorResponse(err.Error()))
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, w2.NewErrorResponse(err.Error()))
	}
	return c.JSON(http.StatusOK, w2.NewSuccessResponse())
}

func (h *Handler) PostArtistDelete(c *echo.Context) error {
	req, err := w2.ParseRemoveGridRequest(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, w2.NewErrorResponse(err.Error()))
	}
	err = h.svc.DeleteArtists(c.Request().Context(), req)
	if errors.Is(err, services.ErrArtistHasRelated) {
		return c.JSON(http.StatusConflict, w2.NewErrorResponse(err.Error()))
	} else if errors.Is(err, services.ErrArtistNotFound) {
		return c.JSON(http.StatusNotFound, w2.NewErrorResponse(err.Error()))
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, w2.NewErrorResponse(err.Error()))
	}
	return c.JSON(http.StatusOK, w2.NewSuccessResponse())
}
