package mappers

import (
	"fmt"

	filtercomparisonEnum "frisboo-bank/openapi-generator-service/pkg/query/models/enums/filter_comparison"
)

func FilterComparisonToSQLOperator(f filtercomparisonEnum.FilterComparison) (string, error) {
	switch f {
	case filtercomparisonEnum.FilterComparisons.EQUAL:
		return "=", nil
	case filtercomparisonEnum.FilterComparisons.NOTEQUAL:
		return "<>", nil
	case filtercomparisonEnum.FilterComparisons.GREATER:
		return ">", nil
	case filtercomparisonEnum.FilterComparisons.GREATEROREQUAL:
		return ">=", nil
	case filtercomparisonEnum.FilterComparisons.LOWER:
		return "<", nil
	case filtercomparisonEnum.FilterComparisons.LOWEROREQUAL:
		return "<=", nil
	case filtercomparisonEnum.FilterComparisons.IN:
		return "IN", nil
	case filtercomparisonEnum.FilterComparisons.ISNULL:
		return "IS NULL", nil
	case filtercomparisonEnum.FilterComparisons.ISNOTNULL:
		return "IS NOT NULL", nil
	case filtercomparisonEnum.FilterComparisons.BETWEEN:
		return "BETWEEN", nil
	case filtercomparisonEnum.FilterComparisons.ISTRUE:
		return "IS TRUE", nil
	case filtercomparisonEnum.FilterComparisons.ISNOTTRUE:
		return "IS NOT TRUE", nil
	case filtercomparisonEnum.FilterComparisons.ISFALSE:
		return "IS FALSE", nil
	case filtercomparisonEnum.FilterComparisons.ISNOTFALSE:
		return "IS NOT FALSE", nil
	case filtercomparisonEnum.FilterComparisons.ISUNKNOWN:
		return "IS UNKNOWN", nil
	case filtercomparisonEnum.FilterComparisons.ISEMPTY:
		return "= ''", nil
	case filtercomparisonEnum.FilterComparisons.ISNOTEMPTY:
		return "!= ''", nil
	default:
		return "", fmt.Errorf("unsupported filter comparison: %s", f.String())
	}
}
