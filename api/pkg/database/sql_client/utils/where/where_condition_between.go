package where

import (
	"fmt"

	"frisboo-bank/openapi-generator-service/pkg/query/models"
)

func betweenCondition(f *models.QueryFilter) (string, map[string]interface{}, error) {
	if len(f.Values) != 2 {
		return "", nil, fmt.Errorf("BETWEEN requires exactly 2 values")
	}

	low := f.Field + "_0"
	high := f.Field + "_1"

	part := fmt.Sprintf("%s BETWEEN :%s AND :%s", f.Field, low, high)

	return part, map[string]interface{}{
		low:  f.Values[0],
		high: f.Values[1],
	}, nil
}
