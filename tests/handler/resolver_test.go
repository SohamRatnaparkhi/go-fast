package handler_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	handler "github.com/sohamratnaparkhi/go-fast/pkg/handler"
)

type testBody struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestBodyResolver_ResolveStruct(t *testing.T) {
	resolver := handler.NewBodyResolver(2, reflect.TypeOf(testBody{}))
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(`{"name":"john","age":21}`))
	ctx := &handler.Context{Request: req}

	value, err := resolver.Resolve(ctx)
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}

	if resolver.FieldIndex() != 2 {
		t.Fatalf("FieldIndex() = %d, want 2", resolver.FieldIndex())
	}

	got, ok := value.Interface().(testBody)
	if !ok {
		t.Fatalf("resolved value type = %T, want testBody", value.Interface())
	}

	if got.Name != "john" || got.Age != 21 {
		t.Fatalf("resolved value = %+v, want {Name:john Age:21}", got)
	}
}

func TestBodyResolver_ResolvePointer(t *testing.T) {
	resolver := handler.NewBodyResolver(0, reflect.TypeOf(&testBody{}))
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(`{"name":"sam","age":33}`))
	ctx := &handler.Context{Request: req}

	value, err := resolver.Resolve(ctx)
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}

	got, ok := value.Interface().(*testBody)
	if !ok {
		t.Fatalf("resolved value type = %T, want *testBody", value.Interface())
	}

	if got.Name != "sam" || got.Age != 33 {
		t.Fatalf("resolved value = %+v, want {Name:sam Age:33}", got)
	}
}

func TestBodyResolver_Resolve_InvalidJSON(t *testing.T) {
	resolver := handler.NewBodyResolver(0, reflect.TypeOf(testBody{}))
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(`{"name":`))
	ctx := &handler.Context{Request: req}

	_, err := resolver.Resolve(ctx)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestHeaderResolver_Resolve(t *testing.T) {
	resolver := handler.NewHeaderResolver(1, "X-Retry", reflect.TypeOf(0))
	req := httptest.NewRequest(http.MethodGet, "/items", nil)
	req.Header.Set("X-Retry", "7")
	ctx := &handler.Context{Request: req}

	value, err := resolver.Resolve(ctx)
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}

	if value.Int() != 7 {
		t.Fatalf("header resolved value = %d, want 7", value.Int())
	}
}

func TestQueryResolver_Resolve(t *testing.T) {
	resolver := handler.NewQueryResolver(1, "active", reflect.TypeOf(true))
	req := httptest.NewRequest(http.MethodGet, "/users?active=true", nil)
	ctx := &handler.Context{Request: req}

	value, err := resolver.Resolve(ctx)
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}

	if !value.Bool() {
		t.Fatal("query resolved value = false, want true")
	}
}

func TestPathVarResolver_Resolve(t *testing.T) {
	resolver := handler.NewPathVarResolver(1, "id", reflect.TypeOf(uint64(0)))
	ctx := &handler.Context{
		Request: httptest.NewRequest(http.MethodGet, "/users/42", nil),
		Params:  map[string]string{"id": "42"},
	}

	value, err := resolver.Resolve(ctx)
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}

	if value.Uint() != 42 {
		t.Fatalf("path var resolved value = %d, want 42", value.Uint())
	}
}

func TestPathVarResolver_Resolve_MissingParam(t *testing.T) {
	resolver := handler.NewPathVarResolver(1, "id", reflect.TypeOf(0))
	ctx := &handler.Context{
		Request: httptest.NewRequest(http.MethodGet, "/users", nil),
		Params:  map[string]string{},
	}

	_, err := resolver.Resolve(ctx)
	if err == nil {
		t.Fatal("expected error for missing path variable, got nil")
	}
}

func TestCookieResolver_Resolve(t *testing.T) {
	resolver := handler.NewCookieResolver(1, "session", reflect.TypeOf(""))
	req := httptest.NewRequest(http.MethodGet, "/profile", nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: "abc123"})
	ctx := &handler.Context{Request: req}

	value, err := resolver.Resolve(ctx)
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}

	if value.String() != "abc123" {
		t.Fatalf("cookie resolved value = %q, want %q", value.String(), "abc123")
	}
}

func TestCookieResolver_Resolve_MissingCookie(t *testing.T) {
	resolver := handler.NewCookieResolver(1, "session", reflect.TypeOf(""))
	ctx := &handler.Context{Request: httptest.NewRequest(http.MethodGet, "/profile", nil)}

	_, err := resolver.Resolve(ctx)
	if err == nil {
		t.Fatal("expected error for missing cookie, got nil")
	}
}
