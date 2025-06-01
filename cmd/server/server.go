package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"glofox/config"
	mapstore "glofox/core"
	route "glofox/internal/gin"
	"glofox/internal/service"

	"github.com/gin-gonic/gin"
)

// Server interface defines the method required to start the application server.
type Server interface {
	RunServer(syMap mapstore.MapStore, lock *sync.Mutex, services service.BusinessService)
}

// server is a concrete implementation of the Server interface.
// It holds the HTTP server, router, port configuration, and server status.
type server struct {
	router *gin.Engine  // Gin router instance
	port   string       // Port the server listens on
	http   *http.Server // Underlying HTTP server
	Status bool         // Indicates if the server is up
	config config.Config
}

// NewServer returns a new instance of the server with the given port.
func NewServer(cfg config.Config) Server {
	return &server{
		config: cfg,
	}
}

// RunServer initializes and starts the server, and listens for termination signals
// to perform graceful shutdown when necessary.
func (serverInfo *server) RunServer(syMap mapstore.MapStore, lock *sync.Mutex, services service.BusinessService) {
	serverInfo.start(syMap, lock, services)
	serverInfo.gracefulShutdown()
}

// start initializes the HTTP server with routing and starts it asynchronously.
func (serverInfo *server) start(syMap mapstore.MapStore, lock *sync.Mutex, services service.BusinessService) {
	serverInfo.http = &http.Server{
		Addr:              ":" + serverInfo.config.Port,                                          // Bind server to specified port
		Handler:           route.NewRouter(syMap, lock, serverInfo.config, services).SetRoutes(), // Set up routing
		ReadHeaderTimeout: 20 * time.Second,                                                      // Prevent slowloris attacks by setting header timeout
	}

	// Start server in a separate goroutine to allow graceful shutdown
	go func(serv server) {
		serv.Status = true
		if err := serv.http.ListenAndServe(); err != nil {
			serv.Status = false
			if errors.Is(err, http.ErrServerClosed) {
				log.Println("Error: server was closed gracefully")
				return
			}
			log.Println("Error: failed to start server:", err)
		} else {
			log.Printf("Server successfully started on port %s\n", serv.port)
		}
	}(*serverInfo)
}

// listenToSignalNotification blocks until an OS shutdown signal is received.
// It listens for SIGINT and SIGTERM.
func listenToSignalNotification() {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // Block until signal is received
}

// gracefulShutdown is triggered after receiving a shutdown signal.
// It gracefully shuts down the server, closing connections properly.
func (serverInfo *server) gracefulShutdown() {
	listenToSignalNotification()

	err := serverInfo.http.Shutdown(context.Background())
	if err != nil {
		log.Fatalf("Graceful shutdown failed: %s. Forcing shutdown...", err.Error())
	} else {
		log.Println("Server shut down gracefully")
	}
}
