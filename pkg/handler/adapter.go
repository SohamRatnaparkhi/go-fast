package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

func Adapt(fn interface{}) (http.HandlerFunc, error) {
	meta, err := Analyze(fn)
	if err != nil {
		return nil, err
	}

	if meta.NumInputs != 1 {
		return nil, fmt.Errorf("handler must have exactly 1 input, got %d", meta.NumInputs)
	}

	inputType := meta.InputTypes[0]
	if inputType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("handler input must be a struct, got %s", inputType.Kind())
	}

	resolvers, bodyFieldIdx, err := buildResolvers(inputType)
	if err != nil {
		return nil, err
	}

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := &Context{Request: r, Response: w, Params: map[string]string{}}
		paramValue := reflect.New(inputType).Elem()

		if bodyFieldIdx >= 0 {
			bodyResolver := resolverByFieldIndex(resolvers, bodyFieldIdx)
			if bodyResolver == nil {
				writeError(w, http.StatusInternalServerError, "body resolver missing")
				return
			}

			val, resolveErr := bodyResolver.Resolve(ctx)
			if resolveErr != nil {
				writeError(w, http.StatusBadRequest, resolveErr.Error())
				return
			}

			if setErr := setResolvedField(paramValue, bodyResolver.FieldIndex(), val); setErr != nil {
				writeError(w, http.StatusBadRequest, setErr.Error())
				return
			}
		}

		for _, resolver := range resolvers {
			if resolver.FieldIndex() == bodyFieldIdx {
				continue
			}

			val, resolveErr := resolver.Resolve(ctx)
			if resolveErr != nil {
				writeError(w, http.StatusBadRequest, resolveErr.Error())
				return
			}

			if setErr := setResolvedField(paramValue, resolver.FieldIndex(), val); setErr != nil {
				writeError(w, http.StatusBadRequest, setErr.Error())
				return
			}
		}

		results := meta.FuncValue.Call([]reflect.Value{paramValue})

		if meta.ReturnsError {
			errVal := results[len(results)-1]
			if !errVal.IsNil() {
				writeError(w, http.StatusInternalServerError, errVal.Interface().(error).Error())
				return
			}
		}

		nonErrorResults := meta.NumOutputs
		if meta.ReturnsError {
			nonErrorResults--
		}

		if nonErrorResults == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if encodeErr := json.NewEncoder(w).Encode(results[0].Interface()); encodeErr != nil {
			writeError(w, http.StatusInternalServerError, encodeErr.Error())
		}
	}, nil
}

func buildResolvers(inputType reflect.Type) ([]FieldResolver, int, error) {
	resolvers := make([]FieldResolver, 0, inputType.NumField())
	bodyFieldIdx := -1

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
		}
	}

	return resolvers, bodyFieldIdx, nil
}

func normalizedJSONTag(tag string) string {
	tag = strings.TrimSpace(tag)
	if tag == "" {
		return ""
	}

	parts := strings.Split(tag, ",")
	return strings.TrimSpace(parts[0])
}

func resolverByFieldIndex(resolvers []FieldResolver, fieldIndex int) FieldResolver {
	for _, resolver := range resolvers {
		if resolver.FieldIndex() == fieldIndex {
			return resolver
		}
	}
	return nil
}

func setResolvedField(target reflect.Value, fieldIndex int, resolvedValue reflect.Value) error {
	field := target.Field(fieldIndex)
	if !field.CanSet() {
		return fmt.Errorf("cannot set field at index %d", fieldIndex)
	}

	if resolvedValue.Type().AssignableTo(field.Type()) {
		field.Set(resolvedValue)
		return nil
	}

	if resolvedValue.Type().ConvertibleTo(field.Type()) {
		field.Set(resolvedValue.Convert(field.Type()))
		return nil
	}

	return fmt.Errorf("resolved type %s cannot be assigned to field type %s", resolvedValue.Type(), field.Type())
}

func writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
