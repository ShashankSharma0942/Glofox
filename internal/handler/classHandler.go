package handler

import (
	"glofox/constants"
	mapstore "glofox/core"
	newError "glofox/errors"
	"glofox/internal/service"
	"glofox/models/dto"
	"glofox/utils"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

// ClassHandler defines the interface for handling class-related HTTP requests.
type ClassHandler interface {
	CreateClass(c *gin.Context)
}

// class is the concrete implementation of ClassHandler.
// It provides logic for handling class creation requests.
// It uses shared map storage, a mutex for thread safety, and a business service for core logic.
type class struct {
	syMap   mapstore.MapStore
	lock    *sync.Mutex
	service service.BusinessService
}

// NewClassHandler creates a new instance of ClassHandler with dependencies injected.
// This sets up the handler to be used in HTTP routing.
func NewClassHandler(syMap mapstore.MapStore, lock *sync.Mutex, services service.BusinessService) ClassHandler {
	return &class{
		syMap:   syMap,
		lock:    lock,
		service: services,
	}
}

// CreateClass handles POST /class endpoint.
// It validates the request payload, invokes the service layer to persist class data,
// and returns a JSON response with the result.
func (class *class) CreateClass(c *gin.Context) {
	var classData dto.Class

	// Bind and validate the incoming JSON payload
	err := c.ShouldBindJSON(&classData)
	if err != nil {
		log.Println(newError.ErrUnmarshalling.Error(), err.Error())
		c.JSON(http.StatusBadRequest, utils.CreateResp(false, newError.ErrUnmarshalling.Error()))
		return
	}

	// Call business logic to handle class creation
	err = class.service.CreateClass(classData)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.CreateResp(false, err.Error()))
		return
	}

	// Respond with success if everything went well
	c.JSON(http.StatusOK, utils.CreateResp(true, constants.ClassSuccess))
}
