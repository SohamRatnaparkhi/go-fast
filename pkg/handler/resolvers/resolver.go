package resolvers

import "reflect"

// FieldIndexProvider exposes the target field index in the input struct.
type FieldIndexProvider interface {
	FieldIndex() int
}

// FieldValueResolver resolves a value for a tagged field from request context.
type FieldValueResolver interface {
	Resolve(ctx *Context) (reflect.Value, error)
}

// FieldResolver resolves one tagged struct field from request context.
type FieldResolver interface {
	FieldIndexProvider
	FieldValueResolver
}
