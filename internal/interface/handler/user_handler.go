package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	domain "github.com/k07g/g1/internal/domain/user"
	usecase "github.com/k07g/g1/internal/usecase/user"
)

type UserHandler struct {
	uc usecase.UseCase
}

func NewUserHandler(uc usecase.UseCase) *UserHandler {
	return &UserHandler{uc: uc}
}

// ---- リクエスト構造体 ----

type createUserRequest struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type updateUserRequest struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// ---- ルーティング ----

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	path := r.URL.Path
	isCollection := path == "/api/users" || path == "/api/users/"

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

func (h *UserHandler) list(w http.ResponseWriter, r *http.Request) {
	users, err := h.uc.List()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeSuccess(w, http.StatusOK, users)
}

func (h *UserHandler) get(w http.ResponseWriter, r *http.Request) {
	id, err := extractUserID(r.URL.Path)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	u, err := h.uc.Get(id)
	if err != nil {
		writeError(w, userErrStatus(err), err.Error())
		return
	}
	writeSuccess(w, http.StatusOK, u)
}

func (h *UserHandler) create(w http.ResponseWriter, r *http.Request) {
	var req createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	u, err := h.uc.Create(usecase.CreateInput{
		Name: req.Name,
		Age:  req.Age,
	})
	if err != nil {
		writeError(w, userErrStatus(err), err.Error())
		return
	}
	writeSuccess(w, http.StatusCreated, u)
}

func (h *UserHandler) update(w http.ResponseWriter, r *http.Request) {
	id, err := extractUserID(r.URL.Path)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req updateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	u, err := h.uc.Update(usecase.UpdateInput{
		ID:   id,
		Name: req.Name,
		Age:  req.Age,
	})
	if err != nil {
		writeError(w, userErrStatus(err), err.Error())
		return
	}
	writeSuccess(w, http.StatusOK, u)
}

func (h *UserHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := extractUserID(r.URL.Path)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.uc.Delete(id); err != nil {
		writeError(w, userErrStatus(err), err.Error())
		return
	}
	writeSuccess(w, http.StatusOK, map[string]string{"message": "deleted"})
}

// ---- ヘルパー ----

func extractUserID(path string) (int, error) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) < 3 {
		return 0, fmt.Errorf("missing id")
	}
	return strconv.Atoi(parts[2])
}

func userErrStatus(err error) int {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, domain.ErrInvalidName), errors.Is(err, domain.ErrInvalidAge):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
