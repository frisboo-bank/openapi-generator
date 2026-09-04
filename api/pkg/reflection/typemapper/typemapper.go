package typemapper

import (
	"fmt"
	"reflect"
)

func GetGenericTypeByT[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}

func GenericInstanceByT[T any]() T {
	typ := GetGenericTypeByT[T]()
	return getInstanceFromType(typ).(T)
}

func GetGenericTypeNameByT[T any]() string {
	t := reflect.TypeOf((*T)(nil)).Elem()
	if t.Kind() != reflect.Ptr {
		return t.Name()
	}

	return fmt.Sprintf("*%s", t.Elem().Name())
}

func getInstanceFromType(typ reflect.Type) any {
	if typ.Kind() == reflect.Pointer {
		return reflect.New(typ.Elem()).Interface()
	}

	return reflect.Zero(typ).Interface()
}
