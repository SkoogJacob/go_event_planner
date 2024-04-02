package routing

import (
	"errors"
	"event_planner_api/authentication"
	"log"
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
	api := router.Group("/api/")
	api.GET("/events", getEvents)
	api.GET("/events/:id", getEvent)
	api.POST("/signup", registerUser)
	api.POST("/login", loginUser)

	needsAuth := api.Group("/", authentication.AuthenticateByToken)
	needsAuth.POST("/events", postEvent)
	needsAuth.PUT("/events/:id", updateEvent)
	needsAuth.DELETE("/events/:id", deleteEvent)
	needsAuth.POST("/events/:id/register", registerUserForEvent)
	needsAuth.DELETE("/events/:id/register", unregisterUserForEvent)
}

func getIdParam(ctx *gin.Context) (int64, error) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "The ID requested could not be parsed to an integer",
			"error":   err,
		})
		log.Printf("Id could not be parsed from path: %s", ctx.Request.URL.Path)
		return -1, err
	} else if id < 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "The passed ID must not be negative",
		})
		log.Printf("negative id value in path: %s", ctx.Request.URL.Path)
		return id, errors.New("the passed ID was negative")
	}
	return id, nil
}
