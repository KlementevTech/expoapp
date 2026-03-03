package public

import (
	"net/http"

	"expo/internal/api/public/middleware"
	"expo/internal/gen/openapi"

	"github.com/gin-gonic/gin"
)

func NewHTTPHandler(versionSvc VersionService) http.Handler {
	gin.SetMode(gin.ReleaseMode)

	handler := gin.New()
	handler.Use(
		gin.Recovery(),
		middleware.Logger(),
	)

	apiGroup := handler.Group("/api")
	openapi.RegisterHandlers(apiGroup, openapi.NewStrictHandler(&strictServerImpl{versionSvc: versionSvc}, nil))
	return handler
}
