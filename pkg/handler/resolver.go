package handler

import (
	"reflect"

	handlerResolvers "github.com/sohamratnaparkhi/go-fast/pkg/handler/resolvers"
)

type FieldIndexProvider = handlerResolvers.FieldIndexProvider
type FieldValueResolver = handlerResolvers.FieldValueResolver
type FieldResolver = handlerResolvers.FieldResolver

type BodyResolver = handlerResolvers.BodyResolver
type HeaderResolver = handlerResolvers.HeaderResolver
type QueryResolver = handlerResolvers.QueryResolver
type PathVarResolver = handlerResolvers.PathVarResolver
type CookieResolver = handlerResolvers.CookieResolver
type FormResolver = handlerResolvers.FormResolver
type FileResolver = handlerResolvers.FileResolver

// NewBodyResolver constructs a resolver for gofast:"body" fields.
func NewBodyResolver(fieldIdx int, fieldType reflect.Type) *BodyResolver {
	return handlerResolvers.NewBodyResolver(fieldIdx, fieldType)
}

// NewHeaderResolver constructs a resolver for gofast:"header:<name>" fields.
func NewHeaderResolver(fieldIdx int, headerName string, fieldType reflect.Type) *HeaderResolver {
	return handlerResolvers.NewHeaderResolver(fieldIdx, headerName, fieldType)
}

// NewQueryResolver constructs a resolver for gofast:"query:<name>" fields.
func NewQueryResolver(fieldIdx int, queryName string, fieldType reflect.Type) *QueryResolver {
	return handlerResolvers.NewQueryResolver(fieldIdx, queryName, fieldType)
}

// NewPathVarResolver constructs a resolver for gofast:"path:<name>" fields.
func NewPathVarResolver(fieldIdx int, paramName string, fieldType reflect.Type) *PathVarResolver {
	return handlerResolvers.NewPathVarResolver(fieldIdx, paramName, fieldType)
}

// NewCookieResolver constructs a resolver for gofast:"cookie:<name>" fields.
func NewCookieResolver(fieldIdx int, cookieName string, fieldType reflect.Type) *CookieResolver {
	return handlerResolvers.NewCookieResolver(fieldIdx, cookieName, fieldType)
}

// NewFormResolver constructs a resolver for gofast:"form:<name>" fields.
func NewFormResolver(fieldIdx int, formName string, fieldType reflect.Type) *FormResolver {
	return handlerResolvers.NewFormResolver(fieldIdx, formName, fieldType)
}

// NewFileResolver constructs a resolver for gofast:"file:<name>" fields.
func NewFileResolver(fieldIdx int, fileName string, maxMemory int64) *FileResolver {
	return handlerResolvers.NewFileResolver(fieldIdx, fileName, maxMemory)
}
