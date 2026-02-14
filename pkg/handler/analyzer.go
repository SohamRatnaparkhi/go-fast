package handler

import (
	"errors"
	"reflect"
)

// errorInterface is used to detect whether the last handler return type is error.
var errorInterface = reflect.TypeOf((*error)(nil)).Elem()

// Analyze inspects a handler function and returns immutable metadata used by Adapt.
//
// The returned metadata is computed once at startup and reused per request.
func Analyze(fn interface{}) (*HandlerMetadata, error) {
	fnType := reflect.TypeOf(fn)
	if fnType == nil || fnType.Kind() != reflect.Func {
		return nil, errors.New("fn is not a function")
	}

	funcValue := reflect.ValueOf(fn)
	funcType := funcValue.Type()

	numInputs := funcType.NumIn()
	numOutputs := funcType.NumOut()

	inputTypes := make([]reflect.Type, numInputs)
	for i := 0; i < numInputs; i++ {
		inputTypes[i] = funcType.In(i)
	}

	outputTypes := make([]reflect.Type, numOutputs)
	for i := 0; i < numOutputs; i++ {
		outputTypes[i] = funcType.Out(i)
	}

	returnsError := false
	if numOutputs > 0 && outputTypes[numOutputs-1].Implements(errorInterface) {
		returnsError = true
	}

	return &HandlerMetadata{
		FuncValue:    funcValue,
		FuncType:     funcType,
		NumInputs:    numInputs,
		NumOutputs:   numOutputs,
		InputTypes:   inputTypes,
		OutputTypes:  outputTypes,
		ReturnsError: returnsError,
	}, nil
}
