package handler

import "reflect"
type HandlerMetadata struct {
	FuncValue reflect.Value
	FuncType reflect.Type

	NumInputs int
	NumOutputs int

	InputTypes []reflect.Type
	OutputTypes []reflect.Type

	ReturnsError bool
}