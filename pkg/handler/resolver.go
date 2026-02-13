package handler

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

type FieldResolver interface {
	Resolve(ctx *Context) (reflect.Value, error)
	FieldIndex() int
}

type BodyResolver struct {
	fieldIdx  int
	fieldType reflect.Type
}

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

type HeaderResolver struct {
	fieldIdx   int
	headerName string
	fieldType  reflect.Type
}

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

type QueryResolver struct {
	fieldIdx   int
	queryName  string
	fieldType  reflect.Type
}

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

type PathVarResolver struct {
	fieldIdx  int
	paramName string
	fieldType reflect.Type
}

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

type CookieResolver struct {
	fieldIdx   int
	cookieName string
	fieldType  reflect.Type
}

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

func convertStringToType(raw string, fieldType reflect.Type) (reflect.Value, error) {
	if fieldType == nil {
		return reflect.Value{}, fmt.Errorf("field type is nil")
	}

	if raw == "" {
		return reflect.Zero(fieldType), nil
	}

	switch fieldType.Kind() {
	case reflect.String:
		return reflect.ValueOf(raw).Convert(fieldType), nil
	case reflect.Bool:
		v, err := strconv.ParseBool(raw)
		if err != nil {
			return reflect.Value{}, err
		}
		value := reflect.New(fieldType).Elem()
		value.SetBool(v)
		return value, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v, err := strconv.ParseInt(raw, 10, fieldType.Bits())
		if err != nil {
			return reflect.Value{}, err
		}
		value := reflect.New(fieldType).Elem()
		value.SetInt(v)
		return value, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		v, err := strconv.ParseUint(raw, 10, fieldType.Bits())
		if err != nil {
			return reflect.Value{}, err
		}
		value := reflect.New(fieldType).Elem()
		value.SetUint(v)
		return value, nil
	case reflect.Float32, reflect.Float64:
		v, err := strconv.ParseFloat(raw, fieldType.Bits())
		if err != nil {
			return reflect.Value{}, err
		}
		value := reflect.New(fieldType).Elem()
		value.SetFloat(v)
		return value, nil
	case reflect.Ptr:
		innerValue, err := convertStringToType(raw, fieldType.Elem())
		if err != nil {
			return reflect.Value{}, err
		}
		ptrValue := reflect.New(fieldType.Elem())
		ptrValue.Elem().Set(innerValue)
		return ptrValue, nil
	default:
		return reflect.Value{}, fmt.Errorf("unsupported field type %s", fieldType)
	}
}