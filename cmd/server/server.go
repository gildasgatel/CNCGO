package main

import (
	"cncgo/api/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	var hand handlers.CncHandler
	hand.SetupRouter(router)
	defer hand.Close()

	router.Run(":8080")
}
