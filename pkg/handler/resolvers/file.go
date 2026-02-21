package resolvers

import (
	"fmt"
	"mime/multipart"
	"reflect"
)

// defaultMaxMemory is the maximum bytes stored in memory for multipart parsing.
// Files beyond this limit are written to temporary files on disk.
const defaultMaxMemory = 32 << 20 // 32 MB

// MultipartFileHeaderType is the reflect.Type for *multipart.FileHeader.
// Exported so the resolver compiler can validate field types at startup.
var MultipartFileHeaderType = reflect.TypeOf((*multipart.FileHeader)(nil))

// FileResolver resolves an uploaded file from a multipart/form-data request.
//
// The field type must be *multipart.FileHeader. The resolver accesses the
// parsed multipart form directly, avoiding an unnecessary file open.
type FileResolver struct {
	fieldIdx  int
	fileName  string
	maxMemory int64
}

var _ FieldResolver = (*FileResolver)(nil)

// NewFileResolver constructs a resolver for gofast:"file:<name>" fields.
// maxMemory controls the maximum bytes stored in memory for multipart parsing;
// pass 0 to use the default (32 MB).
func NewFileResolver(fieldIdx int, fileName string, maxMemory int64) *FileResolver {
	if maxMemory <= 0 {
		maxMemory = defaultMaxMemory
	}
	return &FileResolver{fieldIdx: fieldIdx, fileName: fileName, maxMemory: maxMemory}
}

func (r *FileResolver) FieldIndex() int { return r.fieldIdx }

func (r *FileResolver) Resolve(ctx *Context) (reflect.Value, error) {
	if ctx == nil || ctx.Request == nil {
		return reflect.Value{}, fmt.Errorf("request context is nil")
	}

	if err := ctx.Request.ParseMultipartForm(r.maxMemory); err != nil {
		return reflect.Value{}, fmt.Errorf("resolve file %q: %w", r.fileName, err)
	}

	if ctx.Request.MultipartForm == nil || ctx.Request.MultipartForm.File == nil {
		return reflect.Value{}, fmt.Errorf("resolve file %q: no multipart form data", r.fileName)
	}

	fhs := ctx.Request.MultipartForm.File[r.fileName]
	if len(fhs) == 0 {
		return reflect.Value{}, fmt.Errorf("resolve file %q: file not found", r.fileName)
	}

	return reflect.ValueOf(fhs[0]), nil
}
