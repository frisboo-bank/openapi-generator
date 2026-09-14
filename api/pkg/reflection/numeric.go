package reflection

import "reflect"

// Numeric includes signed integers, unsigned integers, and floats.
type Numeric interface {
	SignedInteger | UnsignedInteger | Float
}

type SignedInteger interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type UnsignedInteger interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

type Integer interface {
	SignedInteger | UnsignedInteger
}

type Float interface {
	~float32 | ~float64
}

// IsNil detects both untyped nil and typed-nil interface cases.
func IsNil(v any) bool {
	if v == nil {
		return true
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Pointer, reflect.Interface, reflect.Slice, reflect.Map, reflect.Func, reflect.Chan:
		return rv.IsNil()
	default:
		return false
	}
}

// IsPossiblyNil reports whether the dynamic kind could represent a nil value.
func IsPossiblyNil(v any) bool {
	if v == nil {
		return true
	}
	k := reflect.TypeOf(v).Kind()
	switch k {
	case reflect.Pointer, reflect.Interface, reflect.Slice, reflect.Map, reflect.Func, reflect.Chan:
		return true
	default:
		return false
	}
}
