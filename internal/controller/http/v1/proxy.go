package v1

import (
	"errors"
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

	handler.POST("/proxies", r.createProxy)

	handler.GET("/proxies", r.getProxyList)
	handler.GET("/proxies/:proxy_id", r.getProxyById)
}

type createProxyRequest struct {
	Protocol string `json:"protocol" binding:"required"  example:"http"`
	Host     string `json:"Host"     binding:"required"  example:"127.0.0.1"`
	Port     int64  `json:"port"     binding:"required"  example:"8080"`
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
		if errors.Is(err, usecase.ErrInvalidData) {
			errorResponse(c, http.StatusBadRequest, err.Error())
		} else {
			errorResponse(c, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	c.JSON(http.StatusOK, createdProxy)
}

type getProxyListRequest struct {
	Offset int64 `form:"offset" example:"22"`
	Limit  int64 `form:"limit,default=20" example:"50"`
}

func (u *ProxyRoutes) getProxyList(c *gin.Context) {
	var req getProxyListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		u.l.Error(err, "http - v1 - getProxyList")
		errorResponse(c, http.StatusBadRequest, "invalid request")
		return
	}

	proxyList, err := u.u.GetProxyList(c, req.Offset, req.Limit)
	if err != nil {
		u.l.Error(err, "http - v1 - getProxyList")
		if errors.Is(err, usecase.ErrInvalidData) {
			errorResponse(c, http.StatusBadRequest, err.Error())
		} else {
			errorResponse(c, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	if proxyList.Proxies == nil {
		proxyList.Proxies = []domain.Proxy{}
	}
	c.JSON(http.StatusOK, proxyList)
}

type getProxyByIdRequest struct {
	ProxyId int64 `uri:"proxy_id" required:"true" example:"22"`
}

func (u *ProxyRoutes) getProxyById(c *gin.Context) {
	var req getProxyByIdRequest
	if err := c.ShouldBindUri(&req); err != nil {
		u.l.Error(err, "http - v1 - getProxyById")
		errorResponse(c, http.StatusBadRequest, "invalid request")
		return
	}

	proxy, err := u.u.GetProxyById(c, req.ProxyId)
	if err != nil {
		u.l.Error(err, "http - v1 - getProxyById")
		if errors.Is(err, usecase.ErrNotFound) {
			errorResponse(c, http.StatusNotFound, "proxy not found")
		} else {
			errorResponse(c, http.StatusInternalServerError, "internal server error")
		}
		return
	}
	c.JSON(http.StatusOK, proxy)
}
