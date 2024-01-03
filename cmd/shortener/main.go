package main

import (
	"github.com/gin-gonic/gin"
	"github.com/marsel1323/url-shortener-go/internal/app"
	"log"
)

func main() {
	r := gin.Default()

	storage := app.NewInMemoryStorage()
	server := app.NewServer(storage)

	r.POST("/", server.HandlePostRequest)
	r.POST("/api/shorten", server.HandleAPIShorten)
	r.GET("/:id", server.HandleGetRequest)

	log.Fatal(r.Run("localhost:8080"))
}
