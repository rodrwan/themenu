package main

import (
	"log"
	"os"

	"github.com/rodrwan/themenu/internal/cqrs"
	"github.com/rodrwan/themenu/internal/web"
)

func main() {
	eventBus := cqrs.NewEventBus()
	apiClient := web.NewAPIClient("http://themenu-api:8080")
	server := web.NewServer(eventBus, apiClient)

	// Iniciar el servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	if err := server.Start(":" + port); err != nil {
		log.Fatal(err)
	}
}
