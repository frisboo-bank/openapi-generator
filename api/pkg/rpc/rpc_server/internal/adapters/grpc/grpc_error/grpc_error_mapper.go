package grpcerror

import (
	"sort"

	applicationerror "frisboo-bank/openapi-generator-service/pkg/application_error"
	applicationerrorcontracts "frisboo-bank/openapi-generator-service/pkg/application_error/contracts"
	environmentenum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/protoadapt"
)

func MapAppErrorToGRPC(
	env environmentenum.Environment,
	err applicationerrorcontracts.AppError,
) *status.Status {
	var st *status.Status
	var details protoadapt.MessageV1

	switch err.Kind() {
	case applicationerrorcontracts.KindNotFound:
		st = status.New(codes.NotFound, err.Message())
		details = mapNotFoundErrorDetails(err)
	case applicationerrorcontracts.KindConflict:
		st = status.New(codes.AlreadyExists, err.Message())
		details = mapAlreadyExistsErrorDetails(err)
	case applicationerrorcontracts.KindVersionMismatch:
		st = status.New(codes.Aborted, err.Message())
		details = mapVersionMismatchErrorDetails(err)
	case applicationerrorcontracts.KindValidationFailed:
		st = status.New(codes.InvalidArgument, err.Message())
		details = mapValidationErrorDetails(err)
	case applicationerrorcontracts.KindBusinessRule:
		st = status.New(codes.FailedPrecondition, err.Message())
	case applicationerrorcontracts.KindUnauthorized:
		st = status.New(codes.Unauthenticated, err.Message())
	case applicationerrorcontracts.KindForbidden:
		st = status.New(codes.PermissionDenied, err.Message())
	case applicationerrorcontracts.KindRateLimited:
		st = status.New(codes.ResourceExhausted, err.Message())
	case applicationerrorcontracts.KindInternal:
		st = status.New(codes.Internal, "internal server error")
		details = mapInternalErrorDetails(err)
	default:
		st = status.New(codes.Unknown, err.Message())
	}

	if details == nil {
		return st
	}

	newSt, attachErr := st.WithDetails(details)
	if attachErr != nil {
		return st
	}
	return newSt
}

func mapNotFoundErrorDetails(err applicationerrorcontracts.AppError) protoadapt.MessageV1 {
	nErr, ok := err.(applicationerror.NotFoundError)
	if !ok {
		return nil
	}

	details := nErr.Details()
	if len(details) == 0 {
		return nil
	}

	return &errdetails.ResourceInfo{
		ResourceType: details["resourceType"],
		ResourceName: details["identifier"],
		Description:  err.Message(),
	}
}

func mapVersionMismatchErrorDetails(err applicationerrorcontracts.AppError) protoadapt.MessageV1 {
	vErr, ok := err.(applicationerror.VersionMismatchError)
	if !ok {
		return nil
	}

	details := vErr.Details()
	if len(details) == 0 {
		return nil
	}

	return &errdetails.ErrorInfo{
		Reason:   "VERSION_MISMATCH",
		Domain:   details["resourceType"],
		Metadata: details,
	}
}

func mapAlreadyExistsErrorDetails(err applicationerrorcontracts.AppError) protoadapt.MessageV1 {
	vErr, ok := err.(applicationerror.ConflictError)
	if !ok {
		return nil
	}

	details := vErr.Details()
	if len(details) == 0 {
		return nil
	}

	return &errdetails.ResourceInfo{
		ResourceType: details["resourceType"],
		ResourceName: details["identifier"],
		Description:  err.Message(),
	}
}

func mapValidationErrorDetails(err applicationerrorcontracts.AppError) protoadapt.MessageV1 {
	vErr, ok := err.(applicationerror.ValidationFailedError)
	if !ok {
		return nil
	}

	errorDetails := vErr.Details()
	if len(errorDetails) == 0 {
		return nil
	}

	br := &errdetails.BadRequest{}
	for k, v := range errorDetails {
		br.FieldViolations = append(br.FieldViolations, &errdetails.BadRequest_FieldViolation{
			Field:       k,
			Description: v,
		})
	}

	sort.Slice(br.FieldViolations, func(i, j int) bool {
		return br.FieldViolations[i].Field < br.FieldViolations[j].Field
	})

	return br
}

func mapInternalErrorDetails(err applicationerrorcontracts.AppError) protoadapt.MessageV1 {
	sErr, ok := err.(applicationerror.InternalError)
	if !ok {
		return nil
	}

	details := sErr.Details()

	domain := details["service"]
	if domain == "" {
		domain = "unknown"
	}

	return &errdetails.ErrorInfo{
		Reason:   err.Error(),
		Domain:   domain,
		Metadata: details,
	}
}
