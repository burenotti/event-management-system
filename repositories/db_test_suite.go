package repositories

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type DBTestSuite struct {
	suite.Suite
	db            *sqlx.DB
	m             *migrate.Migrate
	DbDsn         string
	MigrationsDir string
	MigrationsDsn string
}

func NewDBTestSuite(dbDsn, migrateDsn, migrateDir string) *DBTestSuite {
	return &DBTestSuite{
		Suite:         suite.Suite{},
		DbDsn:         dbDsn,
		MigrationsDir: migrateDir,
		MigrationsDsn: migrateDsn,
	}
}

func (s *DBTestSuite) SetupSuite() {
	var err error

	s.db, err = sqlx.Connect("pgx", s.DbDsn)
	require.NoError(s.T(), err, "failed to connect to database")

	s.m, err = migrate.New(s.MigrationsDir, s.MigrationsDsn)
	require.NoError(s.T(), err, "failed to open migrations")

	err = s.m.Up()
	require.NoError(s.T(), err, "failed to migrate database")
}

func DBTestSuiteFromEnv() *DBTestSuite {
	return NewDBTestSuite(
		os.Getenv("DB_DSN"),
		os.Getenv("MIGRATE_DSN"),
		os.Getenv("MIGRATE_DIR"),
	)
}

func (s *DBTestSuite) TearDownSuite() {
	_ = s.m.Down()
	_ = s.db.Close()
}
