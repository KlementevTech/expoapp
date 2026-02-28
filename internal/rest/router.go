package rest

import (
	"expo/internal/gen/openapi"
	"expo/internal/rest/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(vs VersionService) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(
		gin.Recovery(),
		middleware.Logger(),
	)

	api := router.Group("/api")
	openapi.RegisterHandlers(api, openapi.NewStrictHandler(&server{vs: vs}, nil))

	return router
}
