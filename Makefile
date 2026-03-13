.PHONY: test test/unit test/cover test/race lint fmt vet clean help

# ==============================
# 変数定義
# ==============================

GO       := go
PACKAGES := ./...
COVER_OUT := coverage.out
COVER_HTML := coverage.html

IMAGE_NAME     := g1
IMAGE_TAG      := latest
CONTAINER_NAME := g1
PORT           := 8080

# ==============================
# テスト
# ==============================

## test: すべてのテストを実行
test:
	$(GO) test $(PACKAGES) -v

## test/unit: 並列・タイムアウト付きでテストを実行
test/unit:
	$(GO) test $(PACKAGES) -v -count=1 -timeout 30s

## test/race: レースコンディション検出つきでテストを実行
test/race:
	$(GO) test $(PACKAGES) -v -race -count=1

## test/cover: カバレッジ計測・HTML レポート生成
test/cover:
	$(GO) test $(PACKAGES) -coverprofile=$(COVER_OUT) -covermode=atomic
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
	$(GO) fmt $(PACKAGES)

## vet: 静的解析
vet:
	$(GO) vet $(PACKAGES)

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
# CI 向け（lint + race + cover をまとめて実行）
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

## clean: 生成ファイルを削除
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