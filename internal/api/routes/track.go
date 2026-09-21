package routes

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/dv1x3r/chinook-go/internal/api/services"

	"github.com/dv1x3r/w2go/w2"
	"github.com/labstack/echo/v5"
)

func (h *Handler) GetTrackGrid(c *echo.Context) error {
	artistID, err := strconv.Atoi(c.Param("artistID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, w2.NewErrorResponse(err.Error()))
	}
	albumID, err := strconv.Atoi(c.Param("albumID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, w2.NewErrorResponse(err.Error()))
	}
	req, err := w2.ParseGetGridRequest(c.QueryParam("request"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, w2.NewErrorResponse(err.Error()))
	}
	res, err := h.svc.GetTrackGrid(c.Request().Context(), req, artistID, albumID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, w2.NewErrorResponse(err.Error()))
	}
	return c.JSON(http.StatusOK, res)
}

func (h *Handler) GetTrackDropdown(c *echo.Context) error {
	req, err := w2.ParseGetDropdownRequest(c.QueryParam("request"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, w2.NewErrorResponse(err.Error()))
	}
	res, err := h.svc.GetTrackDropdown(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, w2.NewErrorResponse(err.Error()))
	}
	return c.JSON(http.StatusOK, res)
}

func (h *Handler) PostTrackGrid(c *echo.Context) error {
	req, err := w2.ParseSaveGridRequest[services.Track](c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, w2.NewErrorResponse(err.Error()))
	}
	err = h.svc.SaveTrackGrid(c.Request().Context(), req)
	if errors.Is(err, services.ErrTrackNotFound) {
		return c.JSON(http.StatusNotFound, w2.NewErrorResponse(err.Error()))
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, w2.NewErrorResponse(err.Error()))
	}
	return c.JSON(http.StatusOK, w2.NewSuccessResponse())
}

func (h *Handler) PostTrackForm(c *echo.Context) error {
	req, err := w2.ParseSaveFormRequest[services.Track](c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, w2.NewErrorResponse(err.Error()))
	}
	if req.Record.ID == 0 {
		req.Record.ID, err = h.svc.CreateTrack(c.Request().Context(), req)
	} else {
		err = h.svc.UpdateTrack(c.Request().Context(), req)
	}
	if errors.Is(err, services.ErrTrackNotFound) {
		return c.JSON(http.StatusNotFound, w2.NewErrorResponse(err.Error()))
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, w2.NewErrorResponse(err.Error()))
	}
	return c.JSON(http.StatusOK, w2.NewSuccessResponse())
}

func (h *Handler) PostTrackDelete(c *echo.Context) error {
	req, err := w2.ParseRemoveGridRequest(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, w2.NewErrorResponse(err.Error()))
	}
	err = h.svc.DeleteTracks(c.Request().Context(), req)
	if errors.Is(err, services.ErrTrackHasRelated) {
		return c.JSON(http.StatusConflict, w2.NewErrorResponse(err.Error()))
	} else if errors.Is(err, services.ErrTrackNotFound) {
		return c.JSON(http.StatusNotFound, w2.NewErrorResponse(err.Error()))
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, w2.NewErrorResponse(err.Error()))
	}
	return c.JSON(http.StatusOK, w2.NewSuccessResponse())
}
