package handler

import (
	"L0-wbtech/internal/service"
	"L0-wbtech/pkg/errors"
	"L0-wbtech/pkg/logger/sl"
	"context"
	stdErrors "errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIHandler struct {
	service service.Service
	log     *slog.Logger
}

func New(service service.Service, log *slog.Logger) *APIHandler {
	return &APIHandler{
		service: service,
		log:     log,
	}
}

func (h *APIHandler) RegisterRoutes(router *gin.Engine) {

	router.GET("/order/:order_uid", h.GetOrder)

}

func (h *APIHandler) GetOrder(c *gin.Context) {
	const op = "handler.APIHandler.GetOrder"
	log := h.log.With(
		slog.String("op", op),
	)

	orderUID := c.Param("order_uid")
	if orderUID == "" {
		log.Error("order_uid is empty")
		c.JSON(http.StatusBadRequest, gin.H{"error": "order_uid is required"})
		return
	}

	order, err := h.service.GetOrder(context.Background(), orderUID)
	if err != nil {
		if stdErrors.Is(err, errors.ErrNotFound) {
			log.Warn("order not found", "order_uid", orderUID)
			c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
			return
		}

		log.Error("failed to get order", sl.Err(err), "order_uid", orderUID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, order)
}
