package expo_test

import (
	"testing"

	"expo/internal/api/grpc/expo"
	"expo/internal/api/grpc/expo/mocks"
	expov1 "expo/internal/gen/pb/expo/v1"
	"expo/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const devVersion = "dev"

func TestService_GetVersion(t *testing.T) {
	tests := []struct {
		name    string
		prepare func(mock *mocks.MockVersionService)
		want    *expov1.GetVersionResponse
	}{
		{
			name: "success",
			prepare: func(versions *mocks.MockVersionService) {
				versions.EXPECT().
					GetVersion().
					Return(new(model.Version(devVersion)))
			},
			want: &expov1.GetVersionResponse{
				Version: &expov1.GetVersionResponse_Version{
					Value: "dev",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			vs := mocks.NewMockVersionService(ctrl)
			expoServiceSrv := expo.NewExpoServiceServer(vs)

			tt.prepare(vs)

			got, err := expoServiceSrv.GetVersion(t.Context(), &expov1.GetVersionRequest{})
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
