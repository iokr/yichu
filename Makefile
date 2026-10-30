APP     := yichu
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
BUILD   := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -s -w -X 'main.Version=$(VERSION)' -X 'main.BuildTime=$(BUILD)'

.PHONY: all dev front verify build run cross up down logs backup migrate clean

all: front verify build

front:
	cd frontend && npm ci && npm run build

verify:
	@test -f web/dist/index.html || (echo "❌ web/dist/index.html 缺失" && exit 1)

build:
	CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -o bin/$(APP) ./cmd/yichu

run: all
	./bin/$(APP)

dev:
	@echo "后端: go run ./cmd/yichu   前端: cd frontend && npm run dev"

cross: front verify
	CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build -trimpath -ldflags "$(LDFLAGS)" -o bin/$(APP)-linux-amd64     ./cmd/yichu
	CGO_ENABLED=0 GOOS=linux   GOARCH=arm64 go build -trimpath -ldflags "$(LDFLAGS)" -o bin/$(APP)-linux-arm64     ./cmd/yichu
	CGO_ENABLED=0 GOOS=darwin  GOARCH=arm64 go build -trimpath -ldflags "$(LDFLAGS)" -o bin/$(APP)-darwin-arm64    ./cmd/yichu
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -trimpath -ldflags "$(LDFLAGS)" -o bin/$(APP)-windows-amd64.exe ./cmd/yichu

up:
	docker compose up -d --build

down:
	docker compose down

logs:
	docker compose logs -f yichu

backup:
	docker compose exec -T mysql mysqldump --single-transaction -uyichu -p$$MYSQL_PASSWORD yichu | gzip > backup/yichu_$$(date +%Y%m%d_%H%M).sql.gz

migrate:
	docker compose exec -T mysql mysql --default-character-set=utf8mb4 -uyichu -p$$MYSQL_PASSWORD yichu < migrations/0001_init.up.sql
	docker compose exec -T mysql mysql --default-character-set=utf8mb4 -uyichu -p$$MYSQL_PASSWORD yichu < migrations/0002_seed.up.sql

clean:
	rm -rf bin && rm -rf web/dist/* && touch web/dist/.gitkeep
