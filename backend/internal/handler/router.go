package handler

import (
	"github.com/gin-gonic/gin"
)

func (h *APIHandler) RegisterRoutes(router *gin.Engine) {

	router.Use(corsMiddleware())

	router.GET("/order/:order_uid", h.GetOrder)

}
