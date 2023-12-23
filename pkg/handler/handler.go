package handler

import (
	"github.com/TemaStatham/OrderService/pkg/cache"
	"github.com/TemaStatham/OrderService/pkg/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
	c        *cache.Cache
}

func NewHandler(services *service.Service, c *cache.Cache) *Handler {
	return &Handler{services: services, c: c}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	// orders := router.Group("/orders")
	// {
	// 	orders.POST("", h.getOrders)
	// 	orders.GET("", func(ctx *gin.Context) {
	// 		h.showHTMLPage(ctx, "index.html", 1)
	// 	})
	// }

	router.POST("/orders", h.getOrders)
	router.GET("/orders", func(ctx *gin.Context) {
		h.showHTMLPage(ctx, "index.html", 1)
	})

	return router
}
