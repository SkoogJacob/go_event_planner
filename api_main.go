package main

import (
	"context"
	"errors"
	"event_planner_api/event_db"
	"event_planner_api/models"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	event_db.InitDb("./data/events.db")
	defer event_db.CloseDb()

	router := gin.Default()
	router.GET("/api/events", getEvents)
	router.GET("/api/events/:id", getEvent)
	router.POST("/api/events", postEvent)
	router.PUT("/api/events/:id", updateEvent)
	router.DELETE("/api/events/:id", deleteEvent)
	router.POST("/api/signup", registerUser)
	router.POST("/api/login", loginUser)
	router.POST("/api/events/:id/register", registerUserForEvent)
	router.DELETE("/api/events/:id/register", unregisterUserForEvent)

	server := http.Server{Addr: ":9000", Handler: router}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds")
	}
	log.Println("Server exiting")
}

func getEvent(ctx *gin.Context) {
	id, err := getIdParam(ctx)
	if err != nil {
		return
	}
	event, err := models.GetEvent(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "unable to fetch event",
			"id":      id,
			"error":   err,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "successfully retrieved event",
		"event":   event,
	})
}

func unregisterUserForEvent(c *gin.Context) {

}

func registerUserForEvent(c *gin.Context) {

}

func loginUser(c *gin.Context) {

}

func registerUser(c *gin.Context) {

}

func deleteEvent(c *gin.Context) {

}

func updateEvent(c *gin.Context) {

}

func getEvents(ctx *gin.Context) {
	events, err := models.GetEvents()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to read from internal DB",
			"error":   err,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "displaying all stored events",
		"events":  events,
	})
}

func postEvent(ctx *gin.Context) {
	var event = &models.Event{}
	err := ctx.ShouldBindJSON(event)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message":         "unable to parse json, check that all required fields were included and formatted correctly",
			"required_fields": "name, description, location, date_time",
		})
		return
	}
	err = event.Save()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to save event to internal DB",
			"error":   err,
		})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "event created",
		"event":   *event,
	})
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
