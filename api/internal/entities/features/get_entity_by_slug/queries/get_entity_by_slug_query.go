package queries

import (
	entityv1 "frisboo-bank/openapi-generator-service/internal/shared/grpc/gen/entity/v1"

	vendorValidation "github.com/go-ozzo/ozzo-validation"
)

type GetEntityBySlugQuery struct {
	Slug          string
	IncludeHidden bool
}

func NewGetEntityBySlugQuery(request *entityv1.GetEntityBySlugRequest) (*GetEntityBySlugQuery, error) {
	command := &GetEntityBySlugQuery{
		Slug: request.Slug,
		// IncludeHidden: request.IncludeHidden,
	}

	if err := command.Validate(); err != nil {
		return nil, err
	}

	return command, nil
}

func (g *GetEntityBySlugQuery) Validate() error {
	return vendorValidation.Errors{
		"slug": vendorValidation.Validate(g.Slug, vendorValidation.Required, vendorValidation.Length(3, 250)),
	}.Filter()
}
