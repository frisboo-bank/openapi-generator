package syserrors

func CantBeNilError(name string) error {
	return Newf("%s can't be nil", name)
}

func MustBePositiveError(name string, value any) error {
	return Newf("%s must be >0: got %v", name, value)
}

func CantBeNegativeError(name string, value any) error {
	return Newf("%s cannot be negative: got %v", name, value)
}

func CantBeEmptyError(name string) error {
	return Newf("%s can't be empty", name)
}

func MustBeTrue(name string) error {
	return Newf("%s must be true", name)
}

func MustBeFalse(name string) error {
	return Newf("%s must be false", name)
}

func MustBeOneOf[T comparable](name string, value T, options []T) error {
	return Newf("%s is invalid: got %v, allowed %v", name, value, options)
}
