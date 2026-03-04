package grpc

import (
	"context"
	"time"

	"expo/internal/model"
)

//go:generate mockgen -source=interfaces.go -destination=mocks/mocks_mock.go -package=mocks

type PartRepository interface {
	Create(ctx context.Context, part *model.Part) error
	Get(ctx context.Context, id model.PartID) (*model.Part, error)
	List(ctx context.Context) ([]*model.Part, error)
	Update(ctx context.Context, part *model.Part, oldVersion int64) error
	Delete(ctx context.Context, id model.PartID, oldVersion int64, deletedAt time.Time) error
}
