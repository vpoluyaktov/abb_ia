package utils

import (
	"errors"
	"reflect"
)

// DeepCopy creates a deep copy of the given source object.
func DeepCopy(src interface{}) (interface{}, error) {
	srcValue := reflect.ValueOf(src)
	// Create a new instance of the source struct
	destValue := reflect.New(srcValue.Type().Elem())

	// Copy each field recursively
	err := copyFields(srcValue, destValue)
	if err != nil {
		return nil, err
	}

	return destValue.Elem().Interface(), nil
}

func copyFields(srcValue, destValue reflect.Value) error {
	switch srcValue.Kind() {
	case reflect.Ptr:
		// If the field is a pointer, create a new copy of the pointed value
		if !srcValue.IsNil() {
			newPointer := reflect.New(srcValue.Type().Elem())
			err := copyFields(srcValue.Elem(), newPointer.Elem())
			if err != nil {
				return err
			}
			destValue.Set(newPointer)
		}
	case reflect.Struct:
		// If the field is a struct, copy each field inside it
		for i := 0; i < srcValue.NumField(); i++ {
			srcField := srcValue.Field(i)
			destField := destValue.Field(i)

			if srcField.CanSet() {
				err := copyFields(srcField, destField)
				if err != nil {
					return err
				}
			}
		}
	default:
		// For other types, simply copy the value
		if destValue.Type() == srcValue.Type() {
			destValue.Set(srcValue)
		} else {
			return errors.New("unsupported type")
		}
	}

	return nil
}
