package queries

import (
	entityv1 "frisboo-bank/openapi-generator-service/internal/shared/grpc/gen/entity/v1"

	vendorValidation "github.com/go-ozzo/ozzo-validation"
)

type CreateEntityQuery struct {
	Slug        string
	Name        string
	Description *string
}

func NewCreateEntityQuery(request *entityv1.CreateEntityRequest) (*CreateEntityQuery, error) {
	// var description *string
	// if request.Description != nil {
	// 	description = &request.Description.Value
	// }

	command := &CreateEntityQuery{
		// Slug:        request.Slug,
		// Name:        request.Name,
		// Description: description,
	}

	if err := command.Validate(); err != nil {
		return nil, err
	}

	return command, nil
}

func (c *CreateEntityQuery) Validate() error {
	return vendorValidation.Errors{
		"slug": vendorValidation.Validate(c.Slug, vendorValidation.Required, vendorValidation.Length(3, 250)),
		"name": vendorValidation.Validate(c.Name, vendorValidation.Required, vendorValidation.Length(3, 250)),
		"description": func() error {
			if c.Description == nil {
				return nil
			}
			return vendorValidation.Validate(*c.Description, vendorValidation.Length(3, 500))
		}(),
	}.Filter()
}
