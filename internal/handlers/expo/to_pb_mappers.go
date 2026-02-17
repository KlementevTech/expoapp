package expo

import (
	"expo/internal/model"
	expov1 "expo/pkg/pb/expo/v1"
)

func mapVersionToPb(version *model.Version) *expov1.GetVersionResponse_Version {
	return &expov1.GetVersionResponse_Version{
		Value: string(*version),
	}
}
