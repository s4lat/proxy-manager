package v1

import (
	"errors"
	"net/http"
	"proxy_manager/internal/domain"
	"proxy_manager/internal/usecase"
	"proxy_manager/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
)

type ProxyRoutes struct {
	u usecase.UseCase
	l logger.Interface
}

func newProxyRoutes(handler *gin.RouterGroup, u usecase.UseCase, l logger.Interface) {
	r := &ProxyRoutes{u: u, l: l}

	handler.POST("/proxies", r.createProxy)
	handler.GET("/proxies/:proxyID", r.getProxy)
	handler.PUT("/proxies/:proxyID", r.updateProxy)
	handler.DELETE("/proxies/:proxyID", r.deleteProxy)

	handler.GET("/proxies", r.getProxyList)

	handler.POST("/proxies/occupy", r.occupyMostAvailableProxy)
	handler.POST("/proxies/release", r.releaseProxy)
}

type createProxyRequest struct {
	Protocol       string    `json:"protocol"       binding:"required" example:"http"                     extensions:"x-order=1"`
	Host           string    `json:"Host"           binding:"required" example:"127.0.0.1"                extensions:"x-order=2"`
	Port           int64     `json:"port"           binding:"required" example:"8080"                     extensions:"x-order=3"`
	Username       string    `json:"username"                          example:"login123"                 extensions:"x-order=4"`
	Password       string    `json:"password"                          example:"qwerty1234"               extensions:"x-order=5"`
	ExpirationDate time.Time `json:"expirationDate" binding:"required" example:"2025-02-18T21:54:42.123Z" extensions:"x-order=6"`
}

// createProxy godoc
//
//	@Summary		Create proxy
//	@Description	Creates proxy with given params and returns it
//	@Tags			proxies
//	@Accept			json
//	@Produce		json
//	@Param			request	body		createProxyRequest	true	"Create proxy"
//	@Success		200		{object}	domain.Proxy
//	@Failure		400		{object}	errResponse
//	@Failure		500		{object}	errResponse
//	@Router			/proxies [POST]
func (u *ProxyRoutes) createProxy(c *gin.Context) {
	var req createProxyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		u.l.Error("http - v1 - createProxy - %s", err)
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	createdProxy, err := u.u.CreateProxy(c, domain.Proxy{
		Protocol:       req.Protocol,
		Username:       req.Username,
		Password:       req.Password,
		Host:           req.Host,
		Port:           req.Port,
		ExpirationDate: req.ExpirationDate,
	})
	if err != nil {
		u.l.Error("http - v1 - createProxy - %s", err)
		if errors.Is(err, usecase.ErrInvalidData) {
			errorResponse(c, http.StatusBadRequest, err.Error())
		} else {
			errorResponse(c, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	c.JSON(http.StatusOK, createdProxy)
}

type getProxyRequest struct {
	ProxyID int64 `uri:"proxyID" binding:"required" example:"22"`
}

// getProxy godoc
//
//	@Summary		Get proxy
//	@Description	Returns proxy with given ID
//	@Tags			proxies
//	@Produce		json
//	@Param			proxyID	path		int64	true	"Proxy ID"
//	@Success		200		{object}	domain.Proxy
//	@Failure		400		{object}	errResponse
//	@Failure		404		{object}	errResponse
//	@Failure		500		{object}	errResponse
//	@Router			/proxies/{proxyID} [GET]
func (u *ProxyRoutes) getProxy(c *gin.Context) {
	var req getProxyRequest
	if err := c.ShouldBindUri(&req); err != nil {
		u.l.Error("http - v1 - getProxy - %s", err)
		errorResponse(c, http.StatusBadRequest, "invalid request")
		return
	}

	proxy, err := u.u.GetProxy(c, req.ProxyID)
	if err != nil {
		u.l.Error("http - v1 - getProxy - %s", err)
		if errors.Is(err, usecase.ErrNotFound) {
			errorResponse(c, http.StatusNotFound, "proxy not found")
		} else {
			errorResponse(c, http.StatusInternalServerError, "internal server error")
		}
		return
	}
	c.JSON(http.StatusOK, proxy)
}

type getProxyBeforeUpdateRequest struct {
	ProxyID int64 `uri:"proxyID" binding:"required" example:"22"`
}

type updateProxyRequest struct {
	Protocol       string    `json:"protocol"       binding:"required" example:"http"                     extensions:"x-order=1"`
	Host           string    `json:"host"           binding:"required" example:"127.0.0.1"                extensions:"x-order=2"`
	Port           int64     `json:"port"           binding:"required" example:"8080"                     extensions:"x-order=3"`
	Username       string    `json:"username"                          example:"login123"                 extensions:"x-order=4"`
	Password       string    `json:"password"                          example:"qwerty1234"               extensions:"x-order=5"`
	ExpirationDate time.Time `json:"expirationDate" binding:"required" example:"2025-02-18T21:54:42.123Z" extensions:"x-order=7"`
}

// updateProxy godoc
//
//	@Summary		Update proxy
//	@Description	Updates proxy with given ID
//	@Tags			proxies
//	@Accept			json
//	@Produce		json
//	@Param			proxyID	path		int64				true	"Proxy ID"
//	@Param			request	body		updateProxyRequest	true	"Proxy data"
//	@Success		200		{object}	domain.Proxy
//	@Failure		400		{object}	errResponse
//	@Failure		404		{object}	errResponse
//	@Failure		500		{object}	errResponse
//	@Router			/proxies/{proxyID} [PUT]
func (u *ProxyRoutes) updateProxy(c *gin.Context) {
	var getProxyReq getProxyBeforeUpdateRequest
	if err := c.ShouldBindUri(&getProxyReq); err != nil {
		u.l.Error("http - v1 - updateProxy - %s", err)
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	var updateProxyReq updateProxyRequest
	if err := c.ShouldBindJSON(&updateProxyReq); err != nil {
		u.l.Error("http - v1 - updateProxy - %s", err)
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	updatedProxy, err := u.u.UpdateProxy(c, domain.Proxy{
		ID:             getProxyReq.ProxyID,
		Protocol:       updateProxyReq.Protocol,
		Username:       updateProxyReq.Username,
		Password:       updateProxyReq.Password,
		Host:           updateProxyReq.Host,
		Port:           updateProxyReq.Port,
		ExpirationDate: updateProxyReq.ExpirationDate,
	})
	if err != nil {
		u.l.Error("http - v1 - updateProxy - %s", err)
		if errors.Is(err, usecase.ErrNotFound) {
			errorResponse(c, http.StatusNotFound, "proxy not found")
		} else {
			errorResponse(c, http.StatusInternalServerError, "internal server error")
		}
		return
	}
	c.JSON(http.StatusOK, updatedProxy)
}

type deleteProxyRequest struct {
	ProxyID int64 `uri:"proxyID" binding:"required" example:"22"`
}

// deleteProxy godoc
//
//	@Summary		Delete proxy
//	@Description	Deletes proxy with given ID
//	@Tags			proxies
//	@Produce		json
//	@Param			proxyID	path	int64	true	"Proxy ID"
//	@Success		204		"No content"
//	@Failure		400		{object}	errResponse
//	@Failure		500		{object}	errResponse
//	@Router			/proxies/{proxyID} [DELETE]
func (u *ProxyRoutes) deleteProxy(c *gin.Context) {
	var req deleteProxyRequest
	if err := c.ShouldBindUri(&req); err != nil {
		u.l.Error("http - v1 - deleteProxy - %s", err)
		errorResponse(c, http.StatusBadRequest, "invalid request")
		return
	}

	if err := u.u.DeleteProxy(c, req.ProxyID); err != nil {
		u.l.Error("http - v1 - deleteProxy - %s", err)
		errorResponse(c, http.StatusInternalServerError, "internal server error")
	}
	c.Status(http.StatusNoContent)
}

type getProxyListRequest struct {
	Offset int64 `form:"offset" example:"22"`
	Limit  int64 `form:"limit,default=20" example:"50"`
}

// getProxyList godoc
//
//	@Summary		Get proxy list
//	@Description	Returns proxy list
//	@Tags			proxies
//	@Produce		json
//	@Param			offset	query		int64	false	"Offset in proxy list"
//	@Param			limit	query		int64	false	"Limit of proxy list size"
//	@Success		200		{object}	domain.ProxyList
//	@Failure		400		{object}	errResponse
//	@Failure		500		{object}	errResponse
//	@Router			/proxies [GET]
func (u *ProxyRoutes) getProxyList(c *gin.Context) {
	var req getProxyListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		u.l.Error("http - v1 - getProxyList- %s", err)
		errorResponse(c, http.StatusBadRequest, "invalid request")
		return
	}

	proxyList, err := u.u.GetProxyList(c, req.Offset, req.Limit)
	if err != nil {
		u.l.Error("http - v1 - getProxyList - %s", err)
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

// occupyMostAvailableProxy godoc
//
//	@Summary		Occupy most available proxy
//	@Description	Occupies the most available proxy, returns its info and key to release
//	@Tags			proxies
//	@Produce		json
//	@Success		200	{object}	domain.ProxyOccupy
//	@Failure		400	{object}	errResponse
//	@Failure		404	{object}	errResponse
//	@Failure		500	{object}	errResponse
//	@Router			/proxies/occupy [POST]
func (u *ProxyRoutes) occupyMostAvailableProxy(c *gin.Context) {
	proxyOccupy, err := u.u.OccupyMostAvailableProxy(c)
	if err != nil {
		u.l.Error("http - v1 - occupyMostAvailableProxy - %s", err)
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

// releaseProxy godoc
//
//	@Summary		Release proxy occupy
//	@Description	Releases proxy occupy with given key
//	@Tags			proxies
//	@Produce		json
//	@Accept			json
//	@Param			key	body	releaseProxyRequest	true	"Key of occupy"
//	@Success		204	"No content"
//	@Failure		400	{object}	errResponse
//	@Failure		500	{object}	errResponse
//	@Router			/proxies/release [POST]
func (u *ProxyRoutes) releaseProxy(c *gin.Context) {
	var req releaseProxyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		u.l.Error("http - v1 - releaseProxy - %s", err)
		errorResponse(c, http.StatusBadRequest, "invalid request")
		return
	}

	if err := u.u.ReleaseProxy(c, req.Key); err != nil {
		u.l.Error("http - v1 - releaseProxy - %s", err)
		errorResponse(c, http.StatusInternalServerError, "internal server error")
	}
	c.Status(http.StatusNoContent)
}
