package queries

import (
	entityv1 "frisboo-bank/openapi-generator-service/internal/shared/grpc/gen/entity/v1"
	"frisboo-bank/openapi-generator-service/pkg/validation"

	vendorValidation "github.com/go-ozzo/ozzo-validation"
)

type UpdateEntityQuery struct {
	Slug        string
	Name        *string
	Description *string
	VersionLock int64
}

func NewUpdateEntityQuery(request *entityv1.UpdateEntityRequest) (*UpdateEntityQuery, error) {
	validation.AssertNotNil("request", request)

	command := &UpdateEntityQuery{
		Slug:        request.Slug,
		VersionLock: request.VersionLock,
	}

	// if request.Body != nil {
	// 	if request.Body.Payload != nil {
	// 		command.Name = request.Body.Payload.Name
	// 		command.Description = request.Body.Payload.Description
	// 	}
	// }

	return command, nil
}

func (u *UpdateEntityQuery) Validate() error {
	// err := vendorValidation.Errors{
	// 	"updateMask": vendorValidation.Validate(u.UpdateMask, vendorValidation.Required),
	// }.Filter()
	// if err != nil {
	// 	return err
	// }
	//
	return vendorValidation.Errors{
		"slug": vendorValidation.Validate(u.Slug, vendorValidation.Required, vendorValidation.Length(3, 250)),
		"name": func() error {
			if u.Name == nil {
				return nil
			}
			return vendorValidation.Validate(*u.Name, vendorValidation.Required, vendorValidation.Length(3, 250))
		}(),
		"description": func() error {
			if u.Description == nil {
				return nil
			}
			return vendorValidation.Validate(*u.Description, vendorValidation.Length(3, 500))
		}(),
		"versionLock": vendorValidation.Validate(u.VersionLock, vendorValidation.Required, vendorValidation.Min(1)),
	}.Filter()
}
