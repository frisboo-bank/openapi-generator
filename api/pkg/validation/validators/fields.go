package validators

import (
	"fmt"
	"slices"

	ozzo "github.com/go-ozzo/ozzo-validation"
)

func AllowedField(allowed []string) ozzo.Rule {
	return ozzo.By(func(value any) error {
		s, ok := value.(string)
		if !ok {
			return fmt.Errorf("value %T is not a string", value)
		}
		if !slices.Contains(allowed, s) {
			return fmt.Errorf("field %q is not allowed", s)
		}
		return nil
	})
}
