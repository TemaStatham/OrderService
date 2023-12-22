package handler

import (
	"github.com/TemaStatham/OrderService/pkg/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	orders := router.Group("/orders")
	{
		orders.POST("/:id", h.getOrders)
		orders.GET("/:id", func(ctx *gin.Context) {
			h.showHTMLPage(ctx, "index.html", gin.H{"orderID": ctx.Param("id")})
		})
	}

	return router
}
