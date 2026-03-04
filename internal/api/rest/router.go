package rest

import (
	"net/http"

	"expo/gen/openapi"
	"expo/internal/api/rest/middleware"

	"github.com/gin-gonic/gin"
)

func NewHandler() http.Handler {
	gin.SetMode(gin.ReleaseMode)

	handler := gin.New()
	handler.Use(
		gin.Recovery(),
		middleware.Logger(),
	)

	apiGroup := handler.Group("/api")
	openapi.RegisterHandlers(apiGroup, openapi.NewStrictHandler(&strictServerImpl{}, nil))
	return handler
}
