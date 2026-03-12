# g1

## ディレクトリ構成

```
clean-api/
├── cmd/api/
│   └── main.go                         # エントリポイント・DI組み立て
├── internal/
│   ├── domain/item/
│   │   ├── entity.go                   # エンティティ・ドメインエラー
│   │   └── repository.go               # リポジトリインターフェース
│   ├── usecase/item/
│   │   ├── port.go                     # Input/Output DTO・UseCaseインターフェース
│   │   └── interactor.go               # ビジネスロジック実装
│   ├── interface/handler/
│   │   └── item_handler.go             # HTTPハンドラー
│   └── infrastructure/repository/
│       └── in_memory.go                # リポジトリ実装（インメモリ）
└── go.mod
```

## 依存の方向

```
handler → usecase.UseCase（IF）← interactor
interactor → domain.Repository（IF）← in_memory
```

## テスト戦略

```
domain          → 依存ゼロ。エンティティの振る舞いを直接テスト
infrastructure  → 実装を直接テスト（DBの代わりにインメモリ）
usecase         → mockRepository を注入してDBに依存せずテスト
handler         → mockUseCase を注入して net/http/httptest でテスト
```