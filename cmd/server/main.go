package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dv1x3r/chinook-go/data"
	"github.com/dv1x3r/chinook-go/web"

	"github.com/dv1x3r/chinook-go/internal/api"
	apiroutes "github.com/dv1x3r/chinook-go/internal/api/routes"
	apiservices "github.com/dv1x3r/chinook-go/internal/api/services"

	"github.com/dv1x3r/chinook-go/internal/storefront"
	storefrontservices "github.com/dv1x3r/chinook-go/internal/storefront/services"

	"github.com/dv1x3r/chinook-go/internal/config"
	"github.com/dv1x3r/chinook-go/internal/lib/db"
	"github.com/dv1x3r/chinook-go/internal/lib/logger"

	"github.com/dv1x3r/w2go/w2"
	"github.com/dv1x3r/w2go/w2db"
	"github.com/dv1x3r/w2go/w2lib"
	"github.com/huandu/go-sqlbuilder"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func main() {
	cfg := config.MustLoad()

	// ── Logger ──────────────────────────────────────────────────────────────────
	if cfg.JsonLogger {
		logger.Create(logger.JsonHandler)
	} else {
		logger.Create(logger.TextHandler)
	}
	if cfg.Debug {
		logger.SetLevel(slog.LevelDebug)
	}

	// ── Database store ──────────────────────────────────────────────────────────
	store := db.NewSQLiteStore(logger.Get(), cfg.DBString)
	if err := store.Open(); err != nil {
		logger.Get().Error("unable to connect to the database", "err", err)
		os.Exit(1)
	}
	defer store.Close()

	if err := store.MigrateFile(data.FS, "sql/base.sql"); err != nil {
		logger.Get().Error("unable to initialize the database", "err", err)
		os.Exit(1)
	}

	if err := store.MigrateUp(data.FS, "sql/updates"); err != nil {
		logger.Get().Error("unable to migrate the database", "err", err)
		os.Exit(1)
	}

	w2db.SetFlavor(sqlbuilder.SQLite)

	// ── Server ──────────────────────────────────────────────────────────────────
	e := echo.NewWithConfig(echo.Config{
		Logger:      logger.Get(),
		IPExtractor: echo.ExtractIPFromXFFHeader(),
	})

	e.Use(
		middleware.RequestLogger(),
		middleware.Recover(),
		middleware.RequestID(),
		middleware.Secure(),
	)

	// ── Static assets ───────────────────────────────────────────────────────────
	// w2ui ships with w2go, the admin pages are plain modules loaded from /admin
	scripts := e.Group("", cacheControl(cacheRevalidate))
	scripts.StaticFS("/lib", echo.MustSubFS(w2lib.FS(), "."))
	scripts.StaticFS("/admin", echo.MustSubFS(web.FS, "admin"))

	// the storefront bundle carries a content hash, so it can be cached forever
	bundle := e.Group("", cacheControl(cacheImmutable))
	bundle.StaticFS("/"+web.CLIENT_JS, echo.MustSubFS(web.FS, "client/dist/"+web.CLIENT_JS))
	bundle.StaticFS("/"+web.CLIENT_CSS, echo.MustSubFS(web.FS, "client/dist/"+web.CLIENT_CSS))

	static := e.Group("", cacheControl(cacheStatic))
	static.StaticFS("/site.webmanifest", echo.MustSubFS(web.FS, "static/site.webmanifest"))

	// ── Admin API ───────────────────────────────────────────────────────────────
	apiSvc := apiservices.NewService(logger.Get(), store, cfg)
	apiHandler := apiroutes.NewHandler(apiSvc)
	api.Register(e, apiHandler, store, cfg.Debug, readOnlyMiddleware(cfg.ReadOnly))

	// ── Storefront ──────────────────────────────────────────────────────────────
	storefrontSvc := storefrontservices.NewService(logger.Get(), store, cfg)
	storefront.NewHandler(storefrontSvc).Register(e.Group("", cacheControl(cachePrivate)))

	startServer(e, cfg.ServerAddress)
}

func startServer(e *echo.Echo, address string) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	sc := echo.StartConfig{
		Address:         address,
		GracefulTimeout: 10 * time.Second,
		HideBanner:      true,
		HidePort:        true,
	}
	if err := sc.Start(ctx, e); err != nil {
		e.Logger.Error("failed to start server", "err", err)
	}
}

const (
	// the file name carries a content hash, so this URL can never go stale
	cacheImmutable = "public, max-age=31536000, immutable"
	// unversioned but seldom touched: favicons, bundled images
	cacheStatic = "public, max-age=86400"
	// embedded assets ship unversioned, so they are stored but always revalidated
	cacheRevalidate = "no-cache"
	// rendered pages differ per visitor: cart counter, hidden categories, payment methods
	cachePrivate = "private, no-cache"
)

func cacheControl(value string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			c.Response().Header().Set(echo.HeaderCacheControl, value)
			return next(c)
		}
	}
}

// the public demo runs without a login, so writes are simply turned off
func readOnlyMiddleware(readOnly bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			if readOnly && c.Request().Method != http.MethodGet {
				return c.JSON(http.StatusUnauthorized, w2.NewErrorResponse("Demo works in read-only mode"))
			}
			return next(c)
		}
	}
}
