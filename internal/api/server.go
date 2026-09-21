package api

import (
	"github.com/dv1x3r/chinook-go/internal/api/routes"
	"github.com/dv1x3r/chinook-go/internal/lib/db"

	"github.com/dv1x3r/w2go/w2explorer"
	"github.com/labstack/echo/v5"
)

func Register(e *echo.Echo, h *routes.Handler, store db.Store, debug bool, middlewares ...echo.MiddlewareFunc) {

	// ── Admin dashboard ─────────────────────────────────────────────────────────
	e.GET("/admin", h.GetDashboard)

	v1 := e.Group("/api/v1", middlewares...)

	if debug {
		v1.GET("/sql", func(c *echo.Context) error {
			sqlSchemaHandler := w2explorer.SQLiteSchemaHTTPHandler(store.DB())
			sqlSchemaHandler(c.Response(), c.Request())
			return nil
		})
		v1.POST("/sql", func(c *echo.Context) error {
			sqlExecHandler := w2explorer.SQLExecHTTPHandler(store.DB())
			sqlExecHandler(c.Response(), c.Request())
			return nil
		})
	}

	// ── Artist ──────────────────────────────────────────────────────────────────
	v1.GET("/artist/grid", h.GetArtistGrid)
	v1.GET("/artist/dropdown", h.GetArtistDropdown)
	v1.POST("/artist/form", h.PostArtistForm)
	v1.POST("/artist/delete", h.PostArtistDelete)

	// ── Album ───────────────────────────────────────────────────────────────────
	v1.GET("/artist/:artistID/album/grid", h.GetAlbumGrid)
	v1.GET("/album/dropdown", h.GetAlbumDropdown)
	v1.POST("/album/form", h.PostAlbumForm)
	v1.POST("/album/delete", h.PostAlbumDelete)

	// ── Track ───────────────────────────────────────────────────────────────────
	v1.GET("/artist/:artistID/album/:albumID/track/grid", h.GetTrackGrid)
	v1.GET("/track/dropdown", h.GetTrackDropdown)
	v1.POST("/track/grid", h.PostTrackGrid)
	v1.POST("/track/form", h.PostTrackForm)
	v1.POST("/track/delete", h.PostTrackDelete)

	// ── Genre ───────────────────────────────────────────────────────────────────
	v1.GET("/genre/grid", h.GetGenreGrid)
	v1.GET("/genre/dropdown", h.GetGenreDropdown)
	v1.POST("/genre/form", h.PostGenreForm)
	v1.POST("/genre/delete", h.PostGenreDelete)

	// ── Media Type ──────────────────────────────────────────────────────────────
	v1.GET("/media_type/grid", h.GetMediaTypeGrid)
	v1.GET("/media_type/dropdown", h.GetMediaTypeDropdown)
	v1.POST("/media_type/form", h.PostMediaTypeForm)
	v1.POST("/media_type/delete", h.PostMediaTypeDelete)

	// ── Playlist ────────────────────────────────────────────────────────────────
	v1.GET("/playlist/grid", h.GetPlaylistGrid)
	v1.POST("/playlist/form", h.PostPlaylistForm)
	v1.POST("/playlist/delete", h.PostPlaylistDelete)

	// ── Playlist Track ──────────────────────────────────────────────────────────
	v1.GET("/playlist/:playlistID/track/grid", h.GetPlaylistTrackGrid)
	v1.POST("/playlist/:playlistID/track/add", h.PostPlaylistTrackAdd)
	v1.POST("/playlist/:playlistID/track/delete", h.PostPlaylistTrackDelete)

	// ── Customer ────────────────────────────────────────────────────────────────
	v1.GET("/customer/grid", h.GetCustomerGrid)
	v1.GET("/customer/dropdown", h.GetCustomerDropdown)
	v1.POST("/customer/form", h.PostCustomerForm)
	v1.POST("/customer/delete", h.PostCustomerDelete)

	// ── Employee ────────────────────────────────────────────────────────────────
	v1.GET("/employee/grid", h.GetEmployeeGrid)
	v1.GET("/employee/dropdown", h.GetEmployeeDropdown)
	v1.POST("/employee/form", h.PostEmployeeForm)
	v1.POST("/employee/delete", h.PostEmployeeDelete)

	// ── Invoice ─────────────────────────────────────────────────────────────────
	v1.GET("/invoice/grid", h.GetInvoiceGrid)
	v1.POST("/invoice/form", h.PostInvoiceForm)
	v1.POST("/invoice/delete", h.PostInvoiceDelete)

	// ── Invoice Item ────────────────────────────────────────────────────────────
	v1.GET("/invoice/:invoiceID/item/grid", h.GetInvoiceItemGrid)
	v1.POST("/invoice/:invoiceID/item/form", h.PostInvoiceItemForm)
	v1.POST("/invoice/:invoiceID/item/delete", h.PostInvoiceItemDelete)
}
