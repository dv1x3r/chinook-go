export CGO_ENABLED?=0
GOOSE=go tool goose -dir=./data/sql/updates sqlite ./data_db/lite.db

.PHONY: build
build:
	go build -o ./build/server ./cmd/server/main.go

.PHONY: run
run:
	go run ./cmd/server/main.go

.PHONY: rebuild
rebuild: templ web build

.PHONY: run-rebuild
run-rebuild: templ web run

.PHONY: test
test:
	go test -v ./...

.PHONY: test-race
test-race:
	CGO_ENABLED=1 go test -v -race ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: update
update:
	go get -u ./...
	go mod tidy

.PHONY: templ
templ:
	go tool templ generate --path ./internal

.PHONY: web
web:
	cd web/client && rm -rf dist && mkdir -p dist && \
	bun build ./src/client.js --outdir ./dist --entry-naming "[dir]/[name].[hash].[ext]" --minify && \
	BROWSERSLIST_IGNORE_OLD_DATA=1 bun run tailwindcss -i ./src/client.css -o ./dist/client.css --minify && \
	hash=$$(md5sum dist/client.css | cut -d' ' -f1 | cut -c1-8) && mv dist/client.css dist/client.$$hash.css

.PHONY: clean
clean:
	rm -rf ./build
	find . -name ".DS_Store" -type f -print -delete

.PHONY: db-up
db-up:
	$(GOOSE) up

.PHONY: db-up-to
db-up-to:
	@read -p "Up to version: " VALUE; \
	$(GOOSE) up-to $$VALUE

.PHONY: db-up-by-one
db-up-by-one:
	$(GOOSE) up-by-one

.PHONY: db-down
db-down:
	$(GOOSE) down

.PHONY: db-down-to
db-down-to:
	@read -p "Down to version: " VALUE; \
	$(GOOSE) down-to $$VALUE

.PHONY: db-status
db-status:
	$(GOOSE) status

.PHONY: db-reset
db-reset:
	$(GOOSE) reset

.PHONY: db-create
db-create:
	@read -p "Migration name: " VALUE; \
	$(GOOSE) create "$$VALUE" sql

