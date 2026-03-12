package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	domain "github.com/k07g/g1/internal/domain/item"
	usecase "github.com/k07g/g1/internal/usecase/item"
)

type ItemHandler struct {
	uc usecase.UseCase
}

func NewItemHandler(uc usecase.UseCase) *ItemHandler {
	return &ItemHandler{uc: uc}
}

// ---- リクエスト/レスポンス構造体 ----

type createRequest struct {
	Name  string `json:"name"`
	Price int    `json:"price"`
}

type updateRequest struct {
	Name  string `json:"name"`
	Price int    `json:"price"`
}

type apiResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// ---- ルーティング ----

func (h *ItemHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	path := r.URL.Path
	isCollection := path == "/api/items" || path == "/api/items/"

	if isCollection {
		switch r.Method {
		case http.MethodGet:
			h.list(w, r)
		case http.MethodPost:
			h.create(w, r)
		default:
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.get(w, r)
	case http.MethodPut:
		h.update(w, r)
	case http.MethodDelete:
		h.delete(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

// ---- ハンドラー実装 ----

func (h *ItemHandler) list(w http.ResponseWriter, r *http.Request) {
	items, err := h.uc.List()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeSuccess(w, http.StatusOK, items)
}

func (h *ItemHandler) get(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r.URL.Path)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	item, err := h.uc.Get(id)
	if err != nil {
		writeError(w, domainErrStatus(err), err.Error())
		return
	}
	writeSuccess(w, http.StatusOK, item)
}

func (h *ItemHandler) create(w http.ResponseWriter, r *http.Request) {
	var req createRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	item, err := h.uc.Create(usecase.CreateInput{
		Name:  req.Name,
		Price: req.Price,
	})
	if err != nil {
		writeError(w, domainErrStatus(err), err.Error())
		return
	}
	writeSuccess(w, http.StatusCreated, item)
}

func (h *ItemHandler) update(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r.URL.Path)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req updateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	item, err := h.uc.Update(usecase.UpdateInput{
		ID:    id,
		Name:  req.Name,
		Price: req.Price,
	})
	if err != nil {
		writeError(w, domainErrStatus(err), err.Error())
		return
	}
	writeSuccess(w, http.StatusOK, item)
}

func (h *ItemHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r.URL.Path)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.uc.Delete(id); err != nil {
		writeError(w, domainErrStatus(err), err.Error())
		return
	}
	writeSuccess(w, http.StatusOK, map[string]string{"message": "deleted"})
}

// ---- ヘルパー ----

func extractID(path string) (int, error) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) < 3 {
		return 0, fmt.Errorf("missing id")
	}
	return strconv.Atoi(parts[2])
}

func domainErrStatus(err error) int {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, domain.ErrInvalidName), errors.Is(err, domain.ErrInvalidPrice):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

func setCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func writeSuccess(w http.ResponseWriter, status int, data interface{}) {
	writeJSON(w, status, apiResponse{Success: true, Data: data})
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, apiResponse{Success: false, Error: msg})
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
