//go:generate goenums -f -c ./migration_type.go
package migrationtype

type migrationType int8

const (
	unknown migrationType = iota // invalid
	goose
)
