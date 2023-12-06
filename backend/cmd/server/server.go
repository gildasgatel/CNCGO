package main

import (
	"cncgo/backend/internal/api/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	// Créer une instance de Gin
	router := gin.Default()

	var hand handlers.CncHandler
	hand.SetupRouter(router)
	defer hand.Close()

	// Démarre le serveur sur le port 8080
	router.Run(":8080")
}
