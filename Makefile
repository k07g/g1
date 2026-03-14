.PHONY: test test/unit test/e2e test/all test/cover test/race lint fmt vet clean help \
        docker/build docker/run docker/stop docker/clean

# ==============================
# 変数定義
# ==============================

GO           := go
UNIT_PKGS    := $(shell go list ./... 2>/dev/null | grep -v '/e2e')
E2E_PKGS     := ./e2e/...
ALL_PKGS     := ./...
COVER_OUT    := coverage.out
COVER_HTML   := coverage.html

IMAGE_NAME     := g1
IMAGE_TAG      := latest
CONTAINER_NAME := g1
PORT           := 8080

# ==============================
# テスト
# ==============================

## test: ユニットテストを実行
test:
	$(GO) test $(UNIT_PKGS) -v

## test/unit: 並列・タイムアウト付きでユニットテストを実行
test/unit:
	$(GO) test $(UNIT_PKGS) -v -count=1 -timeout 30s

DATABASE_URL ?= postgres://g1:g1pass@localhost:5432/g1db?sslmode=disable

## test/e2e: E2Eテストを実行（dbコンテナが未起動なら自動起動）
test/e2e:
	@if ! docker compose ps db --status running 2>/dev/null | grep -q db; then \
		echo "▶ dbコンテナを起動します..."; \
		docker compose up -d db; \
		echo "⏳ PostgreSQL の起動を待機中..."; \
		docker compose exec db sh -c 'until pg_isready -U g1 -d g1db; do sleep 1; done'; \
	fi
	DATABASE_URL=$(DATABASE_URL) $(GO) test $(E2E_PKGS) -v -count=1 -timeout 60s

## test/all: ユニット + E2E テストをすべて実行
test/all:
	$(GO) test $(ALL_PKGS) -v -count=1 -timeout 60s

## test/race: レースコンディション検出つきでテストを実行
test/race:
	$(GO) test $(ALL_PKGS) -v -race -count=1

## test/cover: カバレッジ計測・HTML レポート生成（全テスト対象）
test/cover:
	$(GO) test $(ALL_PKGS) -coverprofile=$(COVER_OUT) -covermode=atomic
	$(GO) tool cover -func=$(COVER_OUT)
	$(GO) tool cover -html=$(COVER_OUT) -o $(COVER_HTML)
	@echo "📊 カバレッジレポート: $(COVER_HTML)"

## test/cover/show: カバレッジレポートをブラウザで開く
test/cover/show: test/cover
	open $(COVER_HTML)

# ==============================
# コード品質
# ==============================

## fmt: コードフォーマット
fmt:
	$(GO) fmt $(ALL_PKGS)

## vet: 静的解析
vet:
	$(GO) vet $(ALL_PKGS)

## lint: go vet + フォーマットチェック（gofmt）
lint: vet
	@echo "🔍 gofmt チェック..."
	@UNFORMATTED=$$(gofmt -l .); \
	if [ -n "$$UNFORMATTED" ]; then \
		echo "❌ フォーマットが必要なファイル:"; \
		echo "$$UNFORMATTED"; \
		exit 1; \
	else \
		echo "✅ フォーマット OK"; \
	fi

# ==============================
# CI 向け
# ==============================

## ci: CI環境向け（lint → race → cover）
ci: lint test/race test/cover

# ==============================
# ビルド
# ==============================

## build: バイナリをビルド
build:
	$(GO) build -o bin/server ./cmd

## run: サーバーを起動
run:
	$(GO) run ./cmd

# ==============================
# Docker
# ==============================

## docker/build: Dockerイメージをビルド
docker/build:
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) .

## docker/run: コンテナを起動（ポート8080）
docker/run:
	docker run --rm -d \
		--name $(CONTAINER_NAME) \
		-p $(PORT):8080 \
		$(IMAGE_NAME):$(IMAGE_TAG)
	@echo "🚀 起動: http://localhost:$(PORT)"

## docker/stop: コンテナを停止
docker/stop:
	docker stop $(CONTAINER_NAME)

## docker/clean: イメージ・コンテナを削除
docker/clean:
	-docker stop $(CONTAINER_NAME) 2>/dev/null
	-docker rm $(CONTAINER_NAME) 2>/dev/null
	-docker rmi $(IMAGE_NAME):$(IMAGE_TAG) 2>/dev/null
	@echo "🧹 Docker リソース削除完了"

# ==============================
# クリーンアップ
# ==============================

## clean: 生成ファイルをすべて削除
clean:
	rm -f $(COVER_OUT) $(COVER_HTML) bin/server
	@echo "🧹 クリーンアップ完了"

# ==============================
# ヘルプ
# ==============================

## help: 利用可能なコマンド一覧を表示
help:
	@echo "使い方: make [target]"
	@echo ""
	@grep -E '^## ' Makefile | sed 's/## //' | column -t -s ':'