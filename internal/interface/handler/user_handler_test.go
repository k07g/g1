package handler_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	domain "github.com/k07g/g1/internal/domain/user"
	"github.com/k07g/g1/internal/interface/handler"
	usecase "github.com/k07g/g1/internal/usecase/user"
)

// ---- モックユースケース ----

type mockUserUseCase struct {
	listOut   []*usecase.Output
	getOut    *usecase.Output
	createOut *usecase.Output
	updateOut *usecase.Output
	err       error
}

func (m *mockUserUseCase) List() ([]*usecase.Output, error)   { return m.listOut, m.err }
func (m *mockUserUseCase) Get(_ int) (*usecase.Output, error) { return m.getOut, m.err }
func (m *mockUserUseCase) Create(_ usecase.CreateInput) (*usecase.Output, error) {
	return m.createOut, m.err
}
func (m *mockUserUseCase) Update(_ usecase.UpdateInput) (*usecase.Output, error) {
	return m.updateOut, m.err
}
func (m *mockUserUseCase) Delete(_ int) error { return m.err }

// ---- GET /api/users ----

func TestUserHandler_List_OK(t *testing.T) {
	t.Parallel()
	uc := &mockUserUseCase{
		listOut: []*usecase.Output{
			{ID: 1, Name: "ユーザーA", Age: 20},
			{ID: 2, Name: "ユーザーB", Age: 30},
		},
	}
	h := handler.NewUserHandler(uc)

	rr := doRequest(t, h, http.MethodGet, "/api/users", nil)
	if rr.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusOK)
	}
	resp := decodeResp(t, rr)
	if !resp.Success {
		t.Errorf("success = false, want true")
	}
}

func TestUserHandler_List_InternalError(t *testing.T) {
	t.Parallel()
	uc := &mockUserUseCase{err: errors.New("db error")}
	h := handler.NewUserHandler(uc)

	rr := doRequest(t, h, http.MethodGet, "/api/users", nil)
	if rr.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusInternalServerError)
	}
}

// ---- GET /api/users/{id} ----

func TestUserHandler_Get_OK(t *testing.T) {
	t.Parallel()
	uc := &mockUserUseCase{getOut: &usecase.Output{ID: 1, Name: "山田太郎", Age: 25}}
	h := handler.NewUserHandler(uc)

	rr := doRequest(t, h, http.MethodGet, "/api/users/1", nil)
	if rr.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusOK)
	}
}

func TestUserHandler_Get_NotFound(t *testing.T) {
	t.Parallel()
	uc := &mockUserUseCase{err: domain.ErrNotFound}
	h := handler.NewUserHandler(uc)

	rr := doRequest(t, h, http.MethodGet, "/api/users/9999", nil)
	if rr.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusNotFound)
	}
}

func TestUserHandler_Get_InvalidID(t *testing.T) {
	t.Parallel()
	h := handler.NewUserHandler(&mockUserUseCase{})

	rr := doRequest(t, h, http.MethodGet, "/api/users/abc", nil)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusBadRequest)
	}
}

// ---- POST /api/users ----

func TestUserHandler_Create_OK(t *testing.T) {
	t.Parallel()
	uc := &mockUserUseCase{createOut: &usecase.Output{ID: 1, Name: "新ユーザー", Age: 18}}
	h := handler.NewUserHandler(uc)

	rr := doRequest(t, h, http.MethodPost, "/api/users", map[string]interface{}{
		"name": "新ユーザー", "age": 18,
	})
	if rr.Code != http.StatusCreated {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusCreated)
	}
	resp := decodeResp(t, rr)
	if !resp.Success {
		t.Errorf("success = false, want true")
	}
}

func TestUserHandler_Create_ValidationError(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		ucErr  error
		wantSt int
	}{
		{"名前が空", domain.ErrInvalidName, http.StatusBadRequest},
		{"年齢マイナス", domain.ErrInvalidAge, http.StatusBadRequest},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			uc := &mockUserUseCase{err: tt.ucErr}
			h := handler.NewUserHandler(uc)

			rr := doRequest(t, h, http.MethodPost, "/api/users", map[string]interface{}{
				"name": "", "age": -1,
			})
			if rr.Code != tt.wantSt {
				t.Errorf("status = %d, want %d", rr.Code, tt.wantSt)
			}
		})
	}
}

func TestUserHandler_Create_InvalidBody(t *testing.T) {
	t.Parallel()
	h := handler.NewUserHandler(&mockUserUseCase{})

	req := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewBufferString("not-json"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusBadRequest)
	}
}

// ---- PUT /api/users/{id} ----

func TestUserHandler_Update_OK(t *testing.T) {
	t.Parallel()
	uc := &mockUserUseCase{updateOut: &usecase.Output{ID: 1, Name: "更新後", Age: 21}}
	h := handler.NewUserHandler(uc)

	rr := doRequest(t, h, http.MethodPut, "/api/users/1", map[string]interface{}{
		"name": "更新後", "age": 21,
	})
	if rr.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusOK)
	}
}

func TestUserHandler_Update_NotFound(t *testing.T) {
	t.Parallel()
	uc := &mockUserUseCase{err: domain.ErrNotFound}
	h := handler.NewUserHandler(uc)

	rr := doRequest(t, h, http.MethodPut, "/api/users/9999", map[string]interface{}{
		"name": "x",
	})
	if rr.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusNotFound)
	}
}

// ---- DELETE /api/users/{id} ----

func TestUserHandler_Delete_OK(t *testing.T) {
	t.Parallel()
	h := handler.NewUserHandler(&mockUserUseCase{})

	rr := doRequest(t, h, http.MethodDelete, "/api/users/1", nil)
	if rr.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusOK)
	}
}

func TestUserHandler_Delete_NotFound(t *testing.T) {
	t.Parallel()
	uc := &mockUserUseCase{err: domain.ErrNotFound}
	h := handler.NewUserHandler(uc)

	rr := doRequest(t, h, http.MethodDelete, "/api/users/9999", nil)
	if rr.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusNotFound)
	}
}

// ---- Method Not Allowed ----

func TestUserHandler_MethodNotAllowed(t *testing.T) {
	t.Parallel()
	h := handler.NewUserHandler(&mockUserUseCase{})

	rr := doRequest(t, h, http.MethodPatch, "/api/users", nil)
	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusMethodNotAllowed)
	}
}

// ---- CORS preflight ----

func TestUserHandler_CORS_Preflight(t *testing.T) {
	t.Parallel()
	h := handler.NewUserHandler(&mockUserUseCase{})

	req := httptest.NewRequest(http.MethodOptions, "/api/users", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusNoContent)
	}
	if rr.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Error("missing CORS header")
	}
}
