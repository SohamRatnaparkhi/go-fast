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

// NewBodyResolver constructs a resolver for json:"body" fields.
func NewBodyResolver(fieldIdx int, fieldType reflect.Type) *BodyResolver {
	return handlerResolvers.NewBodyResolver(fieldIdx, fieldType)
}

// NewHeaderResolver constructs a resolver for json:"header:<name>" fields.
func NewHeaderResolver(fieldIdx int, headerName string, fieldType reflect.Type) *HeaderResolver {
	return handlerResolvers.NewHeaderResolver(fieldIdx, headerName, fieldType)
}

// NewQueryResolver constructs a resolver for json:"query:<name>" fields.
func NewQueryResolver(fieldIdx int, queryName string, fieldType reflect.Type) *QueryResolver {
	return handlerResolvers.NewQueryResolver(fieldIdx, queryName, fieldType)
}

// NewPathVarResolver constructs a resolver for json:"path:<name>" fields.
func NewPathVarResolver(fieldIdx int, paramName string, fieldType reflect.Type) *PathVarResolver {
	return handlerResolvers.NewPathVarResolver(fieldIdx, paramName, fieldType)
}

// NewCookieResolver constructs a resolver for json:"cookie:<name>" fields.
func NewCookieResolver(fieldIdx int, cookieName string, fieldType reflect.Type) *CookieResolver {
	return handlerResolvers.NewCookieResolver(fieldIdx, cookieName, fieldType)
}
