package resolvers

import (
	"fmt"
	"reflect"
)

// HeaderResolver resolves a request header into the destination field type.
type HeaderResolver struct {
	fieldIdx   int
	headerName string
	fieldType  reflect.Type
}

var _ FieldResolver = (*HeaderResolver)(nil)

// NewHeaderResolver constructs a resolver for gofast:"header:<name>" fields.
func NewHeaderResolver(fieldIdx int, headerName string, fieldType reflect.Type) *HeaderResolver {
	return &HeaderResolver{fieldIdx: fieldIdx, headerName: headerName, fieldType: fieldType}
}

func (r *HeaderResolver) FieldIndex() int { return r.fieldIdx }

func (r *HeaderResolver) Resolve(ctx *Context) (reflect.Value, error) {
	if ctx == nil || ctx.Request == nil {
		return reflect.Value{}, fmt.Errorf("request context is nil")
	}

	raw := ctx.Request.Header.Get(r.headerName)
	value, err := convertStringToType(raw, r.fieldType)
	if err != nil {
		return reflect.Value{}, fmt.Errorf("resolve header %q: %w", r.headerName, err)
	}

	return value, nil
}
