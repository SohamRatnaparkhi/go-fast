package handler

import (
	"errors"
	"reflect"
)

var errorInterface = reflect.TypeOf((*error)(nil)).Elem()

func Analyze(fn interface{}) (*HandlerMetadata, error) {
	if reflect.TypeOf(fn).Kind() != reflect.Func {
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