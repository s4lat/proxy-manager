package v1

import (
	_ "proxy_manager/docs"
	"proxy_manager/internal/usecase"
	"proxy_manager/pkg/logger"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// NewRouter godoc
// Swagger spec:
//
//	@title			Proxy Manager API
//	@description	Proxy Manager API documentation
//	@version		1.0
//	@BasePath		/api/v1
func NewRouter(handler *gin.Engine, u usecase.UseCase, l logger.Interface, serveSwag bool) {
	handler.Use(JSONLogMiddleware(l))
	handler.Use(gin.Recovery())

	h := handler.Group("/api/v1")
	{
		if serveSwag {
			h.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
		}
		newProxyRoutes(h, u, l)
	}
}
