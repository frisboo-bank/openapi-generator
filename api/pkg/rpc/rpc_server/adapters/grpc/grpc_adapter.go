package grpc

import (
	"context"
	"fmt"
	"net"
	"time"

	environmentEnum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	loggerContracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/rpc/rpc_server/contracts"
	"frisboo-bank/openapi-generator-service/pkg/rpc/rpc_server/models"
	rpcservertype "frisboo-bank/openapi-generator-service/pkg/rpc/rpc_server/models/enums/rpc_server_type"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

var _ contracts.RPCServer = (*grpcAdapter)(nil)

type grpcAdapter struct {
	address               string
	debug                 bool
	healthCheck           *health.Server
	keepaliveTime         time.Duration
	keepaliveTimeout      time.Duration
	listener              net.Listener
	logger                loggerContracts.Logger
	maxConnectionAge      time.Duration
	maxConnectionAgeGrace time.Duration
	maxConnectionIdle     time.Duration
	name                  string
	server                *grpc.Server
	streamInterceptors    []grpc.StreamServerInterceptor
	unaryInterceptors     []grpc.UnaryServerInterceptor
}

func NewGRPCServer(
	name string,
	cfg *models.RPCServerOptions,
	logger loggerContracts.Logger,
	env environmentEnum.Environment,
) contracts.RPCServer {
	return &grpcAdapter{
		address:               cfg.Address(),
		debug:                 cfg.Debug,
		keepaliveTime:         cfg.KeepAliveTime,
		keepaliveTimeout:      cfg.KeepAliveTimeout,
		logger:                logger,
		maxConnectionAge:      cfg.MaxConnectionAge,
		maxConnectionAgeGrace: cfg.MaxConnectionAgeGrace,
		maxConnectionIdle:     cfg.MaxConnectionIdle,
		name:                  name,
	}
}

func (g *grpcAdapter) SetupDefaultMiddlewares() {
	g.unaryInterceptors = append(g.unaryInterceptors)
}

func (g *grpcAdapter) AddMiddlewares(middlewares ...any) {
	for _, middleware := range middlewares {
		switch v := middleware.(type) {
		case grpc.UnaryServerInterceptor:
			g.unaryInterceptors = append(g.unaryInterceptors, v)
		case grpc.StreamServerInterceptor:
			g.streamInterceptors = append(g.streamInterceptors, v)
		default:
			panic(fmt.Errorf("invalid gRPC middleware type %T", middleware))

		}
	}
}

func (g *grpcAdapter) Start(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	listener, err := net.Listen("tcp", g.address)
	if err != nil {
		return fmt.Errorf("grpc-server %q listen %s: %w", g.name, g.address, err)
	}
	g.listener = listener

	opts := []grpc.ServerOption{
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     g.maxConnectionIdle,
			MaxConnectionAge:      g.maxConnectionAge,
			MaxConnectionAgeGrace: g.maxConnectionAgeGrace,
			Time:                  g.keepaliveTime,
			Timeout:               g.keepaliveTimeout,
		}),
	}

	if len(g.unaryInterceptors) > 0 {
		opts = append(opts, grpc.ChainUnaryInterceptor(g.unaryInterceptors...))
	}

	if len(g.streamInterceptors) > 0 {
		opts = append(opts, grpc.ChainStreamInterceptor(g.streamInterceptors...))
	}

	srv := grpc.NewServer(opts...)
	hcSrv := health.NewServer()
	healthpb.RegisterHealthServer(srv, hcSrv)

	if g.debug {
		reflection.Register(srv)
	}

	g.server = srv
	g.healthCheck = hcSrv

	hcSrv.SetServingStatus(g.name, healthpb.HealthCheckResponse_SERVING)
	g.logger.Infof("grpc server %q listening on %s", g.name, g.address)

	go func() {
		if err := srv.Serve(listener); err != nil {
			g.logger.Errorf("grpc server %q stopped with error: %w", g.name, err)
		}
	}()

	return nil
}

func (g *grpcAdapter) Stop(ctx context.Context) error {
	g.logger.Infof("grpc server %q shutting down", g.name)

	if g.healthCheck != nil {
		g.healthCheck.SetServingStatus(g.name, healthpb.HealthCheckResponse_NOT_SERVING)
	}

	done := make(chan struct{})
	go func() {
		g.server.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		g.server.Stop()
		return ctx.Err()
	}
}

func (g *grpcAdapter) Name() string                      { return g.name }
func (g *grpcAdapter) Type() rpcservertype.RpcServerType { return rpcservertype.RpcServerTypes.GRPC }
func (g *grpcAdapter) Logger() loggerContracts.Logger    { return g.logger }
