package storefront

import (
	"net/http"
	"strconv"

	"github.com/dv1x3r/chinook-go/internal/lib/echotempl"
	"github.com/dv1x3r/chinook-go/internal/storefront/services"
	"github.com/dv1x3r/chinook-go/internal/storefront/templates/page/catalog"

	"github.com/labstack/echo/v5"
)

type Handler struct {
	svc *services.Service
}

func NewHandler(svc *services.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(g *echo.Group) {
	g.GET("/", h.GetCatalog)
	g.GET("/playlist/:playlistID", h.GetCatalog)
}

func (h *Handler) GetCatalog(c *echo.Context) error {
	ctx := c.Request().Context()

	playlistID, _ := strconv.Atoi(c.Param("playlistID"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	props := catalog.Props{
		Path:       c.Request().URL.Path,
		PlaylistID: playlistID,
		Request: services.TrackRequest{
			Search:     c.QueryParam("search"),
			SearchBy:   c.QueryParam("search-by"),
			Filters:    c.QueryParam("filters"),
			PlaylistID: playlistID,
			Offset:     offset,
		},
	}

	var err error
	if props.Tracks, err = h.svc.GetTracks(ctx, props.Request); err != nil {
		return err
	}

	// the infinite scroll asks for the next rows only, htmx puts them after the last one
	if offset > 0 {
		return echotempl.Render(c, http.StatusOK, catalog.Rows(props))
	}

	if props.Playlists, err = h.svc.GetPlaylists(ctx); err != nil {
		return err
	}

	if props.AllTracks, err = h.svc.CountTracks(ctx); err != nil {
		return err
	}

	if props.Genres, err = h.svc.GetGenres(ctx); err != nil {
		return err
	}

	if props.MediaTypes, err = h.svc.GetMediaTypes(ctx); err != nil {
		return err
	}

	props.Title, props.Total = "All tracks", props.AllTracks
	if playlistID > 0 {
		playlist, ok := findPlaylist(props.Playlists, playlistID)
		if !ok {
			return echo.ErrNotFound
		}
		props.Title, props.Total = playlist.Name, playlist.Tracks
	}

	return echotempl.Render(c, http.StatusOK, catalog.List(props))
}

func findPlaylist(playlists []services.Playlist, id int) (services.Playlist, bool) {
	for _, playlist := range playlists {
		if playlist.ID == id {
			return playlist, true
		}
	}
	return services.Playlist{}, false
}
