package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"expo/internal/model"
	"expo/internal/storage/sqlite/sqlc"

	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

type PartStorage struct {
	db      *sql.DB
	queries sqlc.Querier
}

func NewPartStorage(db *sql.DB, queries sqlc.Querier) *PartStorage {
	return &PartStorage{
		db:      db,
		queries: queries,
	}
}

func (r *PartStorage) Create(ctx context.Context, p *model.Part) error {
	err := r.queries.CreatePart(ctx, sqlc.CreatePartParams{
		ID:        p.ID,
		Name:      p.Name,
		Version:   int64(p.Version),
		CreatedAt: p.CreatedAt,
	})
	if err != nil {
		if asErr, ok := errors.AsType[*sqlite.Error](err); ok {
			code := asErr.Code()
			if code == sqlite3.SQLITE_CONSTRAINT_PRIMARYKEY ||
				code == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
				return fmt.Errorf("%w: %w", model.ErrAlreadyExists, asErr)
			}
		}
		return fmt.Errorf("sqlite create (sqlc): %w", err)
	}

	return nil
}

func (r *PartStorage) Get(ctx context.Context, id model.PartID) (*model.Part, error) {
	row, err := r.queries.GetPart(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, fmt.Errorf("sqlite find by id (sqlc): %w", err)
	}

	return mapGetPartRowToModel(row), nil
}

func (r *PartStorage) List(ctx context.Context) ([]*model.Part, error) {
	rows, err := r.queries.ListParts(ctx)
	if err != nil {
		return nil, fmt.Errorf("sqlite get list (sqlx): %w", err)
	}

	return mapListPartsRowsToModels(rows), nil
}

func (r *PartStorage) Update(ctx context.Context, part *model.Part, oldVersion int64) error {
	res, err := r.queries.UpdatePart(ctx, sqlc.UpdatePartParams{
		ID:      part.ID,
		Name:    part.Name,
		Version: int64(part.Version),
		UpdatedAt: sql.NullTime{
			Time: *part.UpdatedAt,
		},
		OldVersion: oldVersion,
	})
	if err != nil {
		return fmt.Errorf("sqlite update: %w", err)
	}

	return r.handleUpdateError(ctx, res, part.ID)
}

func (r *PartStorage) Delete(ctx context.Context, id model.PartID, oldVersion int64, deletedAt time.Time) error {
	res, err := r.queries.DeletePart(ctx, sqlc.DeletePartParams{
		ID:         id,
		OldVersion: oldVersion,
		UpdatedAt: sql.NullTime{
			Time: deletedAt,
		},
		DeletedAt: sql.NullTime{
			Time: deletedAt,
		},
	})
	if err != nil {
		return fmt.Errorf("sqlite update: %w", err)
	}

	return r.handleUpdateError(ctx, res, id)
}

func (r *PartStorage) handleUpdateError(ctx context.Context, res sql.Result, id model.PartID) error {
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}

	if count != 0 {
		return nil
	}

	exists, err := r.queries.ExistsPart(ctx, id)
	if err != nil {
		return fmt.Errorf("sqlite check exists: %w", err)
	}

	if exists {
		return model.ErrOptimisticLockConflict
	}

	return model.ErrNotFound
}
