package rest

import (
	"context"

	"expo/internal/model"
	"expo/pkg/openapi"
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
