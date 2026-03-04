package sqlite

import (
	"database/sql"
	"time"

	"expo/internal/model"
	"expo/internal/storage/sqlite/sqlc"
)

func mapGetPartRowToModel(row sqlc.GetPartRow) *model.Part {
	return &model.Part{
		ID:        row.ID,
		Name:      row.Name,
		Version:   int(row.Version),
		CreatedAt: row.CreatedAt,
	}
}

func mapListPartRowToModel(row sqlc.ListPartsRow) *model.Part {
	return &model.Part{
		ID:        row.ID,
		Name:      row.Name,
		Version:   int(row.Version),
		CreatedAt: row.CreatedAt,
		UpdatedAt: mapUpdatedAt(row.UpdatedAt),
	}
}

func mapListPartsRowsToModels(rows []sqlc.ListPartsRow) []*model.Part {
	if len(rows) == 0 {
		return nil
	}

	res := make([]*model.Part, 0, len(rows))
	for _, row := range rows {
		res = append(res, mapListPartRowToModel(row))
	}
	return res
}

func mapUpdatedAt(t sql.NullTime) *time.Time {
	if t.Valid {
		return &t.Time
	}
	return nil
}
