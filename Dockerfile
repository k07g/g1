# ==============================
# Stage 1: ビルド
# ==============================
FROM golang:1.26.1-alpine AS builder

# セキュリティアップデート + ビルドに必要なパッケージ
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# 依存関係のキャッシュを先に解決（ソースより先にコピー）
COPY go.mod go.sum* ./
RUN go mod download

# ソースコードをコピーしてビルド
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /app/bin/server \
    ./cmd

# ==============================
# Stage 2: 実行（最小イメージ）
# ==============================
FROM scratch

# ルート証明書・タイムゾーンデータをコピー
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# バイナリのみコピー
COPY --from=builder /app/bin/server /server

# 非rootユーザーで実行（scratch では直接 UID 指定）
USER 65534:65534

EXPOSE 8080

ENTRYPOINT ["/server"]