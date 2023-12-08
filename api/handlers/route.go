package handlers

import (
	"github.com/gin-gonic/gin"
)

func (h *CncHandler) SetupRouter(r *gin.Engine) {
	r.GET("/", h.hello)
	r.GET("/state", h.state)
	r.POST("/config", h.configConnexion)
	r.POST("/command", h.handelCommand)
	r.POST("/file", h.handelFile)
}
