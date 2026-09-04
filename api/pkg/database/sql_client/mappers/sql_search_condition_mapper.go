package mappers

import (
	"fmt"

	searchconditionEnum "frisboo-bank/openapi-generator-service/pkg/query/models/enums/search_condition"
)

func SearchConditionToSQLOperator(s searchconditionEnum.SearchCondition) (string, error) {
	switch s {
	case searchconditionEnum.SearchConditions.CONTAINS:
		return "ILIKE %%s%", nil
	case searchconditionEnum.SearchConditions.STARTS:
		return "ILIKE %s%", nil
	case searchconditionEnum.SearchConditions.ENDS:
		return "ILIKE %%s", nil
	case searchconditionEnum.SearchConditions.EQUALS:
		return "= %s", nil
	default:
		return "", fmt.Errorf("unsupported search condition: %s", s.String())
	}
}
