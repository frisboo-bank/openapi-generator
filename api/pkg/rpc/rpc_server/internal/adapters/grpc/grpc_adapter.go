package grpc

import (
	"context"
	"fmt"
	"net"
	"time"

	"frisboo-bank/openapi-generator-service/pkg/rpc/rpc_server/contracts"
	"frisboo-bank/openapi-generator-service/pkg/rpc/rpc_server/models"
	"frisboo-bank/openapi-generator-service/pkg/rpc/rpc_server/registrar"
	"frisboo-bank/openapi-generator-service/pkg/validation"

	environmentenum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"

	"frisboo-bank/openapi-generator-service/pkg/rpc/rpc_server/internal/adapters/grpc/interceptors"
	rpcservertype "frisboo-bank/openapi-generator-service/pkg/rpc/rpc_server/models/enums/rpc_server_type"

	grpcctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

var _ contracts.RPCServerAdapter = (*grpcAdapter)(nil)

type grpcAdapter struct {
	address               string
	debug                 bool
	env                   environmentenum.Environment
	logger                loggercontracts.Logger
	maxConnectionAge      time.Duration
	maxConnectionAgeGrace time.Duration
	maxConnectionIdle     time.Duration
	name                  string
	server                *grpc.Server
	serviceManager        *registrar.ServiceManager
	streamInterceptors    []grpc.StreamServerInterceptor
	time                  time.Duration
	timeout               time.Duration
	unaryInterceptors     []grpc.UnaryServerInterceptor
}

func NewGRPCServer(
	name string,
	cfg *models.RPCServerOptions,
	logger loggercontracts.Logger,
	env environmentenum.Environment,
) contracts.RPCServerAdapter {
	validation.AssertNotEmpty("name", name)
	validation.AssertNotNil("cfg", cfg)
	validation.AssertNotNil("logger", logger)
	validation.AssertValidEnum("env", env)

	serviceManager := registrar.NewServiceManager(logger)

	return &grpcAdapter{
		address:               cfg.Address(),
		debug:                 cfg.Debug,
		env:                   env,
		logger:                logger,
		maxConnectionAge:      cfg.MaxConnectionAge,
		maxConnectionAgeGrace: cfg.MaxConnectionAgeGrace,
		maxConnectionIdle:     cfg.MaxConnectionIdle,
		name:                  name,
		serviceManager:        serviceManager,
		time:                  cfg.KeepAliveTime,
		timeout:               cfg.KeepAliveTimeout,
	}
}

func (g *grpcAdapter) SetupDefaultMiddlewares() {
	recoveryInterceptor := interceptors.NewGRPCRecoveryInterceptor(g.logger, g.env)
	errorInterceptor := interceptors.NewGRPCErrorInterceptor(g.logger, g.env)

	g.unaryInterceptors = append(
		g.unaryInterceptors,
		errorInterceptor.UnaryServerInterceptor(),
		grpcctxtags.UnaryServerInterceptor(),
		recoveryInterceptor.UnaryServerInterceptor(),
	)
	g.streamInterceptors = append(
		g.streamInterceptors,
		errorInterceptor.StreamServerInterceptor(),
		recoveryInterceptor.StreamServerInterceptor(),
	)
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

func (g *grpcAdapter) Start(ctx context.Context) (err error) {
	if err := ctx.Err(); err != nil {
		return err
	}

	listener, err := net.Listen("tcp", g.address)
	if err != nil {
		return fmt.Errorf("grpc-server %q listen %s: %w", g.name, g.address, err)
	}

	opts := []grpc.ServerOption{
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     g.maxConnectionIdle,
			MaxConnectionAge:      g.maxConnectionAge,
			MaxConnectionAgeGrace: g.maxConnectionAgeGrace,
			Time:                  g.time,
			Timeout:               g.timeout,
		}),
	}
	if len(g.unaryInterceptors) > 0 {
		opts = append(opts, grpc.ChainUnaryInterceptor(g.unaryInterceptors...))
	}
	if len(g.streamInterceptors) > 0 {
		opts = append(opts, grpc.ChainStreamInterceptor(g.streamInterceptors...))
	}

	grpcServer := grpc.NewServer(opts...)
	if g.debug {
		reflection.Register(grpcServer)
	}

	for _, service := range g.serviceManager.Services() {
		if reg, ok := service.(interface{ RegisterToServer(*grpc.Server) }); ok {
			reg.RegisterToServer(grpcServer)
		}
	}

	g.server = grpcServer

	g.logger.Infof("grpc server %q listening on %s", g.name, g.address)

	go func() {
		if err := g.server.Serve(listener); err != nil {
			g.logger.Errorf("grpc server %q stopped with error: %w", g.name, err)
		}
	}()

	return nil
}

func (g *grpcAdapter) Stop(ctx context.Context) error {
	if g.server == nil {
		return nil
	}

	g.logger.Infof("grpc server shutting down")

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

func (g *grpcAdapter) ListServices() []any {
	return []any{}
}

func (g *grpcAdapter) Name() string                              { return g.name }
func (g *grpcAdapter) Type() rpcservertype.RpcServerType         { return rpcservertype.RpcServerTypes.GRPC }
func (g *grpcAdapter) ServiceManager() *registrar.ServiceManager { return g.serviceManager }
func (g *grpcAdapter) Logger() loggercontracts.Logger            { return g.logger }
