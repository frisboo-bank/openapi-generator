package mediator

import (
	"frisboo-bank/openapi-generator-service/internal/entities/contracts"
	createentityqueries "frisboo-bank/openapi-generator-service/internal/entities/features/create_entity/queries"
	deleteentityqueries "frisboo-bank/openapi-generator-service/internal/entities/features/delete_entity/queries"
	getentitybyslugqueries "frisboo-bank/openapi-generator-service/internal/entities/features/get_entity_by_slug/queries"
	listentitiesqueries "frisboo-bank/openapi-generator-service/internal/entities/features/list_entities/queries"
	updateentityqueries "frisboo-bank/openapi-generator-service/internal/entities/features/update_entity/queries"
	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/mediator"
	mediatorcontracts "frisboo-bank/openapi-generator-service/pkg/mediator/contracts"
)

func ConfigureEntitiesMediator(
	mediatorInstance mediatorcontracts.Mediator,
	entitiesRepository contracts.EntityRepository,
	logger loggercontracts.Logger,
) error {
	if err := mediator.RegisterRequest(
		mediatorInstance,
		createentityqueries.NewCreateEntityHandler(entitiesRepository, logger),
	); err != nil {
		return err
	}

	if err := mediator.RegisterRequest(
		mediatorInstance,
		deleteentityqueries.NewDeleteEntityHandler(entitiesRepository, logger),
	); err != nil {
		return err
	}

	if err := mediator.RegisterRequest(
		mediatorInstance,
		getentitybyslugqueries.NewGetEntityBySlugHandler(entitiesRepository, logger),
	); err != nil {
		return err
	}

	if err := mediator.RegisterRequest(
		mediatorInstance,
		listentitiesqueries.NewListEntitiesHandler(entitiesRepository, logger),
	); err != nil {
		return err
	}

	if err := mediator.RegisterRequest(
		mediatorInstance,
		updateentityqueries.NewUpdateEntityHandler(entitiesRepository, logger),
	); err != nil {
		return err
	}

	return nil
}
