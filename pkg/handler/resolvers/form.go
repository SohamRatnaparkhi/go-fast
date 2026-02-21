package resolvers

import (
	"fmt"
	"reflect"
)

// FormResolver resolves a form field value into the destination field type.
//
// It reads from the POST body only (not URL query parameters) using
// request.PostFormValue. Works with both application/x-www-form-urlencoded
// and multipart/form-data content types.
type FormResolver struct {
	fieldIdx  int
	formName  string
	fieldType reflect.Type
}

var _ FieldResolver = (*FormResolver)(nil)

// NewFormResolver constructs a resolver for gofast:"form:<name>" fields.
func NewFormResolver(fieldIdx int, formName string, fieldType reflect.Type) *FormResolver {
	return &FormResolver{fieldIdx: fieldIdx, formName: formName, fieldType: fieldType}
}

func (r *FormResolver) FieldIndex() int { return r.fieldIdx }

func (r *FormResolver) Resolve(ctx *Context) (reflect.Value, error) {
	if ctx == nil || ctx.Request == nil {
		return reflect.Value{}, fmt.Errorf("request context is nil")
	}

	raw := ctx.Request.PostFormValue(r.formName)
	value, err := convertStringToType(raw, r.fieldType)
	if err != nil {
		return reflect.Value{}, fmt.Errorf("resolve form %q: %w", r.formName, err)
	}

	return value, nil
}
