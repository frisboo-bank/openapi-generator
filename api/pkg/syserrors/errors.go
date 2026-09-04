package syserrors

// CantBeNilError indicates a field was unexpectedly nil.
func CantBeNilError(name string) error {
	return Newf("%s can't be nil", name)
}

// MustBePositiveError indicates a field must be > 0.
func MustBePositiveError(name string, value any) error {
	return Newf("%s must be >0: got %v", name, value)
}

// CantBeNegativeError indicates a field cannot be negative.
func CantBeNegativeError(name string, value any) error {
	return Newf("%s cannot be negative: got %v", name, value)
}

// CantBeEmptyError indicates a string (or collection) cannot be empty.
func CantBeEmptyError(name string) error {
	return Newf("%s can't be empty", name)
}

// MustBeTrue indicates a boolean must be true.
func MustBeTrue(name string) error {
	return Newf("%s must be true", name)
}

// MustBeFalse indicates a boolean must be false.
func MustBeFalse(name string) error {
	return Newf("%s must be false", name)
}

// MustBeOneOf indicates a value must match one of allowed options.
func MustBeOneOf[T comparable](name string, value T, options []T) error {
	return Newf("%s is invalid: got %v, allowed %v", name, value, options)
}
