package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"strings"
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

// --- Form Adapter Tests ---

func TestAdapt_FormFields(t *testing.T) {
	type formInput struct {
		Name  string `json:"form:name"`
		Email string `json:"form:email"`
		Age   int    `json:"form:age"`
	}
	type formOutput struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Age   int    `json:"age"`
	}

	h, err := handler.Adapt(func(req formInput) (*formOutput, error) {
		return &formOutput{Name: req.Name, Email: req.Email, Age: req.Age}, nil
	})
	if err != nil {
		t.Fatalf("Adapt() error = %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/submit",
		strings.NewReader("name=alice&email=alice@test.com&age=30"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	h(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var got formOutput
	if decodeErr := json.NewDecoder(w.Body).Decode(&got); decodeErr != nil {
		t.Fatalf("decode response: %v", decodeErr)
	}

	if got.Name != "alice" || got.Email != "alice@test.com" || got.Age != 30 {
		t.Fatalf("unexpected form mapping: %+v", got)
	}
}

func TestAdapt_FormWithHeaderAndQuery(t *testing.T) {
	type input struct {
		Name   string `json:"form:name"`
		Token  string `json:"header:Authorization"`
		Format string `json:"query:format"`
	}
	type output struct {
		Name   string `json:"name"`
		Token  string `json:"token"`
		Format string `json:"format"`
	}

	h, err := handler.Adapt(func(req input) (*output, error) {
		return &output{Name: req.Name, Token: req.Token, Format: req.Format}, nil
	})
	if err != nil {
		t.Fatalf("Adapt() error = %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/submit?format=json",
		strings.NewReader("name=bob"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer xyz")

	w := httptest.NewRecorder()
	h(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var got output
	if decodeErr := json.NewDecoder(w.Body).Decode(&got); decodeErr != nil {
		t.Fatalf("decode response: %v", decodeErr)
	}

	if got.Name != "bob" || got.Token != "Bearer xyz" || got.Format != "json" {
		t.Fatalf("unexpected mapping: %+v", got)
	}
}

// --- File Adapter Tests ---

func TestAdapt_FileUpload(t *testing.T) {
	type uploadInput struct {
		Avatar *multipart.FileHeader `json:"file:avatar"`
	}
	type uploadOutput struct {
		Filename string `json:"filename"`
		Size     int64  `json:"size"`
	}

	h, err := handler.Adapt(func(req uploadInput) (*uploadOutput, error) {
		return &uploadOutput{
			Filename: req.Avatar.Filename,
			Size:     req.Avatar.Size,
		}, nil
	})
	if err != nil {
		t.Fatalf("Adapt() error = %v", err)
	}

	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	hdr := make(textproto.MIMEHeader)
	hdr.Set("Content-Disposition", `form-data; name="avatar"; filename="photo.png"`)
	hdr.Set("Content-Type", "image/png")
	part, _ := mw.CreatePart(hdr)
	part.Write([]byte("fake-png-data"))
	mw.Close()

	req := httptest.NewRequest(http.MethodPost, "/upload", &buf)
	req.Header.Set("Content-Type", mw.FormDataContentType())

	w := httptest.NewRecorder()
	h(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var got uploadOutput
	if decodeErr := json.NewDecoder(w.Body).Decode(&got); decodeErr != nil {
		t.Fatalf("decode response: %v", decodeErr)
	}

	if got.Filename != "photo.png" {
		t.Fatalf("filename = %q, want %q", got.Filename, "photo.png")
	}
	if got.Size != int64(len("fake-png-data")) {
		t.Fatalf("size = %d, want %d", got.Size, len("fake-png-data"))
	}
}

func TestAdapt_FormAndFile(t *testing.T) {
	type input struct {
		Title    string                `json:"form:title"`
		Document *multipart.FileHeader `json:"file:document"`
	}
	type output struct {
		Title    string `json:"title"`
		Filename string `json:"filename"`
		Size     int64  `json:"size"`
	}

	h, err := handler.Adapt(func(req input) (*output, error) {
		return &output{
			Title:    req.Title,
			Filename: req.Document.Filename,
			Size:     req.Document.Size,
		}, nil
	})
	if err != nil {
		t.Fatalf("Adapt() error = %v", err)
	}

	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	_ = mw.WriteField("title", "My Document")
	hdr := make(textproto.MIMEHeader)
	hdr.Set("Content-Disposition", `form-data; name="document"; filename="report.pdf"`)
	hdr.Set("Content-Type", "application/pdf")
	part, _ := mw.CreatePart(hdr)
	part.Write([]byte("pdf-content"))
	mw.Close()

	req := httptest.NewRequest(http.MethodPost, "/docs", &buf)
	req.Header.Set("Content-Type", mw.FormDataContentType())

	w := httptest.NewRecorder()
	h(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var got output
	if decodeErr := json.NewDecoder(w.Body).Decode(&got); decodeErr != nil {
		t.Fatalf("decode response: %v", decodeErr)
	}

	if got.Title != "My Document" || got.Filename != "report.pdf" {
		t.Fatalf("unexpected mapping: %+v", got)
	}
}

func TestAdapt_BodyAndForm_Error(t *testing.T) {
	type badInput struct {
		Body createUserBody `json:"body"`
		Name string         `json:"form:name"`
	}

	_, err := handler.Adapt(func(req badInput) error { return nil })
	if err == nil {
		t.Fatal("expected error for body+form combination, got nil")
	}
}

func TestAdapt_BodyAndFile_Error(t *testing.T) {
	type badInput struct {
		Body   createUserBody        `json:"body"`
		Avatar *multipart.FileHeader `json:"file:avatar"`
	}

	_, err := handler.Adapt(func(req badInput) error { return nil })
	if err == nil {
		t.Fatal("expected error for body+file combination, got nil")
	}
}

func TestAdapt_FileWrongType_Error(t *testing.T) {
	type badInput struct {
		Avatar string `json:"file:avatar"`
	}

	_, err := handler.Adapt(func(req badInput) error { return nil })
	if err == nil {
		t.Fatal("expected error for file field with wrong type, got nil")
	}
}

func TestAdapt_FileUpload_MissingFile_Returns400(t *testing.T) {
	type uploadInput struct {
		Avatar *multipart.FileHeader `json:"file:avatar"`
	}

	h, err := handler.Adapt(func(req uploadInput) (map[string]string, error) {
		return map[string]string{"filename": req.Avatar.Filename}, nil
	})
	if err != nil {
		t.Fatalf("Adapt() error = %v", err)
	}

	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	hdr := make(textproto.MIMEHeader)
	hdr.Set("Content-Disposition", `form-data; name="other_file"; filename="other.txt"`)
	hdr.Set("Content-Type", "text/plain")
	part, _ := mw.CreatePart(hdr)
	part.Write([]byte("wrong file"))
	mw.Close()

	req := httptest.NewRequest(http.MethodPost, "/upload", &buf)
	req.Header.Set("Content-Type", mw.FormDataContentType())

	w := httptest.NewRecorder()
	h(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

// Ensure unused imports are referenced.
var _ = fmt.Sprintf
var _ = errors.New
