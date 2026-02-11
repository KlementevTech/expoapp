package expo_test

import (
	"context"
	"expoapp/internal/domain"
	"expoapp/internal/web/expo/mocks"
	"testing"

	"expoapp/internal/web/expo"
	expopb "expoapp/pkg/api/expo"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const devVersion = "dev"

func TestService_GetInfo(t *testing.T) {
	tests := []struct {
		name    string
		prepare func(mock *mocks.MockVersionService)
		want    *expopb.GetInfoResponse
	}{
		{
			name: "success",
			prepare: func(versions *mocks.MockVersionService) {
				versions.EXPECT().
					GetVersion().
					Return(domain.NewVersion(devVersion))
			},
			want: &expopb.GetInfoResponse{
				Version: &expopb.GetInfoResponse_Version{
					Version: "dev",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			versions := mocks.NewMockVersionService(ctrl)
			expoService := expo.NewService(versions)

			tt.prepare(versions)

			got, err := expoService.GetInfo(context.Background(), &expopb.GetInfoRequest{})
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
