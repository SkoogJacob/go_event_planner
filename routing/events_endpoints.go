package routing

import (
	"event_planner_api/models"
	"net/http"
	"strings"
	"time"

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
		})
		return
	}
	update := make(map[string]string, 4)
	err = ctx.ShouldBindJSON(&update)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "unable to parse data as json"})
		return
	}
	var changes = 0
	if len(strings.TrimSpace(update["name"])) > 0 {
		event.Name = strings.TrimSpace(update["name"])
		changes++
	}
	if len(strings.TrimSpace(update["description"])) > 0 {
		event.Description = strings.TrimSpace(update["description"])
		changes++
	}
	if len(strings.TrimSpace(update["location"])) > 0 {
		event.Description = strings.TrimSpace(update["location"])
		changes++
	}
	if len(strings.TrimSpace(update["date_time"])) > 0 {
		date := strings.TrimSpace(update["date_time"])
		var dateTime time.Time
		dateTime, err = time.Parse(time.RFC3339, date)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "date_time field could not be parsed",
				"error":   err,
			})
			return
		}
		event.DateTime = dateTime
		changes++
	}
	if changes == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "The data submitted would not result in any changes",
		})
		return
	}
	err = event.UpdateEvent()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "error occurred when attempting to update event"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "event was successfully updated",
		"event":   event,
	})
}

func deleteEvent(ctx *gin.Context) {

}
