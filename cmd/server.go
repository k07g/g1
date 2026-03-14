package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/k07g/g1/internal/infrastructure/repository"
	"github.com/k07g/g1/internal/interface/handler"
	itemUsecase "github.com/k07g/g1/internal/usecase/item"
	userUsecase "github.com/k07g/g1/internal/usecase/user"
)

func main() {
	// ---- 依存性の注入（外側から内側へ） ----
	//
	//  [infrastructure]  →  [usecase]  →  [interface/handler]
	//   Repo                 Interactor      Handler
	//       ↓                    ↓
	//  domain.Repository   usecase.UseCase  （インターフェース経由）

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL が設定されていません")
	}

	itemRepo, err := repository.NewPostgreSQLRepository(dsn)
	if err != nil {
		log.Fatalf("PostgreSQL接続エラー(items): %v", err)
	}
	userRepo, err := repository.NewPostgreSQLUserRepository(dsn)
	if err != nil {
		log.Fatalf("PostgreSQL接続エラー(users): %v", err)
	}
	log.Println("🐘 PostgreSQLリポジトリを使用")

	itemHandler := handler.NewItemHandler(itemUsecase.NewInteractor(itemRepo))
	userHandler := handler.NewUserHandler(userUsecase.NewInteractor(userRepo))

	// ---- ルーティング ----
	mux := http.NewServeMux()
	mux.Handle("/api/items", itemHandler)
	mux.Handle("/api/items/", itemHandler)
	mux.Handle("/api/users", userHandler)
	mux.Handle("/api/users/", userHandler)
	mux.HandleFunc("/health", healthHandler)

	// ---- サーバー起動 ----
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      loggingMiddleware(mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("🚀 サーバー起動: http://localhost:8080")
	log.Println("📋 エンドポイント一覧:")
	log.Println("   GET    /api/items       - アイテム一覧取得")
	log.Println("   POST   /api/items       - アイテム新規作成")
	log.Println("   GET    /api/items/{id}  - アイテム1件取得")
	log.Println("   PUT    /api/items/{id}  - アイテム更新")
	log.Println("   DELETE /api/items/{id}  - アイテム削除")
	log.Println("   GET    /api/users       - ユーザー一覧取得")
	log.Println("   POST   /api/users       - ユーザー新規作成")
	log.Println("   GET    /api/users/{id}  - ユーザー1件取得")
	log.Println("   PUT    /api/users/{id}  - ユーザー更新")
	log.Println("   DELETE /api/users/{id}  - ユーザー削除")
	log.Println("   GET    /health          - ヘルスチェック")

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("サーバーエラー: %v", err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("[%s] %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}
