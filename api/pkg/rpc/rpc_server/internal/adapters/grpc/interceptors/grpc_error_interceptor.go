package interceptors

import (
	"context"

	applicationerror "frisboo-bank/openapi-generator-service/pkg/application_error"
	environmentenum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	grpcerror "frisboo-bank/openapi-generator-service/pkg/rpc/rpc_server/internal/adapters/grpc/grpc_error"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type GRPCErrorInterceptor struct {
	env    environmentenum.Environment
	logger loggercontracts.Logger
}

func NewGRPCErrorInterceptor(
	logger loggercontracts.Logger,
	env environmentenum.Environment,
) *GRPCErrorInterceptor {
	return &GRPCErrorInterceptor{
		logger: logger,
		env:    env,
	}
}

func (i *GRPCErrorInterceptor) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (res any, err error) {
		res, err = handler(ctx, req)
		if err == nil {
			return res, nil
		}
		return nil, i.mapError(ctx, err).Err()
	}
}

func (i *GRPCErrorInterceptor) StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		err = handler(srv, ss)
		if err == nil {
			return nil
		}
		return i.mapError(ss.Context(), err).Err()
	}
}

func (i *GRPCErrorInterceptor) mapError(ctx context.Context, err error) *status.Status {
	if gs, ok := status.FromError(err); ok {
		return gs
	}

	appErr, ok := applicationerror.IsAppError(err)
	if !ok {
		appErr = applicationerror.NewInternalErrorWrap(ctx, err, "internal server error", nil)
	}
	return grpcerror.MapAppErrorToGRPC(i.env, appErr)
}
