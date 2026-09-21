package validation

import (
	vendorValidation "github.com/go-ozzo/ozzo-validation"
)

func FlattenErrors(errs vendorValidation.Errors) map[string]string {
	if len(errs) == 0 {
		return nil
	}

	result := make(map[string]string)
	flattenErrorsWithPrefix("", errs, result)

	if len(result) == 0 {
		return nil
	}

	return result
}

func flattenErrorsWithPrefix(prefix string, errs vendorValidation.Errors, result map[string]string) {
	for field, err := range errs {
		if err == nil {
			continue
		}

		key := field
		if prefix != "" {
			key = prefix + "." + key
		}

		if nestedErr, ok := err.(vendorValidation.Errors); ok {
			flattenErrorsWithPrefix(key, nestedErr, result)
		} else {
			result[key] = err.Error()
		}
	}
}
