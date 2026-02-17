package expo_test

import (
	"testing"

	"expo/internal/handlers/expo"
	"expo/internal/handlers/expo/mocks"
	"expo/internal/model"
	expov1 "expo/pkg/pb/expo/v1"

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
			expoService := expo.NewService(vs)

			tt.prepare(vs)

			got, err := expoService.GetVersion(t.Context(), &expov1.GetVersionRequest{})
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
