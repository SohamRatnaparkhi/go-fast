package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	handler "github.com/sohamratnaparkhi/go-fast/pkg/handler"
)

type createUserBody struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type createUserInput struct {
	Body    createUserBody `json:"body"`
	Token   string         `json:"header:Authorization"`
	Active  bool           `json:"query:active"`
	Session string         `json:"cookie:session"`
}

type createUserOutput struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Token   string `json:"token"`
	Active  bool   `json:"active"`
	Session string `json:"session"`
}

func TestAdapt_Success(t *testing.T) {
	h, err := handler.Adapt(func(req createUserInput) (*createUserOutput, error) {
		return &createUserOutput{
			Name:    req.Body.Name,
			Email:   req.Body.Email,
			Token:   req.Token,
			Active:  req.Active,
			Session: req.Session,
		}, nil
	})
	if err != nil {
		t.Fatalf("Adapt() error = %v", err)
	}

	body := bytes.NewBufferString(`{"name":"john","email":"john@test.com"}`)
	req := httptest.NewRequest(http.MethodPost, "/users?active=true", body)
	req.Header.Set("Authorization", "Bearer abc")
	req.AddCookie(&http.Cookie{Name: "session", Value: "sess-1"})

	w := httptest.NewRecorder()
	h(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var got createUserOutput
	if decodeErr := json.NewDecoder(w.Body).Decode(&got); decodeErr != nil {
		t.Fatalf("decode response: %v", decodeErr)
	}

	if got.Name != "john" || got.Email != "john@test.com" {
		t.Fatalf("unexpected body mapping: %+v", got)
	}
	if got.Token != "Bearer abc" || !got.Active || got.Session != "sess-1" {
		t.Fatalf("unexpected resolver mapping: %+v", got)
	}
}

func TestAdapt_HandlerReturnsError(t *testing.T) {
	h, err := handler.Adapt(func(req createUserInput) (*createUserOutput, error) {
		return nil, errors.New("boom")
	})
	if err != nil {
		t.Fatalf("Adapt() error = %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/users?active=true", bytes.NewBufferString(`{"name":"john","email":"john@test.com"}`))
	req.Header.Set("Authorization", "Bearer abc")
	req.AddCookie(&http.Cookie{Name: "session", Value: "sess-1"})

	w := httptest.NewRecorder()
	h(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestAdapt_RequiresSingleStructInput(t *testing.T) {
	if _, err := handler.Adapt(func(a, b string) {}); err == nil {
		t.Fatal("expected error for multiple inputs, got nil")
	}

	if _, err := handler.Adapt(func(a string) {}); err == nil {
		t.Fatal("expected error for non-struct input, got nil")
	}
}

func TestAdapt_PathFieldWithoutParams_Returns400(t *testing.T) {
	type input struct {
		ID int `json:"path:id"`
	}

	h, err := handler.Adapt(func(req input) (map[string]int, error) {
		return map[string]int{"id": req.ID}, nil
	})
	if err != nil {
		t.Fatalf("Adapt() error = %v", err)
	}

	w := httptest.NewRecorder()
	h(w, httptest.NewRequest(http.MethodGet, "/users/42", nil))

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}
