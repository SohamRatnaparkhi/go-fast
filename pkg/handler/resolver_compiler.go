package handler

import (
	"fmt"
	"reflect"
	"strings"

	handlerResolvers "github.com/sohamratnaparkhi/go-fast/pkg/handler/resolvers"
)

// buildResolvers compiles resolver instances for tagged fields in inputType.
//
// It returns both resolver list and the index of the body field (if any).
// The body index is tracked separately because request body is a one-shot reader
// and should be resolved first.
func buildResolvers(inputType reflect.Type) ([]FieldResolver, int, error) {
	resolvers := make([]FieldResolver, 0, inputType.NumField())
	bodyFieldIdx := -1
	hasFormOrFile := false

	for i := 0; i < inputType.NumField(); i++ {
		field := inputType.Field(i)
		tag := normalizedJSONTag(field.Tag.Get("json"))
		if tag == "" || tag == "-" {
			continue
		}

		if !field.IsExported() {
			return nil, -1, fmt.Errorf("field %q is tagged but not exported", field.Name)
		}

		switch {
		case tag == "body":
			if bodyFieldIdx >= 0 {
				return nil, -1, fmt.Errorf("multiple body fields found: %d and %d", bodyFieldIdx, i)
			}
			bodyFieldIdx = i
			resolvers = append(resolvers, NewBodyResolver(i, field.Type))

		case strings.HasPrefix(tag, "header:"):
			name := strings.TrimPrefix(tag, "header:")
			if name == "" {
				return nil, -1, fmt.Errorf("header tag name cannot be empty for field %q", field.Name)
			}
			resolvers = append(resolvers, NewHeaderResolver(i, name, field.Type))

		case strings.HasPrefix(tag, "query:"):
			name := strings.TrimPrefix(tag, "query:")
			if name == "" {
				return nil, -1, fmt.Errorf("query tag name cannot be empty for field %q", field.Name)
			}
			resolvers = append(resolvers, NewQueryResolver(i, name, field.Type))

		case strings.HasPrefix(tag, "path:"):
			name := strings.TrimPrefix(tag, "path:")
			if name == "" {
				return nil, -1, fmt.Errorf("path tag name cannot be empty for field %q", field.Name)
			}
			resolvers = append(resolvers, NewPathVarResolver(i, name, field.Type))

		case strings.HasPrefix(tag, "cookie:"):
			name := strings.TrimPrefix(tag, "cookie:")
			if name == "" {
				return nil, -1, fmt.Errorf("cookie tag name cannot be empty for field %q", field.Name)
			}
			resolvers = append(resolvers, NewCookieResolver(i, name, field.Type))

		case strings.HasPrefix(tag, "form:"):
			name := strings.TrimPrefix(tag, "form:")
			if name == "" {
				return nil, -1, fmt.Errorf("form tag name cannot be empty for field %q", field.Name)
			}
			hasFormOrFile = true
			resolvers = append(resolvers, NewFormResolver(i, name, field.Type))

		case strings.HasPrefix(tag, "file:"):
			name := strings.TrimPrefix(tag, "file:")
			if name == "" {
				return nil, -1, fmt.Errorf("file tag name cannot be empty for field %q", field.Name)
			}
			if field.Type != handlerResolvers.MultipartFileHeaderType {
				return nil, -1, fmt.Errorf("file field %q must be *multipart.FileHeader, got %s", field.Name, field.Type)
			}
			hasFormOrFile = true
			resolvers = append(resolvers, NewFileResolver(i, name))
		}
	}

	if bodyFieldIdx >= 0 && hasFormOrFile {
		return nil, -1, fmt.Errorf("cannot combine body resolver with form/file resolvers: body consumes request body as JSON, form/file consume it as multipart or url-encoded data")
	}

	return resolvers, bodyFieldIdx, nil
}

// normalizedJSONTag returns the first comma-delimited segment of a json tag.
func normalizedJSONTag(tag string) string {
	tag = strings.TrimSpace(tag)
	if tag == "" {
		return ""
	}

	parts := strings.Split(tag, ",")
	return strings.TrimSpace(parts[0])
}

// resolverByFieldIndex returns resolver for a given struct field index.
func resolverByFieldIndex(resolvers []FieldResolver, fieldIndex int) FieldResolver {
	for _, resolver := range resolvers {
		if resolver.FieldIndex() == fieldIndex {
			return resolver
		}
	}
	return nil
}
