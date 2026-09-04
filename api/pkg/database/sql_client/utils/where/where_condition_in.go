package where

import (
	"fmt"
	"strings"

	"frisboo-bank/openapi-generator-service/pkg/query/models"
)

func inCondition(f *models.QueryFilter) (string, map[string]interface{}, error) {
	if len(f.Values) == 0 {
		return "", nil, fmt.Errorf("IN requires at least 1 value")
	}

	keys := make([]string, len(f.Values))
	vals := make(map[string]interface{}, len(f.Values))

	for i, v := range f.Values {
		key := fmt.Sprintf("%s_%d", f.Field, i)
		keys[i] = ":" + key
		vals[key] = v
	}

	return fmt.Sprintf(
		"%s IN (%s)",
		f.Field,
		strings.Join(keys, ", "),
	), vals, nil
}
