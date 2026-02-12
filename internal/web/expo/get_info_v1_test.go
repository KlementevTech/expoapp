package expo_test

import (
	"context"
	"expoapp/internal/domain"
	"expoapp/internal/web/expo/mocks"
	"testing"

	"expoapp/internal/web/expo"
	"expoapp/pkg/api/pb"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const devVersion = "dev"

func TestService_GetInfoV1(t *testing.T) {
	tests := []struct {
		name    string
		prepare func(mock *mocks.MockVersionProvider)
		want    *pb.GetInfoV1Response
	}{
		{
			name: "success",
			prepare: func(versions *mocks.MockVersionProvider) {
				versions.EXPECT().
					GetVersion().
					Return(domain.NewVersion(devVersion))
			},
			want: &pb.GetInfoV1Response{
				Version: &pb.GetInfoV1Response_Version{
					Version: "dev",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			versions := mocks.NewMockVersionProvider(ctrl)
			expoService := expo.NewService(versions)

			tt.prepare(versions)

			got, err := expoService.GetInfoV1(context.Background(), &pb.GetInfoV1Request{})
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
