package resolvers

import (
	"fmt"
	"reflect"
)

// PathVarResolver resolves a path parameter into the destination field type.
type PathVarResolver struct {
	fieldIdx  int
	paramName string
	fieldType reflect.Type
}

var _ FieldResolver = (*PathVarResolver)(nil)

// NewPathVarResolver constructs a resolver for gofast:"path:<name>" fields.
func NewPathVarResolver(fieldIdx int, paramName string, fieldType reflect.Type) *PathVarResolver {
	return &PathVarResolver{fieldIdx: fieldIdx, paramName: paramName, fieldType: fieldType}
}

func (r *PathVarResolver) FieldIndex() int { return r.fieldIdx }

func (r *PathVarResolver) Resolve(ctx *Context) (reflect.Value, error) {
	if ctx == nil || ctx.Request == nil {
		return reflect.Value{}, fmt.Errorf("request context is nil")
	}

	raw, ok := ctx.Params[r.paramName]
	if !ok {
		return reflect.Value{}, fmt.Errorf("path variable %q not found", r.paramName)
	}

	value, err := convertStringToType(raw, r.fieldType)
	if err != nil {
		return reflect.Value{}, fmt.Errorf("resolve path variable %q: %w", r.paramName, err)
	}

	return value, nil
}
