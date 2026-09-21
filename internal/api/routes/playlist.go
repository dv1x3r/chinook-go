package routes

import (
	"errors"
	"net/http"

	"github.com/dv1x3r/chinook-go/internal/api/services"

	"github.com/dv1x3r/w2go/w2"
	"github.com/labstack/echo/v5"
)

func (h *Handler) GetPlaylistGrid(c *echo.Context) error {
	req, err := w2.ParseGetGridRequest(c.QueryParam("request"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, w2.NewErrorResponse(err.Error()))
	}
	res, err := h.svc.GetPlaylistGrid(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, w2.NewErrorResponse(err.Error()))
	}
	return c.JSON(http.StatusOK, res)
}

func (h *Handler) PostPlaylistForm(c *echo.Context) error {
	req, err := w2.ParseSaveFormRequest[services.Playlist](c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, w2.NewErrorResponse(err.Error()))
	}
	if req.Record.ID == 0 {
		req.Record.ID, err = h.svc.CreatePlaylist(c.Request().Context(), req)
	} else {
		err = h.svc.UpdatePlaylist(c.Request().Context(), req)
	}
	if errors.Is(err, services.ErrPlaylistNotFound) {
		return c.JSON(http.StatusNotFound, w2.NewErrorResponse(err.Error()))
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, w2.NewErrorResponse(err.Error()))
	}
	return c.JSON(http.StatusOK, w2.NewSuccessResponse())
}

func (h *Handler) PostPlaylistDelete(c *echo.Context) error {
	req, err := w2.ParseRemoveGridRequest(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, w2.NewErrorResponse(err.Error()))
	}
	err = h.svc.DeletePlaylists(c.Request().Context(), req)
	if errors.Is(err, services.ErrPlaylistHasRelated) {
		return c.JSON(http.StatusConflict, w2.NewErrorResponse(err.Error()))
	} else if errors.Is(err, services.ErrPlaylistNotFound) {
		return c.JSON(http.StatusNotFound, w2.NewErrorResponse(err.Error()))
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, w2.NewErrorResponse(err.Error()))
	}
	return c.JSON(http.StatusOK, w2.NewSuccessResponse())
}
