package routing

import (
	"event_planner_api/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

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

func updateEvent(ctx *gin.Context) {
	id, err := getIdParam(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "passed bad id in the url, must be a positive number",
		})
		return
	}
	event, err := models.GetEvent(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "unable to get the event",
			"error":   err,
		})
	}
	update := make(map[string]interface{}, 5)
	ctx.ShouldBindJSON(update)
	log.Println(update, "\n", event)
	ctx.JSON(200, gin.H{"message": "under development"})
}

func deleteEvent(ctx *gin.Context) {

}
