package resolvers

import (
	"fmt"
	"reflect"
)

// CookieResolver resolves a cookie value into the destination field type.
type CookieResolver struct {
	fieldIdx   int
	cookieName string
	fieldType  reflect.Type
}

var _ FieldResolver = (*CookieResolver)(nil)

// NewCookieResolver constructs a resolver for json:"cookie:<name>" fields.
func NewCookieResolver(fieldIdx int, cookieName string, fieldType reflect.Type) *CookieResolver {
	return &CookieResolver{fieldIdx: fieldIdx, cookieName: cookieName, fieldType: fieldType}
}

func (r *CookieResolver) FieldIndex() int { return r.fieldIdx }

func (r *CookieResolver) Resolve(ctx *Context) (reflect.Value, error) {
	if ctx == nil || ctx.Request == nil {
		return reflect.Value{}, fmt.Errorf("request context is nil")
	}

	cookie, err := ctx.Request.Cookie(r.cookieName)
	if err != nil {
		return reflect.Value{}, fmt.Errorf("resolve cookie %q: %w", r.cookieName, err)
	}

	value, err := convertStringToType(cookie.Value, r.fieldType)
	if err != nil {
		return reflect.Value{}, fmt.Errorf("resolve cookie %q: %w", r.cookieName, err)
	}

	return value, nil
}
