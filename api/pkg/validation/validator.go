package validation

import (
	"frisboo-bank/openapi-generator-service/pkg/reflection"
	"frisboo-bank/openapi-generator-service/pkg/syserrors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func EnumOneOf[T reflection.EnumValue](options reflection.EnumContainer[T]) validation.RuleFunc {
	opts := make([]T, 0)

	for o := range options.All() {
		opts = append(opts, o)
	}

	return func(value any) error {
		v, ok := value.(T)
		if !ok {
			return syserrors.New("Must be a valid enum")
		}
		if !v.IsValid() {
			return syserrors.MustBeOneOf("", v, opts)
		}
		return nil
	}
}
