package queries

import (
	entityv1 "frisboo-bank/openapi-generator-service/internal/shared/grpc/gen/entity/v1"

	vendorValidation "github.com/go-ozzo/ozzo-validation"
)

type DeleteEntityQuery struct {
	Slug        string
	VersionLock int64
}

func NewDeleteEntityQuery(request *entityv1.DeleteEntityRequest) (*DeleteEntityQuery, error) {
	command := &DeleteEntityQuery{
		// Slug:        request.Slug,
		// VersionLock: request.VersionLock,
	}
	if err := command.Validate(); err != nil {
		return nil, err
	}
	return command, nil
}

func (d *DeleteEntityQuery) Validate() error {
	return vendorValidation.Errors{
		"slug":        vendorValidation.Validate(d.Slug, vendorValidation.Required, vendorValidation.Length(3, 250)),
		"versionLock": vendorValidation.Validate(d.VersionLock, vendorValidation.Required, vendorValidation.Min(1)),
	}.Filter()
}
