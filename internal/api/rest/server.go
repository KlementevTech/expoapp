package rest

import (
	"context"

	"expo/gen/openapi"
)

type strictServerImpl struct{}

func (s *strictServerImpl) GetVersion(
	_ context.Context,
	_ openapi.GetVersionRequestObject,
) (openapi.GetVersionResponseObject, error) {
	return &openapi.GetVersion200JSONResponse{
		Version: "dev",
	}, nil
}
