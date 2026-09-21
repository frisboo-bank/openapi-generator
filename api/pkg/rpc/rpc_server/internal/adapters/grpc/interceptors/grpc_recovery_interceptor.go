package interceptors

import (
	"context"
	"fmt"

	applicationerror "frisboo-bank/openapi-generator-service/pkg/application_error"
	environmentEnum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"

	"google.golang.org/grpc"
)

type GRPCRecoveryInterceptor struct {
	env    environmentEnum.Environment
	logger loggercontracts.Logger
}

func NewGRPCRecoveryInterceptor(
	logger loggercontracts.Logger,
	env environmentEnum.Environment,
) *GRPCRecoveryInterceptor {
	return &GRPCRecoveryInterceptor{
		logger: logger,
		env:    env,
	}
}

func (i *GRPCRecoveryInterceptor) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (res any, err error) {
		defer func() {
			if r := recover(); r != nil {
				res = nil
				err = i.recoverToAppError(ctx, r, info.FullMethod)
			}
		}()

		return handler(ctx, req)
	}
}

func (i *GRPCRecoveryInterceptor) StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = i.recoverToAppError(ss.Context(), r, info.FullMethod)
			}
		}()

		return handler(srv, ss)
	}
}

func (i *GRPCRecoveryInterceptor) recoverToAppError(ctx context.Context, r any, method string) error {
	var panicErr error
	switch v := r.(type) {
	case error:
		panicErr = v
	case string:
		panicErr = fmt.Errorf("%s", v)
	default:
		panicErr = fmt.Errorf("%v", v)
	}

	// stack := stacktrace.Capture(3)
	// i.logger.Error(
	// 	"panic recovered in gRPC handler",
	// 	"method", method,
	// 	"error", panicErr,
	// 	"stack", stack,
	// )

	return applicationerror.NewInternalError(ctx, fmt.Sprintf("panic: %v\n", panicErr), nil)
}
