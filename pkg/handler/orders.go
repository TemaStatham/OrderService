package handler

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Message struct {
	ID string `json:"id" binding:"required"`
}

func (h *Handler) showHTMLPage(c *gin.Context, templateName string, data interface{}) {
	tmpl, err := template.ParseFiles("./static/" + templateName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	tmpl.Execute(c.Writer, data)
}

func (h *Handler) getOrders(c *gin.Context) {
	var requestData Message

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, f := h.c.Get(requestData.ID)
	if f {
		o := *order
		c.JSON(http.StatusOK, gin.H{
			"message": "Data received successfully",
			"order":   o,
		})
		return
	}
	fmt.Print("check order in db\n")

	order, err := h.services.GetOrder(requestData.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error db": err.Error()})
		return
	}
	o := *order
	c.JSON(http.StatusOK, gin.H{
		"message": "Data received successfully",
		"order":   o,
	})
}
