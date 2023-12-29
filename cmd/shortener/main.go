package main

import (
	"github.com/marsel1323/url-shortener-go/internal/app"
	"log"
)

func main() {
	storage := app.NewInMemoryStorage()
	server := app.NewServer(storage)

	log.Fatal(server.ListenAndServe())
}
