package services

import (
	"context"

	createentityendpoints "frisboo-bank/openapi-generator-service/internal/entities/features/create_entity/endpoints"
	deleteentityendpoints "frisboo-bank/openapi-generator-service/internal/entities/features/delete_entity/endpoints"
	getentitybyslugendpoints "frisboo-bank/openapi-generator-service/internal/entities/features/get_entity_by_slug/endpoints"
	listentitiesendpoints "frisboo-bank/openapi-generator-service/internal/entities/features/list_entities/endpoints"
	updateentityendpoints "frisboo-bank/openapi-generator-service/internal/entities/features/update_entity/endpoints"
	entityv1 "frisboo-bank/openapi-generator-service/internal/shared/grpc/gen/entity/v1"
	"frisboo-bank/openapi-generator-service/pkg/rpc/rpc_server/registrar"
	"frisboo-bank/openapi-generator-service/pkg/validation"

	"google.golang.org/grpc"
)

var (
	_ entityv1.EntityServiceServer = (*EntityGRPCServerService)(nil)
	_ registrar.Service            = (*EntityGRPCServerService)(nil)
)

type EntityGRPCServerService struct {
	entityv1.UnimplementedEntityServiceServer
	createEntityEndpoint    *createentityendpoints.CreateEntityEndpoint
	deleteEntityEndpoint    *deleteentityendpoints.DeleteEntityEndpoint
	getEntityBySlugEndpoint *getentitybyslugendpoints.GetEntityBySlugEndpoint
	listEntitiesEndpoint    *listentitiesendpoints.ListEntitiesEndpoint
	updateEntityEndpoint    *updateentityendpoints.UpdateEntityEndpoint
}

func NewEntityGRPCServerService(
	createEntityEndpoint *createentityendpoints.CreateEntityEndpoint,
	deleteEntityEndpoint *deleteentityendpoints.DeleteEntityEndpoint,
	getEntityBySlugEndpoint *getentitybyslugendpoints.GetEntityBySlugEndpoint,
	listEntitiesEndpoint *listentitiesendpoints.ListEntitiesEndpoint,
	updateEntityEndpoint *updateentityendpoints.UpdateEntityEndpoint,
) entityv1.EntityServiceServer {
	validation.AssertNotNil("createEntityEndpoint", createEntityEndpoint)
	validation.AssertNotNil("deleteEntityEndpoint", deleteEntityEndpoint)
	validation.AssertNotNil("getEntityBySlugEndpoint", getEntityBySlugEndpoint)
	validation.AssertNotNil("listEntitiesEndpoint", listEntitiesEndpoint)
	validation.AssertNotNil("updateEntityEndpoint", updateEntityEndpoint)

	return &EntityGRPCServerService{
		createEntityEndpoint:    createEntityEndpoint,
		deleteEntityEndpoint:    deleteEntityEndpoint,
		getEntityBySlugEndpoint: getEntityBySlugEndpoint,
		listEntitiesEndpoint:    listEntitiesEndpoint,
		updateEntityEndpoint:    updateEntityEndpoint,
	}
}

func (e *EntityGRPCServerService) CreateEntity(ctx context.Context, req *entityv1.CreateEntityRequest) (*entityv1.CreateEntityResponse, error) {
	return e.createEntityEndpoint.CreateEntity(ctx, req)
}

func (e *EntityGRPCServerService) DeleteEntity(ctx context.Context, req *entityv1.DeleteEntityRequest) (*entityv1.DeleteEntityResponse, error) {
	return e.deleteEntityEndpoint.DeleteEntity(ctx, req)
}

func (e *EntityGRPCServerService) GetEntityBySlug(ctx context.Context, req *entityv1.GetEntityBySlugRequest) (*entityv1.GetEntityBySlugResponse, error) {
	return e.getEntityBySlugEndpoint.GetEntityBySlug(ctx, req)
}

func (e *EntityGRPCServerService) ListEntities(ctx context.Context, req *entityv1.ListEntitiesRequest) (*entityv1.ListEntitiesResponse, error) {
	return e.listEntitiesEndpoint.ListEntities(ctx, req)
}

func (e *EntityGRPCServerService) UpdateEntity(ctx context.Context, req *entityv1.UpdateEntityRequest) (*entityv1.UpdateEntityResponse, error) {
	return e.updateEntityEndpoint.UpdateEntity(ctx, req)
}

func (e *EntityGRPCServerService) RegisterToServer(server *grpc.Server) {
	server.RegisterService(&entityv1.EntityService_ServiceDesc, e)
}

func (e *EntityGRPCServerService) Name() string                         { return "entity.v1.EntityService" }
func (e *EntityGRPCServerService) HealthStatus(ctx context.Context) any { return nil }
func (e *EntityGRPCServerService) Start(ctx context.Context) error      { return nil }
func (e *EntityGRPCServerService) Stop(ctx context.Context) error       { return nil }
