//go:build ignore

package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type testBody struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestBodyResolver_ResolveStruct(t *testing.T) {
	resolver := NewBodyResolver(2, reflect.TypeOf(testBody{}))
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(`{"name":"john","age":21}`))
	ctx := &Context{Request: req}

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
	resolver := NewBodyResolver(0, reflect.TypeOf(&testBody{}))
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(`{"name":"sam","age":33}`))
	ctx := &Context{Request: req}

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
	resolver := NewBodyResolver(0, reflect.TypeOf(testBody{}))
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(`{"name":`))
	ctx := &Context{Request: req}

	_, err := resolver.Resolve(ctx)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestHeaderResolver_Resolve(t *testing.T) {
	resolver := NewHeaderResolver(1, "X-Retry", reflect.TypeOf(0))
	req := httptest.NewRequest(http.MethodGet, "/items", nil)
	req.Header.Set("X-Retry", "7")
	ctx := &Context{Request: req}

	value, err := resolver.Resolve(ctx)
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}

	if value.Int() != 7 {
		t.Fatalf("header resolved value = %d, want 7", value.Int())
	}
}

func TestQueryResolver_Resolve(t *testing.T) {
	resolver := NewQueryResolver(1, "active", reflect.TypeOf(true))
	req := httptest.NewRequest(http.MethodGet, "/users?active=true", nil)
	ctx := &Context{Request: req}

	value, err := resolver.Resolve(ctx)
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}

	if !value.Bool() {
		t.Fatal("query resolved value = false, want true")
	}
}

func TestPathVarResolver_Resolve(t *testing.T) {
	resolver := NewPathVarResolver(1, "id", reflect.TypeOf(uint64(0)))
	ctx := &Context{
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
	resolver := NewPathVarResolver(1, "id", reflect.TypeOf(0))
	ctx := &Context{
		Request: httptest.NewRequest(http.MethodGet, "/users", nil),
		Params:  map[string]string{},
	}

	_, err := resolver.Resolve(ctx)
	if err == nil {
		t.Fatal("expected error for missing path variable, got nil")
	}
}

func TestCookieResolver_Resolve(t *testing.T) {
	resolver := NewCookieResolver(1, "session", reflect.TypeOf(""))
	req := httptest.NewRequest(http.MethodGet, "/profile", nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: "abc123"})
	ctx := &Context{Request: req}

	value, err := resolver.Resolve(ctx)
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}

	if value.String() != "abc123" {
		t.Fatalf("cookie resolved value = %q, want %q", value.String(), "abc123")
	}
}

func TestCookieResolver_Resolve_MissingCookie(t *testing.T) {
	resolver := NewCookieResolver(1, "session", reflect.TypeOf(""))
	ctx := &Context{Request: httptest.NewRequest(http.MethodGet, "/profile", nil)}

	_, err := resolver.Resolve(ctx)
	if err == nil {
		t.Fatal("expected error for missing cookie, got nil")
	}
}

func TestConvertStringToType(t *testing.T) {
	t.Run("pointer int", func(t *testing.T) {
		ptrType := reflect.TypeOf((*int)(nil))
		value, err := convertStringToType("9", ptrType)
		if err != nil {
			t.Fatalf("convertStringToType() error = %v", err)
		}
		if value.Elem().Int() != 9 {
			t.Fatalf("pointer inner value = %d, want 9", value.Elem().Int())
		}
	})

	t.Run("empty gives zero", func(t *testing.T) {
		value, err := convertStringToType("", reflect.TypeOf(0))
		if err != nil {
			t.Fatalf("convertStringToType() error = %v", err)
		}
		if value.Int() != 0 {
			t.Fatalf("value = %d, want 0", value.Int())
		}
	})

	t.Run("nil type", func(t *testing.T) {
		_, err := convertStringToType("x", nil)
		if err == nil {
			t.Fatal("expected error for nil field type, got nil")
		}
	})

	t.Run("unsupported type", func(t *testing.T) {
		_, err := convertStringToType("x", reflect.TypeOf(testBody{}))
		if err == nil {
			t.Fatal("expected error for unsupported type, got nil")
		}
	})
}
