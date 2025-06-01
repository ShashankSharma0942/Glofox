package route

import (
	"net/http"
	"sync"

	"glofox/config"
	mapstore "glofox/core"
	"glofox/internal/handler"
	"glofox/internal/service"

	"github.com/gin-gonic/gin"
)

// router holds dependencies and the Gin engine for defining and managing routes.
type router struct {
	gin      *gin.Engine
	syMap    mapstore.MapStore
	lock     *sync.Mutex
	cfg      config.Config
	services service.BusinessService
}

// NewRouter initializes a new router with provided dependencies.
// It prepares the Gin engine and returns the router wrapper.
func NewRouter(syMap mapstore.MapStore, lock *sync.Mutex, cfg config.Config, services service.BusinessService) *router {
	return &router{
		gin:      gin.Default(),
		syMap:    syMap,
		lock:     lock,
		services: services,
		cfg:      cfg,
	}
}

// SetRoutes defines the API endpoints and attaches route groups.
// It returns the configured HTTP handler for the server to use.
func (router *router) SetRoutes() http.Handler {
	baseGrp := router.gin.Group(router.cfg.BaseRoute)
	{
		router.Class(baseGrp)
		router.Booking(baseGrp)
	}
	return router.gin.Handler()
}

// Class registers the endpoint for class creation under the given route group.
func (router *router) Class(rg *gin.RouterGroup) {
	handle := handler.NewClassHandler(router.syMap, router.lock, router.services)
	{
		rg.POST("/class", handle.CreateClass) // POST /class to create a new class
	}
}

// Booking registers the endpoint for class booking under the given route group.
func (router *router) Booking(rg *gin.RouterGroup) {
	handle := handler.NewBookingHandler(router.syMap, router.lock, router.services)
	{
		rg.POST("/booking", handle.CreateBooking) // POST /booking to book a class
	}
}
