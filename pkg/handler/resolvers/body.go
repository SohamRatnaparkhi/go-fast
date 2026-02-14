package resolvers

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// BodyResolver decodes request JSON body into a struct-typed field.
type BodyResolver struct {
	fieldIdx  int
	fieldType reflect.Type
}

var _ FieldResolver = (*BodyResolver)(nil)

// NewBodyResolver constructs a resolver for json:"body" fields.
func NewBodyResolver(fieldIdx int, fieldType reflect.Type) *BodyResolver {
	return &BodyResolver{fieldIdx: fieldIdx, fieldType: fieldType}
}

func (r *BodyResolver) FieldIndex() int { return r.fieldIdx }

func (r *BodyResolver) Resolve(ctx *Context) (reflect.Value, error) {
	if ctx == nil || ctx.Request == nil {
		return reflect.Value{}, fmt.Errorf("request context is nil")
	}

	if r.fieldType.Kind() == reflect.Ptr {
		instance := reflect.New(r.fieldType.Elem())
		if err := json.NewDecoder(ctx.Request.Body).Decode(instance.Interface()); err != nil {
			return reflect.Value{}, fmt.Errorf("decode body: %w", err)
		}
		return instance, nil
	}

	instance := reflect.New(r.fieldType)
	if err := json.NewDecoder(ctx.Request.Body).Decode(instance.Interface()); err != nil {
		return reflect.Value{}, fmt.Errorf("decode body: %w", err)
	}

	return instance.Elem(), nil
}
