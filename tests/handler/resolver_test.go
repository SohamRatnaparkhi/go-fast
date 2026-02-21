package handler_test

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"reflect"
	"strings"
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

// --- Form Resolver Tests ---

func TestFormResolver_Resolve(t *testing.T) {
	resolver := handler.NewFormResolver(1, "username", reflect.TypeOf(""))
	req := httptest.NewRequest(http.MethodPost, "/submit",
		strings.NewReader("username=alice&email=alice@test.com"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	ctx := &handler.Context{Request: req}

	value, err := resolver.Resolve(ctx)
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}

	if value.String() != "alice" {
		t.Fatalf("form resolved value = %q, want %q", value.String(), "alice")
	}
}

func TestFormResolver_Resolve_TypeConversion(t *testing.T) {
	resolver := handler.NewFormResolver(1, "age", reflect.TypeOf(0))
	req := httptest.NewRequest(http.MethodPost, "/submit",
		strings.NewReader("age=25"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	ctx := &handler.Context{Request: req}

	value, err := resolver.Resolve(ctx)
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}

	if value.Int() != 25 {
		t.Fatalf("form resolved value = %d, want 25", value.Int())
	}
}

func TestFormResolver_Resolve_MissingValue(t *testing.T) {
	resolver := handler.NewFormResolver(1, "missing", reflect.TypeOf(""))
	req := httptest.NewRequest(http.MethodPost, "/submit",
		strings.NewReader("username=alice"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	ctx := &handler.Context{Request: req}

	value, err := resolver.Resolve(ctx)
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}

	if value.String() != "" {
		t.Fatalf("form resolved value = %q, want empty string", value.String())
	}
}

func TestFormResolver_Resolve_MultipartFormData(t *testing.T) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	_ = w.WriteField("username", "bob")
	_ = w.WriteField("age", "30")
	w.Close()

	resolver := handler.NewFormResolver(1, "username", reflect.TypeOf(""))
	req := httptest.NewRequest(http.MethodPost, "/submit", &buf)
	req.Header.Set("Content-Type", w.FormDataContentType())
	ctx := &handler.Context{Request: req}

	value, err := resolver.Resolve(ctx)
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}

	if value.String() != "bob" {
		t.Fatalf("form resolved value = %q, want %q", value.String(), "bob")
	}
}

func TestFormResolver_FieldIndex(t *testing.T) {
	resolver := handler.NewFormResolver(3, "name", reflect.TypeOf(""))
	if resolver.FieldIndex() != 3 {
		t.Fatalf("FieldIndex() = %d, want 3", resolver.FieldIndex())
	}
}

// --- File Resolver Tests ---

func newMultipartFileRequest(t *testing.T, fieldName, fileName, content string) *http.Request {
	t.Helper()
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name=%q; filename=%q`, fieldName, fileName))
	h.Set("Content-Type", "application/octet-stream")
	part, err := w.CreatePart(h)
	if err != nil {
		t.Fatalf("create part: %v", err)
	}
	part.Write([]byte(content))
	w.Close()

	req := httptest.NewRequest(http.MethodPost, "/upload", &buf)
	req.Header.Set("Content-Type", w.FormDataContentType())
	return req
}

func TestFileResolver_Resolve(t *testing.T) {
	resolver := handler.NewFileResolver(1, "avatar", 0)
	req := newMultipartFileRequest(t, "avatar", "photo.png", "fake-image-data")
	ctx := &handler.Context{Request: req}

	value, err := resolver.Resolve(ctx)
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}

	fh, ok := value.Interface().(*multipart.FileHeader)
	if !ok {
		t.Fatalf("resolved value type = %T, want *multipart.FileHeader", value.Interface())
	}

	if fh.Filename != "photo.png" {
		t.Fatalf("file name = %q, want %q", fh.Filename, "photo.png")
	}

	if fh.Size != int64(len("fake-image-data")) {
		t.Fatalf("file size = %d, want %d", fh.Size, len("fake-image-data"))
	}
}

func TestFileResolver_Resolve_MissingFile(t *testing.T) {
	resolver := handler.NewFileResolver(1, "avatar", 0)
	req := newMultipartFileRequest(t, "document", "doc.pdf", "content")
	ctx := &handler.Context{Request: req}

	_, err := resolver.Resolve(ctx)
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestFileResolver_Resolve_NotMultipart(t *testing.T) {
	resolver := handler.NewFileResolver(1, "avatar", 0)
	req := httptest.NewRequest(http.MethodPost, "/upload",
		strings.NewReader("plain body"))
	ctx := &handler.Context{Request: req}

	_, err := resolver.Resolve(ctx)
	if err == nil {
		t.Fatal("expected error for non-multipart request, got nil")
	}
}

func TestFileResolver_FieldIndex(t *testing.T) {
	resolver := handler.NewFileResolver(5, "doc", 0)
	if resolver.FieldIndex() != 5 {
		t.Fatalf("FieldIndex() = %d, want 5", resolver.FieldIndex())
	}
}
