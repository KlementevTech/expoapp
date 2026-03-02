package rest

import (
	"context"

	"expo/internal/gen/openapi"
	"expo/internal/model"
)

type VersionService interface {
	GetVersion() *model.Version
}

type strictServerImpl struct {
	versionSvc VersionService
}

func (s *strictServerImpl) GetVersion(
	_ context.Context,
	_ openapi.GetVersionRequestObject,
) (openapi.GetVersionResponseObject, error) {
	return &openapi.GetVersion200JSONResponse{
		Version: string(*s.versionSvc.GetVersion()),
	}, nil
}
