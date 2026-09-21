package routes

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/dv1x3r/chinook-go/internal/api/services"

	"github.com/dv1x3r/w2go/w2"
	"github.com/labstack/echo/v5"
)

func (h *Handler) GetPlaylistTrackGrid(c *echo.Context) error {
	playlistID, err := strconv.Atoi(c.Param("playlistID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, w2.NewErrorResponse(err.Error()))
	}
	req, err := w2.ParseGetGridRequest(c.QueryParam("request"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, w2.NewErrorResponse(err.Error()))
	}
	res, err := h.svc.GetPlaylistTrackGrid(c.Request().Context(), req, playlistID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, w2.NewErrorResponse(err.Error()))
	}
	return c.JSON(http.StatusOK, res)
}

func (h *Handler) PostPlaylistTrackAdd(c *echo.Context) error {
	playlistID, err := strconv.Atoi(c.Param("playlistID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, w2.NewErrorResponse(err.Error()))
	}
	req, err := w2.ParseSaveFormRequest[services.PlaylistTrackForm](c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, w2.NewErrorResponse(err.Error()))
	}
	err = h.svc.AddPlaylistTrack(c.Request().Context(), req, playlistID)
	if errors.Is(err, services.ErrPlaylistTrackExists) {
		return c.JSON(http.StatusConflict, w2.NewErrorResponse(err.Error()))
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, w2.NewErrorResponse(err.Error()))
	}
	return c.JSON(http.StatusOK, w2.NewSuccessResponse())
}

func (h *Handler) PostPlaylistTrackDelete(c *echo.Context) error {
	playlistID, err := strconv.Atoi(c.Param("playlistID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, w2.NewErrorResponse(err.Error()))
	}
	req, err := w2.ParseRemoveGridRequest(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, w2.NewErrorResponse(err.Error()))
	}
	err = h.svc.RemovePlaylistTracks(c.Request().Context(), req, playlistID)
	if errors.Is(err, services.ErrPlaylistTrackNotFound) {
		return c.JSON(http.StatusNotFound, w2.NewErrorResponse(err.Error()))
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, w2.NewErrorResponse(err.Error()))
	}
	return c.JSON(http.StatusOK, w2.NewSuccessResponse())
}
