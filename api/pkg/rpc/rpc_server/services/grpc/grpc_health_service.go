package grpc

import (
	"context"
	"fmt"

	loggerContracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/rpc/rpc_server/registrar"
	"frisboo-bank/openapi-generator-service/pkg/validation"

	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	"google.golang.org/grpc/health"
)

var _ registrar.Service = (*grpcHealthService)(nil)

type grpcHealthService struct {
	logger      loggerContracts.Logger
	server      *health.Server
	serviceName string
}

func NewGRPCHealthService(serviceName string, logger loggerContracts.Logger) registrar.Service {
	validation.AssertNotEmpty("serviceName", serviceName)
	validation.AssertNotNil("logger", logger)

	return &grpcHealthService{
		logger:      logger,
		server:      health.NewServer(),
		serviceName: serviceName,
	}
}

func (g *grpcHealthService) Name() string {
	return fmt.Sprintf("health.%s", g.serviceName)
}

func (g *grpcHealthService) RegisterToServer(server *grpc.Server) {
	healthpb.RegisterHealthServer(server, g.server)
}

func (g *grpcHealthService) Start(ctx context.Context) error {
	g.server.SetServingStatus(g.serviceName, healthpb.HealthCheckResponse_SERVING)
	g.logger.Infof("health service for %q set to SERVING", g.serviceName)
	return nil
}

func (g *grpcHealthService) Stop(ctx context.Context) error {
	g.server.SetServingStatus(g.serviceName, healthpb.HealthCheckResponse_NOT_SERVING)
	g.logger.Infof("health service for %q set to NOT_SERVING", g.serviceName)
	return nil
}

func (g *grpcHealthService) HealthStatus(ctx context.Context) interface{} {
	// return g.server.GetServingStatus(g.serviceName)
	return nil
}
