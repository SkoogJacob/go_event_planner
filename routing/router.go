package routing

import (
	"errors"
	"event_planner_api/authentication"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func MakeServer(serverAddress string) *http.Server {
	router := gin.Default()
	registerRoutes(router)
	server := http.Server{Addr: serverAddress, Handler: router}
	return &server
}

func registerRoutes(router *gin.Engine) {
	router.GET("/api/events", getEvents)
	router.GET("/api/events/:id", getEvent)
	router.POST("/api/signup", registerUser)
	router.POST("/api/login", loginUser)
	needsAuth := router.Group("/", authentication.AuthenticateByToken)
	needsAuth.POST("/api/events", postEvent)
	needsAuth.PUT("/api/events/:id", updateEvent)
	needsAuth.DELETE("/api/events/:id", deleteEvent)
	needsAuth.POST("/api/events/:id/register", registerUserForEvent)
	needsAuth.DELETE("/api/events/:id/register", unregisterUserForEvent)
}

func getIdParam(ctx *gin.Context) (int64, error) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "The ID requested could not be parsed to an integer",
			"error":   err,
		})
		return -1, err
	} else if id < 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "The passed ID must not be negative",
		})
		return id, errors.New("the passed ID was negative")
	}
	return id, nil
}
