package routing

import (
	"event_planner_api/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func unregisterUserForEvent(c *gin.Context) {

}

func registerUserForEvent(c *gin.Context) {

}

func loginUser(c *gin.Context) {
	var user models.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message":           "could not decode request, are all obligatory fields included?",
			"obligatory_fields": "email, password",
		})
		return
	}
	authenticated, err := user.Login()
	if !authenticated {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "the credentials did not match"})
		log.Printf("Error in authenticating user %v: %v\n", user, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "login successful!"})
}

func registerUser(c *gin.Context) {
	var user models.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message":           "could not decode request, are all obligatory fields included?",
			"obligatory_fields": "email, password",
		})
		return
	}
	err = user.Save()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "could not save the user to the database",
		})
		log.Printf("Error in attempting to save user to database: %v\n", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "successfully created user",
	})
}
