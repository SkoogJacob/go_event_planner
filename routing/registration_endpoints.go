package routing

import (
	"event_planner_api/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func registerUserForEvent(c *gin.Context) {
	eventId, err := getIdParam(c)
	if err != nil {
		return
	}
	event, err := models.GetEvent(eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "unable to fetch the event"})
		log.Printf("unable to fetch event with id %s: %s\n", eventId, err)
	}
	userId := c.GetInt64("user_id")
	err = event.Register(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not register for event"})
		log.Printf("Error in trying to register user with id %s for event with id %s: %s\n",
			userId, event.ID, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "registration successful"})
}

func unregisterUserForEvent(c *gin.Context) {
	eventId, err := getIdParam(c)
	if err != nil {
		return
	}
	event, err := models.GetEvent(eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "unable to fetch the event"})
		log.Printf("unable to fetch event with id %s: %s\n", eventId, err)
	}
	userId := c.GetInt64("user_id")
	err = event.DeRegister(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "unable to de-register from event"})
		log.Printf("Unable to de-register user from event: %s", err)
	}
	c.JSON(http.StatusOK, gin.H{"message": "de-registered from event"})
}
