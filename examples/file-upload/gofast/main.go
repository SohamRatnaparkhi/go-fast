package main

import (
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/sohamratnaparkhi/go-fast/pkg/handler"
)

type UploadResponse struct {
	Title    string `json:"title"`
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
}

// UploadDocument handles a multipart form with both text fields and a file.
//
// The form: tag reads text fields, the file: tag reads uploaded files.
// The file field must be *multipart.FileHeader — validated at startup.
func UploadDocument(req struct {
	Title    string                `gofast:"form:title"`
	Document *multipart.FileHeader `gofast:"file:document"`
}) (*UploadResponse, error) {
	return &UploadResponse{
		Title:    req.Title,
		Filename: req.Document.Filename,
		Size:     req.Document.Size,
	}, nil
}

func main() {
	h, err := handler.Adapt(UploadDocument)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/upload", h)
	fmt.Println("go-fast server on :8080")
	fmt.Println(`
File upload with form fields — zero boilerplate:

  curl -X POST localhost:8080/upload \
    -F 'title=My Report' \
    -F 'document=@report.pdf'`)
	_ = http.ListenAndServe(":8080", nil)
}
