package devcontainers

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"frisboo-bank/openapi-generator-service/pkg/database/migration"
	"frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	defaultPostgresImage    = "postgres:alpine"
	defaultPostgresDbName   = "testdb"
	defaultPostgresUsername = "postgres"
	defaultPostgresPassword = "postgres"
)

type PostgresTestContainerOptions struct {
	Image                    string
	DBName                   string
	Username                 string
	Password                 string
	InitScripts              []string
	ConfigFile               string
	AdditionalWaitStrategies []wait.Strategy
}

type PostgresTestContainer struct {
	opts      PostgresTestContainerOptions
	container *postgres.PostgresContainer
	dsn       string
	ctx       context.Context
}

func NewPostgresTestContainer(opts PostgresTestContainerOptions) *PostgresTestContainer {
	if opts.Image == "" {
		opts.Image = defaultPostgresImage
	}
	if opts.DBName == "" {
		opts.DBName = defaultPostgresDbName
	}
	if opts.Username == "" {
		opts.Username = defaultPostgresUsername
	}
	if opts.Password == "" {
		opts.Password = defaultPostgresPassword
	}

	return &PostgresTestContainer{
		opts: opts,
		ctx:  context.Background(),
	}
}

func (p *PostgresTestContainer) Run(t *testing.T) {
	t.Helper()

	containerOpts := []testcontainers.ContainerCustomizer{
		postgres.WithDatabase(p.opts.DBName),
		postgres.WithUsername(p.opts.Username),
		postgres.WithPassword(p.opts.Password),
	}

	if len(p.opts.InitScripts) > 0 {
		containerOpts = append(containerOpts, postgres.WithInitScripts(p.opts.InitScripts...))
	}

	if p.opts.ConfigFile != "" {
		containerOpts = append(containerOpts, postgres.WithConfigFile(p.opts.ConfigFile))
	}

	if len(p.opts.AdditionalWaitStrategies) > 0 {
		containerOpts = append(containerOpts, testcontainers.WithWaitStrategy(
			wait.ForAll(p.opts.AdditionalWaitStrategies...),
		))
	} else {
		containerOpts = append(containerOpts, postgres.BasicWaitStrategies())
	}

	container, err := postgres.Run(p.ctx, p.opts.Image, containerOpts...)
	require.NoError(t, err, "failed to start postgres testcontainer")

	p.container = container

	testcontainers.CleanupContainer(t, container)

	dsn, err := container.ConnectionString(p.ctx, "sslmode=disable")
	require.NoError(t, err, "failed to get postgres connection string")
	p.dsn = dsn
}

func (p *PostgresTestContainer) DSN() string {
	if p.dsn == "" {
		panic("PostgresTestContainer.DSN called before Run")
	}
	return p.dsn
}

func (p *PostgresTestContainer) DB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("pgx", p.DSN())
	require.NoError(t, err, "failed to open sql.DB")

	t.Cleanup(func() {
		_ = db.Close()
	})

	return db
}

func (p *PostgresTestContainer) SQLX(t *testing.T) *sqlx.DB {
	t.Helper()

	db, err := sqlx.Connect("pgx", p.DSN())
	require.NoError(t, err, "failed to open sqlx.DB")

	t.Cleanup(func() {
		_ = db.Close()
	})

	return db
}

func (p *PostgresTestContainer) ExecSQL(t *testing.T, sql string) {
	t.Helper()

	db := p.DB(t)
	defer db.Close()

	_, err := db.Exec(sql)
	require.NoError(t, err, "failed to execute SQL")
}

func (p *PostgresTestContainer) MigrateFn(t *testing.T, migrateFn func(ctx context.Context, db *sqlx.DB) error) {
	t.Helper()

	db := p.SQLX(t)
	defer db.Close()

	ctx, cancel := context.WithTimeout(p.ctx, 30*time.Second)
	defer cancel()

	err := migrateFn(ctx, db)
	require.NoError(t, err, "migration failed")
}

func (p *PostgresTestContainer) Migrate(t *testing.T, migrationDir string) {
	t.Helper()

	migrator, err := migration.CreateMigrationForTests(
		"test",
		p.DB(t),
		migrationDir,
		environment.Environments.TESTING,
	)
	require.NoError(t, err, "failed to create test migrator")

	err = migrator.Up(p.ctx, 0)
	require.NoError(t, err, "migrate up failed")
}

func (p *PostgresTestContainer) Container() *postgres.PostgresContainer {
	return p.container
}
