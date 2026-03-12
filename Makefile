.PHONY: test test/unit test/cover test/race lint fmt vet clean help

# ==============================
# 変数定義
# ==============================

GO       := go
PACKAGES := ./...
COVER_OUT := coverage.out
COVER_HTML := coverage.html

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
	$(GO) build -o bin/api ./cmd/api

## run: サーバーを起動
run:
	$(GO) run ./cmd/api

# ==============================
# クリーンアップ
# ==============================

## clean: 生成ファイルを削除
clean:
	rm -f $(COVER_OUT) $(COVER_HTML) bin/api
	@echo "🧹 クリーンアップ完了"

# ==============================
# ヘルプ
# ==============================

## help: 利用可能なコマンド一覧を表示
help:
	@echo "使い方: make [target]"
	@echo ""
	@grep -E '^## ' Makefile | sed 's/## //' | column -t -s ':'