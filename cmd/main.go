package main

import (
	"glofox/cmd/server"
	"glofox/config"
	"glofox/constants"
	mapstore "glofox/core"
	"glofox/internal/service"
	"log"
	"sync"
)

// main is the entry point of the application.
// It initializes shared resources, sets up the business services,
// and starts the HTTP server.
func main() {
	// Load configuration from config.json
	cfg, err := config.LoadConfig(constants.FilePath)
	if err != nil {
		log.Fatalf(constants.Failepath, err)
	}
	// Initialize a shared mutex for synchronizing access to the map store
	lock := &sync.Mutex{}

	// Create a thread-safe map store instance
	reqMap := mapstore.NewMuMapStore(lock)

	// Initialize the application's business logic layer with shared state
	services := service.InitializeService(reqMap, lock, *cfg)

	// Create a new HTTP server using the configured port
	newServer := server.NewServer(*cfg)

	// Start the server and listen for incoming requests
	newServer.RunServer(reqMap, lock, services)
}
