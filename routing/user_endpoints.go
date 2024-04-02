package routing

import (
	"event_planner_api/authentication"
	"event_planner_api/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

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
	err = user.Login()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "the credentials did not match"})
		log.Printf("Error in authenticating user %v: %v\n", user, err)
		return
	}
	token, err := authentication.GenerateToken(user.Email, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error during login"})
		log.Printf("Error in getting JWT token: %v\n", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "login successful!", "token": token})
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
	token, err := authentication.GenerateToken(user.Email, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "user was created successfully, but the server failed to generate an auth token"})
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "successfully created user",
		"token":   token,
	})
}
