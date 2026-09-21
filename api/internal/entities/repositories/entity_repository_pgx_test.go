package repositories_test

import (
	"frisboo-bank/openapi-generator-service/internal/entities/contracts"
	"frisboo-bank/openapi-generator-service/internal/entities/repositories"
	environmentenum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	"frisboo-bank/openapi-generator-service/pkg/logger"
	"frisboo-bank/openapi-generator-service/pkg/tests/devcontainers"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"
)

type EntityRepositoryPgxTestSuite struct {
	suite.Suite
	pg         *devcontainers.PostgresTestContainer
	db         *sqlx.DB
	repository contracts.EntitySQLRepository
}

func (s *EntityRepositoryPgxTestSuite) SetupSuite() {
	s.pg = devcontainers.NewPostgresTestContainer(devcontainers.PostgresTestContainerOptions{})
	s.pg.Run(s.T())

	s.db = s.pg.SQLX(s.T())
	// require.NoError(s.T(), s.runMigrations())

	s.repository = repositories.NewEntityRepositoryPgx(
		testSQLXAdapter{db: s.db},
		logger.CreateNoopLogger("test", environmentenum.Environments.TESTING),
	)
}
