package sqlite_test

import (
	"context"
	"database/sql"
	"testing"

	"expo/internal/storage"
	"expo/internal/storage/sqlite"
	"expo/internal/storage/sqlite/sqlc"
	"expo/migrations/sqlite/schema"

	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/suite"
	_ "modernc.org/sqlite"
)

type SQLiteTestSuite struct {
	storage.PartStorageSuite

	db      *sql.DB
	queries sqlc.Querier
}

func (s *SQLiteTestSuite) SetupTest() {
	db, err := sql.Open("sqlite", "file::memory:?cache=shared")
	if err != nil {
		s.T().Fatalf("failed to open sqlite: %v", err)
	}

	goose.SetBaseFS(schema.FS)

	if err = goose.SetDialect("sqlite3"); err != nil {
		s.T().Fatalf("failed to set goose dialect: %v", err)
	}

	if err = goose.Up(db, "."); err != nil {
		s.T().Fatalf("goose up: %v", err)
	}

	s.queries = sqlc.New(db)

	s.Repo = sqlite.NewPartStorage(
		db,
		s.queries,
	)
}

func (s *SQLiteTestSuite) TearDownTest() {
	_ = s.queries.CleanParts(context.Background())

	if s.db != nil {
		_ = s.db.Close()
	}
}

func TestSQLiteSuite(t *testing.T) {
	suite.Run(t, new(SQLiteTestSuite))
}
