package query

import (
	"fmt"
	"frisboo-bank/openapi-generator-service/pkg/validation/validators"

	ozzo "github.com/go-ozzo/ozzo-validation"
)

func (q Query) Validate(allowed AllowedQueryFields) error {
	for i, f := range q.Filters {
		if err := f.Validate(allowed.Filter); err != nil {
			return fmt.Errorf("filter[%d]: %w", i, err)
		}
	}
	for i, o := range q.Orders {
		if err := o.Validate(allowed.Order); err != nil {
			return fmt.Errorf("order[%d]: %w", i, err)
		}
	}
	if q.Searches != nil {
		if err := q.Searches.Validate(allowed.Search); err != nil {
			return fmt.Errorf("searches: %w", err)
		}
	}
	if q.Pagination != nil {
		if err := q.Pagination.Validate(); err != nil {
			return fmt.Errorf("pagination: %w", err)
		}
	}
	return nil
}

func (f *Filter) Validate(allowed []string) error {
	rules := []*ozzo.FieldRules{
		ozzo.Field(&f.Field, ozzo.Required, validators.AllowedField(allowed)),
		ozzo.Field(&f.Comparator, validators.ValidEnum()),
	}
	if f.Comparator.Options().RequiresValue {
		rules = append(rules, ozzo.Field(&f.Values, ozzo.Required))
	}
	if err := ozzo.ValidateStruct(f, rules...); err != nil {
		return err
	}
	return f.Comparator.ValidateValues(f.Values)
}

func (o *Order) Validate(allowed []string) error {
	return ozzo.ValidateStruct(o,
		ozzo.Field(&o.Field, ozzo.Required, validators.AllowedField(allowed)),
		ozzo.Field(&o.Direction, ozzo.Required, validators.ValidEnum()),
	)
}

func (sg *SearchGroup) Validate(allowed []string) error {
	if sg == nil || sg.IsEmpty() {
		return nil
	}

	if err := ozzo.ValidateStruct(sg, ozzo.Field(&sg.Operator, ozzo.Required, validators.ValidEnum())); err != nil {
		return err
	}

	for i, c := range sg.Conditions {
		if err := c.Validate(allowed); err != nil {
			return fmt.Errorf("condition[%d]: %w", i, err)
		}
	}

	for i, g := range sg.Groups {
		if g == nil {
			continue
		}
		if err := g.Validate(allowed); err != nil {
			return fmt.Errorf("group[%d]: %w", i, err)
		}
	}
	return nil
}

func (s *Search) Validate(allowed []string) error {
	return ozzo.ValidateStruct(s,
		ozzo.Field(&s.Field, ozzo.Required, validators.AllowedField(allowed)),
		ozzo.Field(&s.Value, ozzo.Required),
		ozzo.Field(&s.Operator, validators.ValidEnum()),
	)
}

func (p *Pagination) Validate() error {
	return ozzo.ValidateStruct(p,
		ozzo.Field(&p.Cursor, ozzo.Length(0, 2048)),
		ozzo.Field(&p.Limit, ozzo.Min(1), ozzo.Max(MaxPaginationLimit)),
	)
}
