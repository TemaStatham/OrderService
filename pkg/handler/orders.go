package handler

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) showHTMLPage(c *gin.Context, templateName string, data interface{}) {
	tmpl, err := template.ParseFiles("./static/" + templateName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	tmpl.Execute(c.Writer, data)
}

func (h *Handler) getOrders(c *gin.Context) {

}
