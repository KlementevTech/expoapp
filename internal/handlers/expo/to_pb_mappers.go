package expo

import (
	expov1 "expo/internal/gen/pb/expo/v1"
	"expo/internal/model"
)

func mapVersionToPb(version *model.Version) *expov1.GetVersionResponse_Version {
	return &expov1.GetVersionResponse_Version{
		Value: string(*version),
	}
}
