package e2e

import (
	"net/http"

	"github.com/k07g/g1/internal/interface/handler"
)

func newServeMux(h *handler.ItemHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/api/items", h)
	mux.Handle("/api/items/", h)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})
	return mux
}
