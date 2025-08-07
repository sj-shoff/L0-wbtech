package handler

import (
	"L0-wbtech/internal/service"
	"L0-wbtech/pkg/errors"
	"L0-wbtech/pkg/logger/sl"
	"context"
	stdErrors "errors"
	"log/slog"
	"net/http"
	"time"

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

	router.Use(corsMiddleware())

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

	start := time.Now()

	order, err := h.service.GetOrder(context.Background(), orderUID)

	dataFetchTime := time.Since(start)
	c.Set("data_fetch_start", start)

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

	log.Debug("Data fetch completed",
		"order_uid", orderUID,
		"data_fetch_time", dataFetchTime,
	)
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		c.Next()
	}
}
