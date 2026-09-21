# CRUD Dashboard and Storefront for Chinook

A small CRUD demo in go: a w2ui admin dashboard and an htmx storefront
over the [SQLite sample database](https://www.sqlitetutorial.net/sqlite-sample-database/).

## Stack

- **Backend**: Go with Echo v5, templ, [w2go](https://github.com/dv1x3r/w2go), SQLite with goose migrations;
- **Admin**: w2ui, served as plain ES modules;
- **Storefront**: tailwind, htmx and alpine, bundled with bun;

## Routes

- `/` - the storefront, tracks of the sample catalog;
- `/admin` - the w2ui dashboard;
- `/api/v1` - the admin API the dashboard talks to;

## Layout

```
cmd/server             - entry point: config, database, routes
data/sql/base.sql      - the sample Chinook database, applied to an empty file
data/sql/updates       - goose migrations applied on top of it
internal/api           - admin API: routes and w2go services
internal/storefront    - storefront: handler, services, templ templates
internal/lib           - database store, logger, helpers
web/admin              - admin dashboard widgets
web/client             - storefront sources, bundled into web/client/dist
```

## Prerequisites

- **Go >= 1.27**
- **bun**, only to rebuild the storefront bundle

## Build

```sh
cp .env.example .env  # default settings
make run              # build and start the server
```

Note that `web/` is embedded into the binary, so the server has to be rebuilt
for a change in the admin widgets or in the storefront bundle to show up.

## Make targets

- `make build` / `make run` - compile or run the server
- `make templ` - regenerate the templ templates
- `make web` - rebuild the storefront bundle into `web/client/dist`
- `make db-status` / `make db-up` / `make db-create` - goose migrations in `data/sql/updates`
