<div align="center">

```
 ██████╗  ██╗
██╔════╝ ███║
██║  ███╗╚██║
██║   ██║ ██║
╚██████╔╝ ██║
 ╚═════╝  ╚═╝
```

**Clean Architecture · Go REST API**

[![Go](https://img.shields.io/badge/Go-1.26+-00ADD8?style=flat-square&logo=go&logoColor=white)](https://golang.org)
[![Architecture](https://img.shields.io/badge/Architecture-Clean-blueviolet?style=flat-square)](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
[![License](https://img.shields.io/badge/License-MIT-green?style=flat-square)](LICENSE)

</div>

---

## 概要

クリーンアーキテクチャの原則に従った Go 製 REST API のリファレンス実装。
**依存の方向を一方向に保つ**ことで、テスタブルかつ変更に強い設計を実現する。

---

## ディレクトリ構成

```
g1/
├── cmd/
│   └── server.go                    # エントリポイント・DI組み立て
│
├── internal/
│   ├── domain/item/
│   │   ├── entity.go              # エンティティ・ドメインエラー
│   │   └── repository.go          # リポジトリインターフェース
│   │
│   ├── usecase/item/
│   │   ├── port.go                # Input/Output DTO・UseCaseインターフェース
│   │   └── interactor.go          # ビジネスロジック実装
│   │
│   ├── interface/handler/
│   │   └── item_handler.go        # HTTPハンドラー
│   │
│   └── infrastructure/repository/
│       └── postgres.go            # リポジトリ実装（PostgreSQL）
│
├── e2e/                             # E2Eテスト
│   ├── items_test.go
│   ├── server_test.go
│   ├── helper_test.go
│   └── router.go
│
├── docker-compose.yml               # PostgreSQL コンテナ定義
├── Dockerfile                       # マルチステージビルド
└── go.mod
```

---

## アーキテクチャ

依存は**外側 → 内側**の一方向のみ。内側のレイヤーは外側を一切知らない。

```
┌─────────────────────────────────────────────────────┐
│                   interface/handler                  │
│                  HTTP リクエスト受付                  │
└───────────────────────┬─────────────────────────────┘
                        │ depends on
                        ▼
┌─────────────────────────────────────────────────────┐
│                  usecase (interface)                 │
│              ビジネスロジックの境界定義               │
└──────────┬────────────────────────┬─────────────────┘
           │ implements              │ depends on
           ▼                        ▼
┌──────────────────┐    ┌──────────────────────────────┐
│   interactor     │    │     domain (entity/repo IF)  │
│ ビジネスロジック  │───▶│      ドメインの核心          │
└──────────────────┘    └──────────────────────────────┘
                                    ▲
                                    │ implements
┌─────────────────────────────────────────────────────┐
│              infrastructure/repository               │
│                  インフラ詳細（PostgreSQL）           │
└─────────────────────────────────────────────────────┘
```

### 依存ルール（一行まとめ）

```
handler → usecase.UseCase（IF）← interactor → domain.Repository（IF）← postgres
```

---

## テスト戦略

各レイヤーを**独立してテスト**できる構造になっている。

| レイヤー | 戦略 | 外部依存 |
|---|---|---|
| `domain` | エンティティの振る舞いを直接テスト | **ゼロ** |
| `usecase` | `mockRepository` を注入 | DBに依存しない |
| `handler` | `mockUseCase` を注入、`httptest` で検証 | HTTPサーバー不要 |
| `e2e` | 全レイヤー結合、実PostgreSQLに接続 | PostgreSQL（必須） |

---

## Getting Started

### 前提

- [Docker](https://docs.docker.com/get-docker/) がインストールされていること

### サーバー起動

```bash
# リポジトリをクローン
git clone https://github.com/k07g/g1.git
cd g1

# PostgreSQL + アプリを起動
docker compose up
```

### テスト実行

```bash
# ユニットテスト
make test/unit

# E2Eテスト（PostgreSQL を自動起動）
make test/e2e
```

> `make test/e2e` は db コンテナが未起動の場合、自動的に `docker compose up -d db` を実行します。

---

<div align="center">

*Simplicity is the ultimate sophistication.*

</div>
