// Package handler provides the reflection-based HTTP adapter used by go-fast.
//
// It analyzes a user function once at startup, builds field resolvers from
// struct tags, and exposes a standard net/http handler for runtime execution.
package handler
