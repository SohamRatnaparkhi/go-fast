package resolvers

import (
	"fmt"
	"reflect"
	"strconv"
)

func convertStringToType(raw string, fieldType reflect.Type) (reflect.Value, error) {
	if fieldType == nil {
		return reflect.Value{}, fmt.Errorf("field type is nil")
	}

	if raw == "" {
		return reflect.Zero(fieldType), nil
	}

	switch fieldType.Kind() {
	case reflect.String:
		return reflect.ValueOf(raw).Convert(fieldType), nil
	case reflect.Bool:
		v, err := strconv.ParseBool(raw)
		if err != nil {
			return reflect.Value{}, err
		}
		value := reflect.New(fieldType).Elem()
		value.SetBool(v)
		return value, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v, err := strconv.ParseInt(raw, 10, fieldType.Bits())
		if err != nil {
			return reflect.Value{}, err
		}
		value := reflect.New(fieldType).Elem()
		value.SetInt(v)
		return value, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		v, err := strconv.ParseUint(raw, 10, fieldType.Bits())
		if err != nil {
			return reflect.Value{}, err
		}
		value := reflect.New(fieldType).Elem()
		value.SetUint(v)
		return value, nil
	case reflect.Float32, reflect.Float64:
		v, err := strconv.ParseFloat(raw, fieldType.Bits())
		if err != nil {
			return reflect.Value{}, err
		}
		value := reflect.New(fieldType).Elem()
		value.SetFloat(v)
		return value, nil
	case reflect.Ptr:
		innerValue, err := convertStringToType(raw, fieldType.Elem())
		if err != nil {
			return reflect.Value{}, err
		}
		ptrValue := reflect.New(fieldType.Elem())
		ptrValue.Elem().Set(innerValue)
		return ptrValue, nil
	default:
		return reflect.Value{}, fmt.Errorf("unsupported field type %s", fieldType)
	}
}
