package migration

import (
	"database/sql"

	"frisboo-bank/openapi-generator-service/pkg/database/migration/contracts"
	migrationinternal "frisboo-bank/openapi-generator-service/pkg/database/migration/internal"
	environmentenum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
)

func CreateMigrationForTests(
	name string,
	db *sql.DB,
	migrationDir string,
	env environmentenum.Environment,
) (contracts.Migration, error) {
	return migrationinternal.CreateMigrationForTests(name, db, migrationDir, env)
}
