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
	handler.GET("/proxies/:proxy_id", r.getProxy)
	handler.DELETE("/proxies/:proxy_id", r.deleteProxy)

	handler.GET("/proxies", r.getProxyList)

	handler.POST("/proxies/occupy", r.occupyMostAvailableProxy)
	handler.POST("/proxies/release", r.releaseProxy)
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
		u.l.Error(err.Error(), "http - v1 - createProxy")
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
		u.l.Error(err.Error(), "http - v1 - createProxy")
		if errors.Is(err, usecase.ErrInvalidData) {
			errorResponse(c, http.StatusBadRequest, err.Error())
		} else {
			errorResponse(c, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	c.JSON(http.StatusOK, createdProxy)
}

type getProxyByIdRequest struct {
	ProxyId int64 `uri:"proxy_id" binding:"required" example:"22"`
}

func (u *ProxyRoutes) getProxy(c *gin.Context) {
	var req getProxyByIdRequest
	if err := c.ShouldBindUri(&req); err != nil {
		u.l.Error(err.Error(), "http - v1 - getProxy")
		errorResponse(c, http.StatusBadRequest, "invalid request")
		return
	}

	proxy, err := u.u.GetProxy(c, req.ProxyId)
	if err != nil {
		u.l.Error(err, "http - v1 - getProxy")
		if errors.Is(err, usecase.ErrNotFound) {
			errorResponse(c, http.StatusNotFound, "proxy not found")
		} else {
			errorResponse(c, http.StatusInternalServerError, "internal server error")
		}
		return
	}
	c.JSON(http.StatusOK, proxy)
}

type deleteProxyRequest struct {
	ProxyId int64 `uri:"proxy_id" binding:"required" example:"22"`
}

func (u *ProxyRoutes) deleteProxy(c *gin.Context) {
	var req deleteProxyRequest
	if err := c.ShouldBindUri(&req); err != nil {
		u.l.Error(err.Error(), "http - v1 - deleteProxy")
		errorResponse(c, http.StatusBadRequest, "invalid request")
		return
	}

	if err := u.u.DeleteProxy(c, req.ProxyId); err != nil {
		u.l.Error(err.Error(), "http - v1 - deleteProxy")
		errorResponse(c, http.StatusInternalServerError, "internal server error")
	}
	c.Status(http.StatusNoContent)
}

type getProxyListRequest struct {
	Offset int64 `form:"offset" example:"22"`
	Limit  int64 `form:"limit,default=20" example:"50"`
}

func (u *ProxyRoutes) getProxyList(c *gin.Context) {
	var req getProxyListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		u.l.Error(err.Error(), "http - v1 - getProxyList")
		errorResponse(c, http.StatusBadRequest, "invalid request")
		return
	}

	proxyList, err := u.u.GetProxyList(c, req.Offset, req.Limit)
	if err != nil {
		u.l.Error(err.Error(), "http - v1 - getProxyList")
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

func (u *ProxyRoutes) occupyMostAvailableProxy(c *gin.Context) {
	proxyOccupy, err := u.u.OccupyMostAvailableProxy(c)
	if err != nil {
		u.l.Error(err.Error(), "http - v1 - occupyMostAvailableProxy")
		if errors.Is(err, usecase.ErrNotFound) {
			errorResponse(c, http.StatusNotFound, "not found any available proxy")
		} else {
			errorResponse(c, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	c.JSON(http.StatusOK, proxyOccupy)
}

type releaseProxyRequest struct {
	Key string `json:"key" binding:"required" example:"91af856e-f788-4e83-908e-153399961f35"`
}

func (u *ProxyRoutes) releaseProxy(c *gin.Context) {
	var req releaseProxyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		u.l.Error(err.Error(), "http - v1 - releaseProxy")
		errorResponse(c, http.StatusBadRequest, "invalid request")
		return
	}

	if err := u.u.ReleaseProxy(c, req.Key); err != nil {
		u.l.Error(err.Error(), "http - v1 - releaseProxy")
		errorResponse(c, http.StatusInternalServerError, "internal server error")
	}
	c.Status(http.StatusNoContent)
}
