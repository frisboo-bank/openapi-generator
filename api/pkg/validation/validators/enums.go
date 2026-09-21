package validators

import (
	"fmt"

	ozzo "github.com/go-ozzo/ozzo-validation"
)

type enumInterface interface {
	IsValid() bool
	String() string
}

func ValidEnum(allowed ...enumInterface) ozzo.Rule {
	return ozzo.By(func(value any) error {
		enum, ok := value.(enumInterface)
		if !ok {
			return fmt.Errorf("value %T does not implement IsValid()", value)
		}
		if !enum.IsValid() {
			return fmt.Errorf("invalid enum value: %s", enum.String())
		}
		if len(allowed) == 0 {
			return nil
		}
		for _, item := range allowed {
			if item.String() == enum.String() {
				return nil
			}
		}
		return fmt.Errorf("enum value %s is not allowed (allowed: %v)", enum.String(), allowed)
	})
}
