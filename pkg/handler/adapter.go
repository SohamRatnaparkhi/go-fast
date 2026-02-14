package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
)

// Adapt validates and compiles a user handler function into an http.HandlerFunc.
//
// The returned closure reuses precomputed metadata and field resolvers so that
// expensive reflection analysis happens once at startup, not on every request.
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
		ctx := &Context{Request: r, Params: map[string]string{}}
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
