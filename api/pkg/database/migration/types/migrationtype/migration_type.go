package migrationtype

type (
	migrationType int8
)

const (
	unknown migrationType = iota // invalid
	goose
)
