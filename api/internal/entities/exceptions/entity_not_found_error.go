package exceptions

import (
	"context"
	applicationerror "frisboo-bank/openapi-generator-service/pkg/application_error"
)

type EntityNotFoundError interface {
	applicationerror.NotFoundError
}

type entityNotFoundError struct {
	applicationerror.NotFoundError
}

func NewEntityNotFoundError(ctx context.Context, entitySlug string, extraDetails map[string]string) EntityNotFoundError {
	return &entityNotFoundError{
		NotFoundError: applicationerror.NewNotFoundError(ctx, "entity", entitySlug, extraDetails),
	}
}
