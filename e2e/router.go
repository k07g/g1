package e2e

import (
	"net/http"

	"github.com/k07g/g1/internal/interface/handler"
)

func newServeMux(itemHandler *handler.ItemHandler, userHandler *handler.UserHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/api/items", itemHandler)
	mux.Handle("/api/items/", itemHandler)
	mux.Handle("/api/users", userHandler)
	mux.Handle("/api/users/", userHandler)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})
	return mux
}
