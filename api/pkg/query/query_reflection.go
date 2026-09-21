package query

import (
	"fmt"
	"frisboo-bank/openapi-generator-service/pkg/validation"
	"reflect"
	"slices"
	"strings"
)

type AllowedQueryFields struct {
	Filter []string
	Order  []string
	Search []string
}

func AllowedQueryFieldsFor[T any]() AllowedQueryFields {
	var v T
	rt := reflect.TypeOf(v)

	validation.Assert(rt != nil, fmt.Errorf("model can't get the type"))

	if rt.Kind() == reflect.Pointer {
		rt = rt.Elem()
	}

	validation.Assert(rt.Kind() == reflect.Struct, fmt.Errorf("model must be a struct, got %s", rt.Kind()))

	res := AllowedQueryFields{}

	for f := range rt.Fields() {
		tag := f.Tag.Get("query")
		if tag == "" {
			continue
		}

		queryParts := strings.SplitN(tag, ":", 2)
		validation.Assert(len(queryParts) == 2, fmt.Errorf("invalid query tag %q on field %s.%s: expected format \"name:cap1,cap2\"", tag, rt.Name(), f.Name))

		fieldName := strings.TrimSpace(queryParts[0])
		validation.Assert(fieldName != "", fmt.Errorf("invalid query tag %q on field %s.%s: expected format \"name:cap1,cap2\"", tag, rt.Name(), f.Name))

		capabilities := strings.SplitSeq(queryParts[1], ",")

		for c := range capabilities {
			capability := strings.TrimSpace(c)
			switch capability {
			case "filter":
				if slices.Contains(res.Filter, fieldName) {
					panic(fmt.Errorf("field %q already present in filter capabilities for %s.%s", fieldName, rt.Name(), f.Name))
				}
				res.Filter = append(res.Filter, fieldName)
			case "order":
				if slices.Contains(res.Order, fieldName) {
					panic(fmt.Errorf("field %q already present in order capabilities for %s.%s", fieldName, rt.Name(), f.Name))
				}
				res.Order = append(res.Order, fieldName)
			case "search":
				if slices.Contains(res.Search, fieldName) {
					panic(fmt.Errorf("field %q already present in search capabilities for %s.%s", fieldName, rt.Name(), f.Name))
				}
				res.Search = append(res.Search, fieldName)
			default:
				panic(fmt.Errorf(
					"invalid capability %q in query tag on field %s.%s (allowed: filter, order, search)",
					capability, rt.Name(), f.Name,
				))
			}
		}
	}

	return res
}
