package reflection

import "reflect"

func GetKind(v any) reflect.Kind {
	return reflect.ValueOf(v).Kind()
}

func IsPointer(v any) bool {
	return GetKind(v) == reflect.Pointer
}
