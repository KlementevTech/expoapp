package rest

import (
	"net/http"

	"expo/internal/gen/openapi"
	"expo/internal/rest/middleware"

	"github.com/gin-gonic/gin"
)

func NewHTTPRouter(versionSvc VersionService) http.Handler {
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
