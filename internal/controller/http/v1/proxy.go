package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"proxy_manager/internal/domain"
	"proxy_manager/internal/usecase"
	"proxy_manager/pkg/logger"
)

// TODO:
//   Прокинуть логгер

type ProxyRoutes struct {
	u usecase.UseCase
	l logger.Interface
}

func newProxyRoutes(handler *gin.RouterGroup, u usecase.UseCase, l logger.Interface) {
	r := &ProxyRoutes{u: u, l: l}

	handler.POST("/proxy", r.createProxy)
}

type createProxyRequest struct {
	Protocol string `json:"protocol" binding:"required"  example:"http"`
	Host     string `json:"Host"     binding:"required"  example:"127.0.0.1"`
	Port     int    `json:"port"     binding:"required"  example:"8080"`
	Username string `json:"username" example:"login123"`
	Password string `json:"password" example:"qwerty1234"`
}

func (u *ProxyRoutes) createProxy(c *gin.Context) {
	var req createProxyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		u.l.Error(err, "http - v1 - createProxy")
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	createdProxy, err := u.u.CreateProxy(c, domain.Proxy{
		Protocol: req.Protocol,
		Username: req.Username,
		Password: req.Password,
		Host:     req.Host,
		Port:     req.Port,
	})
	if err != nil {
		u.l.Error(err, "http - v1 - createProxy")
		errorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	c.JSON(http.StatusOK, createdProxy)
}
