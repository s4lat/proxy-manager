package v1

import (
	"github.com/gin-gonic/gin"
	"proxy_manager/internal/usecase"
	"proxy_manager/pkg/logger"
)

func NewRouter(handler *gin.Engine, u usecase.UseCase, l logger.Interface) {
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	h := handler.Group("/api/v1")
	{
		newProxyRoutes(h, u, l)
	}
}
