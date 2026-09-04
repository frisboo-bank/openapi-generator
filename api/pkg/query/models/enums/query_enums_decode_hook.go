package enums

import (
	"reflect"

	filtercomparisonEnum "frisboo-bank/openapi-generator-service/pkg/query/models/enums/filter_comparison"
	logicaloperatorEnum "frisboo-bank/openapi-generator-service/pkg/query/models/enums/logical_operator"
	orderdirectionEnum "frisboo-bank/openapi-generator-service/pkg/query/models/enums/order_direction"
	searchconditionEnum "frisboo-bank/openapi-generator-service/pkg/query/models/enums/search_condition"

	"github.com/go-viper/mapstructure/v2"
)

func PaginationEnumsDecodeHook() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data any) (any, error) {
		switch t {
		case reflect.TypeFor[filtercomparisonEnum.FilterComparison]():
			return filtercomparisonEnum.ParseFilterComparison(data)
		case reflect.TypeFor[logicaloperatorEnum.LogicalOperator]():
			return logicaloperatorEnum.ParseLogicalOperator(data)
		case reflect.TypeFor[orderdirectionEnum.OrderDirection]():
			return orderdirectionEnum.ParseOrderDirection(data)
		case reflect.TypeFor[searchconditionEnum.SearchCondition]():
			return searchconditionEnum.ParseSearchCondition(data)
		}

		return data, nil
	}
}
