package resolvers

import (
	"fmt"
	"reflect"
)

// QueryResolver resolves a query string value into the destination field type.
type QueryResolver struct {
	fieldIdx  int
	queryName string
	fieldType reflect.Type
}

var _ FieldResolver = (*QueryResolver)(nil)

// NewQueryResolver constructs a resolver for json:"query:<name>" fields.
func NewQueryResolver(fieldIdx int, queryName string, fieldType reflect.Type) *QueryResolver {
	return &QueryResolver{fieldIdx: fieldIdx, queryName: queryName, fieldType: fieldType}
}

func (r *QueryResolver) FieldIndex() int { return r.fieldIdx }

func (r *QueryResolver) Resolve(ctx *Context) (reflect.Value, error) {
	if ctx == nil || ctx.Request == nil {
		return reflect.Value{}, fmt.Errorf("request context is nil")
	}

	raw := ctx.Request.URL.Query().Get(r.queryName)
	value, err := convertStringToType(raw, r.fieldType)
	if err != nil {
		return reflect.Value{}, fmt.Errorf("resolve query %q: %w", r.queryName, err)
	}

	return value, nil
}
