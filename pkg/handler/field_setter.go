package handler

import (
	"fmt"
	"reflect"
)

// setResolvedField sets a resolved value onto target struct field by index.
//
// The function accepts assignable types directly and also supports conversion
// when reflect determines conversion is safe.
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
