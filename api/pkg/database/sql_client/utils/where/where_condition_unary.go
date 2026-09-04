package where

import (
	"fmt"

	"frisboo-bank/openapi-generator-service/pkg/query/models"
)

func unaryCondition(f *models.QueryFilter, operator string) (string, map[string]interface{}, error) {
	if len(f.Values) != 1 {
		return "", nil, fmt.Errorf("operator %s requires exactly 1 value", f.Comparator.String())
	}

	return fmt.Sprintf("%s %s :%s", f.Field, operator, f.Field),
		map[string]interface{}{f.Field: f.Values[0]},
		nil
}
