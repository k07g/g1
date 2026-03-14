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

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go&logoColor=white)](https://golang.org)
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
│       └── in_memory.go           # リポジトリ実装（インメモリ）
│
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
│                  インフラ詳細（DB等）                 │
└─────────────────────────────────────────────────────┘
```

### 依存ルール（一行まとめ）

```
handler → usecase.UseCase（IF）← interactor → domain.Repository（IF）← in_memory
```

---

## テスト戦略

各レイヤーを**独立してテスト**できる構造になっている。

| レイヤー | 戦略 | 外部依存 |
|---|---|---|
| `domain` | エンティティの振る舞いを直接テスト | **ゼロ** |
| `infrastructure` | インメモリ実装を直接テスト | なし（DBの代替） |
| `usecase` | `mockRepository` を注入 | DBに依存しない |
| `handler` | `mockUseCase` を注入、`httptest` で検証 | HTTPサーバー不要 |

---

## Getting Started

```bash
# リポジトリをクローン
git clone https://github.com/k07g/g1.git
cd g1

# 依存関係の取得
go mod tidy

# サーバー起動
go run ./cmd

# テスト実行
go test ./...
```

---

<div align="center">

*Simplicity is the ultimate sophistication.*

</div>