# コントリビューションガイド

このプロジェクトへのコントリビューションを歓迎します。Issue の報告・機能提案・Pull Request など、あらゆる形での参加を感謝します。

---

## 行動規範

参加にあたっては [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) を遵守してください。

---

## バグ報告・機能提案

[GitHub Issues](https://github.com/k07g/g1/issues) からお知らせください。

- **バグ報告**: 再現手順・期待する動作・実際の動作・環境情報を記載してください。
- **機能提案**: 背景・解決したい課題・提案する解決策を記載してください。

---

## Pull Request

### 前提

- [Docker](https://docs.docker.com/get-docker/) がインストールされていること
- Go 1.26 以上がインストールされていること

### 手順

```bash
# 1. フォーク & クローン
git clone https://github.com/<your-name>/g1.git
cd g1

# 2. ブランチを作成
git checkout -b feature/your-feature-name

# 3. 変更を加える

# 4. フォーマット & 静的解析
make fmt
make vet

# 5. ユニットテストを実行
make test/unit

# 6. E2Eテストを実行（PostgreSQL を自動起動）
make test/e2e

# 7. コミット & プッシュ
git commit -m "feat: your feature description"
git push origin feature/your-feature-name
```

その後、`main` ブランチへの Pull Request を作成してください。

### コミットメッセージ

[Conventional Commits](https://www.conventionalcommits.org/) に準拠してください。

| プレフィックス | 用途 |
|---|---|
| `feat` | 新機能 |
| `fix` | バグ修正 |
| `refactor` | リファクタリング |
| `test` | テストの追加・修正 |
| `docs` | ドキュメントの変更 |
| `chore` | ビルド・CI などの変更 |

### PR のチェックリスト

- [ ] `make fmt` でフォーマット済み
- [ ] `make vet` が通ること
- [ ] `make test/unit` が通ること
- [ ] `make test/e2e` が通ること
- [ ] 新機能・バグ修正にはテストを追加している

---

## アーキテクチャ

変更を加える前に [README.md](README.md) のアーキテクチャセクションを確認してください。各レイヤーの依存方向を守ることが重要です。

```
handler → usecase（IF）← interactor → domain.Repository（IF）← postgres
```

新しいエンティティを追加する場合は、以下の順で実装してください。

1. `internal/domain/<entity>/` — エンティティ・リポジトリインターフェース
2. `internal/usecase/<entity>/` — Input/Output DTO・UseCase インターフェース・Interactor
3. `internal/interface/handler/` — HTTP ハンドラー
4. `internal/infrastructure/repository/` — PostgreSQL 実装
5. `cmd/server.go` — DI の組み立てとルーティング登録
6. `e2e/` — E2E テスト
