package authentication

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func AuthenticateByToken(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "not authorized"})
		return
	}
	_, userId, err := VerifyToken(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "unauthorized, bad token"})
		log.Printf("Error verifying token: %v\n", err)
		return
	}

	c.Set("user_id", userId)

	c.Next()
}
