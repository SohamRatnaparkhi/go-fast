package handler

import "reflect"

// HandlerMetadata contains precomputed information about a user handler
// function. It is created once by Analyze and reused by Adapt for each request.
type HandlerMetadata struct {
	FuncValue reflect.Value
	FuncType  reflect.Type

	NumInputs  int
	NumOutputs int

	InputTypes  []reflect.Type
	OutputTypes []reflect.Type

	ReturnsError bool
}
