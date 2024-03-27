package main

import (
	"context"
	"errors"
	"event_planner_api/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	router := gin.Default()
	router.GET("/events", getEvents)
	router.POST("/events", postEvent)

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

func getEvents(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "displaying all stored events",
		"events":  models.GetEvents(),
	})
}

func postEvent(ctx *gin.Context) {
	var event *models.Event = &models.Event{}
	err := ctx.ShouldBindJSON(event)
	if err == nil {
		event.Save()
		ctx.JSON(http.StatusCreated, gin.H{
			"message": "event created",
			"event":   *event,
		})
		return
	}
	ctx.JSON(http.StatusBadRequest, gin.H{
		"message":         "unable to parse json, check that all required fields were included and formatted correctly",
		"required_fields": "name, description, location, date_time",
	})
}
