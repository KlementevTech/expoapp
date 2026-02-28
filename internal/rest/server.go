package rest

import (
	"context"

	"expo/internal/gen/openapi"
	"expo/internal/model"
)

type VersionService interface {
	GetVersion() *model.Version
}

type server struct {
	vs VersionService
}

func (s *server) GetVersion(
	_ context.Context,
	_ openapi.GetVersionRequestObject,
) (openapi.GetVersionResponseObject, error) {
	return &openapi.GetVersion200JSONResponse{
		Version: string(*s.vs.GetVersion()),
	}, nil
}
