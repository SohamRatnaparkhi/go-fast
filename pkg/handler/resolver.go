package handler

import (
	"reflect"

	handlerresolvers "github.com/sohamratnaparkhi/go-fast/pkg/handler/resolvers"
)

type FieldIndexProvider = handlerresolvers.FieldIndexProvider
type FieldValueResolver = handlerresolvers.FieldValueResolver
type FieldResolver = handlerresolvers.FieldResolver

type BodyResolver = handlerresolvers.BodyResolver
type HeaderResolver = handlerresolvers.HeaderResolver
type QueryResolver = handlerresolvers.QueryResolver
type PathVarResolver = handlerresolvers.PathVarResolver
type CookieResolver = handlerresolvers.CookieResolver

// NewBodyResolver constructs a resolver for json:"body" fields.
func NewBodyResolver(fieldIdx int, fieldType reflect.Type) *BodyResolver {
	return handlerresolvers.NewBodyResolver(fieldIdx, fieldType)
}

// NewHeaderResolver constructs a resolver for json:"header:<name>" fields.
func NewHeaderResolver(fieldIdx int, headerName string, fieldType reflect.Type) *HeaderResolver {
	return handlerresolvers.NewHeaderResolver(fieldIdx, headerName, fieldType)
}

// NewQueryResolver constructs a resolver for json:"query:<name>" fields.
func NewQueryResolver(fieldIdx int, queryName string, fieldType reflect.Type) *QueryResolver {
	return handlerresolvers.NewQueryResolver(fieldIdx, queryName, fieldType)
}

// NewPathVarResolver constructs a resolver for json:"path:<name>" fields.
func NewPathVarResolver(fieldIdx int, paramName string, fieldType reflect.Type) *PathVarResolver {
	return handlerresolvers.NewPathVarResolver(fieldIdx, paramName, fieldType)
}

// NewCookieResolver constructs a resolver for json:"cookie:<name>" fields.
func NewCookieResolver(fieldIdx int, cookieName string, fieldType reflect.Type) *CookieResolver {
	return handlerresolvers.NewCookieResolver(fieldIdx, cookieName, fieldType)
}
