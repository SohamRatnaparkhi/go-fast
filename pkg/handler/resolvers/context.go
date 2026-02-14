package resolvers

import "net/http"

// Context carries request-scoped values used by field resolvers.
//
// Params is expected to be populated by a router for path-variable resolution.
// When no router integration is present, Params may be empty.
type Context struct {
	Request *http.Request
	Params  map[string]string
}
