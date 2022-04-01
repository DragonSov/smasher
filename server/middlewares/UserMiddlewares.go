package middlewares

import (
	"github.com/DragonSov/smasher/server/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"strings"
)

type UserMiddlewares struct {
	Service services.UserService
}

func (m UserMiddlewares) CheckAuthorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Getting the authorization token from the headers
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(403, gin.H{
				"status":  "error",
				"message": "For this action you need to sign in",
			})
			c.Abort()
			return
		}

		// Getting a clean (no Bearer) authorization token
		authToken := strings.Split(authHeader, " ")
		if len(authToken) != 2 || authToken[0] != "Bearer" {
			c.JSON(403, gin.H{
				"status":  "error",
				"message": "Invalid Authorization header. Please try to sign in again",
			})
			c.Abort()
			return
		}

		// Receiving token claims (JWT decoding)
		tokenClaims, err := services.DecodeJWT(authToken[1])
		if err != nil {
			c.JSON(403, gin.H{
				"status":  "error",
				"message": "Invalid JWT token. Please try to sign in again",
			})
			c.Abort()
			return
		}

		// Parsing the user ID from the JWT subject
		userId, err := uuid.Parse(tokenClaims.Subject)
		if err != nil {
			c.JSON(500, gin.H{
				"status":  "error",
				"message": "Failed to get the user ID. Try again later",
			})
			c.Abort()
			return
		}

		// Selecting user by ID
		user, err := m.Service.SelectUserByUUID(userId)
		if err != nil || user == nil {
			c.JSON(403, gin.H{
				"status":  "error",
				"message": "User not found. Please try to sign in again",
			})
			c.Abort()
			return
		}

		// Setting up a custom model in context
		c.Set("authUser", *user)
		c.Next()
	}
}
